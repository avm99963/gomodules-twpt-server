package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	pb "gomodules.avm99963.com/twpt-server/api_proto"
	db "gomodules.avm99963.com/twpt-server/internal/db"
)

var (
	port              = flag.Int("port", 10000, "The server port")
	dbDsn             = flag.String("db", "", "MySQL/MariaDB database data source name (DSN)") // https://github.com/go-sql-driver/mysql#dsn-data-source-name
	dbConnMaxLifetime = flag.Int("dbConnMaxLifetime", 3*60, "Maximum amount of time a connection to the database may be reused in seconds.")
	dbMaxOpenConns    = flag.Int("dbMaxOpenConns", 5, "Maximum number of open connections to the database.")
	dbMaxIdleConns    = flag.Int("dbMaxIdleConns", *dbMaxOpenConns, "Maximum number of connections to the database in the idle connection pool.")
)

type killSwitchServiceServer struct {
	pb.UnimplementedKillSwitchServiceServer
	dbPool *sql.DB
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
	return s
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
	if req.GetKillSwitch().GetFeature().GetId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "feature.id must be set.")
	}

	err := db.EnableKillSwitch(s.dbPool, ctx, req.KillSwitch)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}
	res := &pb.EnableKillSwitchResponse{}
	return res, nil
}

func (s *killSwitchServiceServer) DisableKillSwitch(ctx context.Context, req *pb.DisableKillSwitchRequest) (*pb.DisableKillSwitchResponse, error) {
	if req.GetKillSwitchId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "kill_switch_id must be set.")
	}

	err := db.DisableKillSwitch(s.dbPool, ctx, req.KillSwitchId)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}
	res := &pb.DisableKillSwitchResponse{}
	return res, nil
}

func (s *killSwitchServiceServer) ListAuthorizedUsers(ctx context.Context, req *pb.ListAuthorizedUsersRequest) (*pb.ListAuthorizedUsersResponse, error) {
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
	if req.GetUser().GetGoogleUid() == "" && req.GetUser().GetEmail() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "At least one of google_uid or email must be set.")
	}

	err := db.AddAuthorizedUser(s.dbPool, ctx, req.User)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}
	res := &pb.AddAuthorizedUserResponse{}
	return res, nil
}

func (s *killSwitchServiceServer) UpdateAuthorizedUser(ctx context.Context, req *pb.UpdateAuthorizedUserRequest) (*pb.UpdateAuthorizedUserResponse, error) {
	if req.GetUserId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user_id must be greater than 0.")
	}

	if req.GetUser().GetGoogleUid() == "" && req.GetUser().GetEmail() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "At least one of google_uid or email must be set.")
	}

	err := db.UpdateAuthorizedUser(s.dbPool, ctx, req.UserId, req.User)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, err.Error())
	}
	res := &pb.UpdateAuthorizedUserResponse{}
	return res, nil
}

func (s *killSwitchServiceServer) DeleteAuthorizedUser(ctx context.Context, req *pb.DeleteAuthorizedUserRequest) (*pb.DeleteAuthorizedUserResponse, error) {
	if req.GetUserId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user_id must be greater than 0.")
	}

	err := db.DeleteAuthorizedUser(s.dbPool, ctx, req.UserId)
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
