package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/ReneKroon/ttlcache/v2"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/api/idtoken"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	pb "gomodules.avm99963.com/twpt-server/api_proto"
	db "gomodules.avm99963.com/twpt-server/internal/db"
)

var (
	authorizedUsersCacheTime = 15 * time.Minute
)

var (
	port              = flag.Int("port", 10000, "The server port")
	dbDsn             = flag.String("db", "", "MySQL/MariaDB database data source name (DSN)") // https://github.com/go-sql-driver/mysql#dsn-data-source-name
	dbConnMaxLifetime = flag.Int("dbConnMaxLifetime", 3*60, "Maximum amount of time a connection to the database may be reused in seconds.")
	dbMaxOpenConns    = flag.Int("dbMaxOpenConns", 5, "Maximum number of open connections to the database.")
	dbMaxIdleConns    = flag.Int("dbMaxIdleConns", *dbMaxOpenConns, "Maximum number of connections to the database in the idle connection pool.")
	jwtAudience       = flag.String("jwtAudience", "", "JWT audience string.")
)

type killSwitchServiceServer struct {
	pb.UnimplementedKillSwitchServiceServer
	dbPool       *sql.DB
	jwtValidator *idtoken.Validator
	cache        ttlcache.SimpleCache
}

func newKillSwitchServiceServer() *killSwitchServiceServer {
	s := &killSwitchServiceServer{}
	db, err := sql.Open("mysql", *dbDsn)
	if err != nil {
		log.Fatalf("unable to open database connection: %v", err)
	}
	db.SetConnMaxLifetime(time.Duration(*dbConnMaxLifetime) * time.Second)
	db.SetMaxOpenConns(*dbMaxOpenConns)
	db.SetMaxIdleConns(*dbMaxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

	s.dbPool = db

	client := &http.Client{}
	validator, err := idtoken.NewValidator(context.Background(), idtoken.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("unable to start idtoken validator: %v", err)
	}
	s.jwtValidator = validator

	s.cache = ttlcache.NewCache()

	return s
}

func getAuthenticatedUser(s *killSwitchServiceServer, ctx context.Context) (*pb.KillSwitchAuthorizedUser, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "getAuthenticatedUser: can't retrieve metadata from incoming request")
	}

	authorization := md.Get("authorization")
	if len(authorization) < 1 {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
	}

	token := authorization[0]
	payload, err := s.jwtValidator.Validate(ctx, token, *jwtAudience)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "getAuthenticatedUser: can't parse idtoken")
	}

	var authorizedUsers []*pb.KillSwitchAuthorizedUser
	var okAssertion bool

	authorizedUsersInterface, err := s.cache.Get("AuthorizedUsers")
	if err == nil {
		authorizedUsers, okAssertion = authorizedUsersInterface.([]*pb.KillSwitchAuthorizedUser)
	}
	if err != nil || !okAssertion {
		freshAuthorizedUsers, err := db.ListAuthorizedUsers(s.dbPool, ctx)
		if err != nil {
			log.Printf("getAuthenticatedUser: error while getting authorized users: %v", err)
			return nil, status.Errorf(codes.Internal, "getAuthenticatedUser: can't get list of authorized users")
		}
		s.cache.SetWithTTL("AuthorizedUsers", freshAuthorizedUsers, authorizedUsersCacheTime)
		authorizedUsers = freshAuthorizedUsers
	}

	// Check if the current user is one of the authorized users, and if so return it.
	for _, u := range authorizedUsers {
		if u.GetGoogleUid() != "" && u.GetGoogleUid() == payload.Subject {
			return u, nil
		}

		email, emailExists := payload.Claims["email"]
		emailVerified, emailVerifiedExists := payload.Claims["email_verified"]

		if u.GetEmail() != "" && emailVerifiedExists && emailExists && emailVerified == true && email == u.GetEmail() {
			return u, nil
		}
	}

	return nil, status.Errorf(codes.PermissionDenied, "User is not part of the authorized users list.")
}

func userHasAccessLevel(requiredLevel pb.KillSwitchAuthorizedUser_AccessLevel, user *pb.KillSwitchAuthorizedUser) bool {
	return requiredLevel <= user.GetAccessLevel()
}

func errorWhenMissingAccess(requiredLevel pb.KillSwitchAuthorizedUser_AccessLevel, user *pb.KillSwitchAuthorizedUser) error {
	if userHasAccessLevel(requiredLevel, user) {
		return nil
	}

	return status.Errorf(codes.PermissionDenied, "User has lower access level than the action that they are trying to perform.")
}

func (s *killSwitchServiceServer) GetKillSwitchStatus(ctx context.Context, req *pb.GetKillSwitchStatusRequest) (*pb.GetKillSwitchStatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "Unimplemented method.")
}

func (s *killSwitchServiceServer) GetKillSwitchOverview(ctx context.Context, req *pb.GetKillSwitchOverviewRequest) (*pb.GetKillSwitchOverviewResponse, error) {
	killSwitches, err := db.ListKillSwitches(s.dbPool, ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}
	res := &pb.GetKillSwitchOverviewResponse{
		KillSwitches: killSwitches,
	}
	return res, nil
}

func (s *killSwitchServiceServer) SyncFeatures(ctx context.Context, req *pb.SyncFeaturesRequest) (*pb.SyncFeaturesResponse, error) {
	// This method requires authentication
	authenticatedUser, err := getAuthenticatedUser(s, ctx)
	if err != nil {
		return nil, err
	}
	err = errorWhenMissingAccess(pb.KillSwitchAuthorizedUser_ACCESS_LEVEL_ADMIN, authenticatedUser)
	if err != nil {
		return nil, err
	}

	log.Println("Syncing features...")

	for _, feature := range req.Features {
		existingFeature, err := db.GetFeatureByCodename(s.dbPool, ctx, feature.Codename)
		if err != nil {
			return nil, status.Errorf(codes.Unavailable, err.Error())
		}
		// If the feature didn't exist in the db, add it. Otherwise, update it if applicable.
		if existingFeature == nil {
			if err := db.AddFeature(s.dbPool, ctx, feature); err != nil {
				return nil, status.Error(codes.Unavailable, err.Error())
			}
		} else {
			canonicalExistingFeature := *existingFeature
			canonicalExistingFeature.Id = 0
			if !proto.Equal(&canonicalExistingFeature, feature) {
				if err := db.UpdateFeature(s.dbPool, ctx, existingFeature.Id, feature); err != nil {
					return nil, status.Error(codes.Unavailable, err.Error())
				}
			}
		}
	}

	res := &pb.SyncFeaturesResponse{}
	return res, nil
}

func (s *killSwitchServiceServer) ListFeatures(ctx context.Context, req *pb.ListFeaturesRequest) (*pb.ListFeaturesResponse, error) {
	features, err := db.ListFeatures(s.dbPool, ctx, req.WithDeprecatedFeatures)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}
	res := &pb.ListFeaturesResponse{
		Features: features,
	}
	return res, nil
}

func (s *killSwitchServiceServer) EnableKillSwitch(ctx context.Context, req *pb.EnableKillSwitchRequest) (*pb.EnableKillSwitchResponse, error) {
	// This method requires authentication
	authenticatedUser, err := getAuthenticatedUser(s, ctx)
	if err != nil {
		return nil, err
	}
	err = errorWhenMissingAccess(pb.KillSwitchAuthorizedUser_ACCESS_LEVEL_ACTIVATOR, authenticatedUser)
	if err != nil {
		return nil, err
	}

	if req.GetKillSwitch().GetFeature().GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "feature.id must be set.")
	}

	err = db.EnableKillSwitch(s.dbPool, ctx, req.KillSwitch, authenticatedUser)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}
	res := &pb.EnableKillSwitchResponse{}
	return res, nil
}

func (s *killSwitchServiceServer) DisableKillSwitch(ctx context.Context, req *pb.DisableKillSwitchRequest) (*pb.DisableKillSwitchResponse, error) {
	// This method requires authentication
	authenticatedUser, err := getAuthenticatedUser(s, ctx)
	if err != nil {
		return nil, err
	}
	err = errorWhenMissingAccess(pb.KillSwitchAuthorizedUser_ACCESS_LEVEL_ACTIVATOR, authenticatedUser)
	if err != nil {
		return nil, err
	}

	if req.GetKillSwitchId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "kill_switch_id must be set.")
	}

	err = db.DisableKillSwitch(s.dbPool, ctx, req.KillSwitchId, authenticatedUser)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}
	res := &pb.DisableKillSwitchResponse{}
	return res, nil
}

func (s *killSwitchServiceServer) ListAuthorizedUsers(ctx context.Context, req *pb.ListAuthorizedUsersRequest) (*pb.ListAuthorizedUsersResponse, error) {
	// This method requires authentication
	authenticatedUser, err := getAuthenticatedUser(s, ctx)
	if err != nil {
		return nil, err
	}
	err = errorWhenMissingAccess(pb.KillSwitchAuthorizedUser_ACCESS_LEVEL_ACTIVATOR, authenticatedUser)
	if err != nil {
		return nil, err
	}

	users, err := db.ListAuthorizedUsers(s.dbPool, ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}
	res := &pb.ListAuthorizedUsersResponse{
		Users: users,
	}
	return res, nil
}

func (s *killSwitchServiceServer) AddAuthorizedUser(ctx context.Context, req *pb.AddAuthorizedUserRequest) (*pb.AddAuthorizedUserResponse, error) {
	// This method requires authentication
	authenticatedUser, err := getAuthenticatedUser(s, ctx)
	if err != nil {
		return nil, err
	}
	err = errorWhenMissingAccess(pb.KillSwitchAuthorizedUser_ACCESS_LEVEL_ADMIN, authenticatedUser)
	if err != nil {
		return nil, err
	}

	if req.GetUser().GetGoogleUid() == "" && req.GetUser().GetEmail() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "At least one of google_uid or email must be set.")
	}

	err = db.AddAuthorizedUser(s.dbPool, ctx, req.User, authenticatedUser)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}
	res := &pb.AddAuthorizedUserResponse{}
	return res, nil
}

func (s *killSwitchServiceServer) UpdateAuthorizedUser(ctx context.Context, req *pb.UpdateAuthorizedUserRequest) (*pb.UpdateAuthorizedUserResponse, error) {
	// This method requires authentication
	authenticatedUser, err := getAuthenticatedUser(s, ctx)
	if err != nil {
		return nil, err
	}
	err = errorWhenMissingAccess(pb.KillSwitchAuthorizedUser_ACCESS_LEVEL_ADMIN, authenticatedUser)
	if err != nil {
		return nil, err
	}

	if req.GetUserId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user_id must be greater than 0.")
	}

	if req.GetUser().GetGoogleUid() == "" && req.GetUser().GetEmail() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "At least one of google_uid or email must be set.")
	}

	err = db.UpdateAuthorizedUser(s.dbPool, ctx, req.UserId, req.User, authenticatedUser)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}
	res := &pb.UpdateAuthorizedUserResponse{}
	return res, nil
}

func (s *killSwitchServiceServer) DeleteAuthorizedUser(ctx context.Context, req *pb.DeleteAuthorizedUserRequest) (*pb.DeleteAuthorizedUserResponse, error) {
	// This method requires authentication
	authenticatedUser, err := getAuthenticatedUser(s, ctx)
	if err != nil {
		return nil, err
	}
	err = errorWhenMissingAccess(pb.KillSwitchAuthorizedUser_ACCESS_LEVEL_ADMIN, authenticatedUser)
	if err != nil {
		return nil, err
	}

	if req.GetUserId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user_id must be greater than 0.")
	}

	err = db.DeleteAuthorizedUser(s.dbPool, ctx, req.UserId, authenticatedUser)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}
	res := &pb.DeleteAuthorizedUserResponse{}
	return res, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterKillSwitchServiceServer(grpcServer, newKillSwitchServiceServer())
	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)
	grpcServer.Serve(lis)
}
