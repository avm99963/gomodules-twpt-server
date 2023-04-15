package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/johnsiilver/getcert"

	pb "gomodules.avm99963.com/twpt-server/api_proto"
)

var (
	grpcEndpoint = flag.String("grpcEndpoint", "", "gRPC endpoint address.")
	jwt          = flag.String("jwt", "", "JWT credentials.")
	insecure     = flag.Bool("insecure", false, "Set if the connection to the gRPC endpoint is insecure.")
)

type Features map[string]Feature

type Feature struct {
	DefaultValue   interface{} `json:"defaultValue"`
	Context        string      `json:"context"`
	KillSwitchType string      `json:"killSwitchType"`
}

func convertStringTypeToPb(context string) pb.Feature_Type {
	switch context {
	case "option":
		return pb.Feature_TYPE_OPTION

	case "experiment":
		return pb.Feature_TYPE_EXPERIMENT

	case "internalKillSwitch":
		return pb.Feature_TYPE_INTERNAL_KILL_SWITCH

	case "deprecated":
		return pb.Feature_TYPE_DEPRECATED

	default:
		return pb.Feature_TYPE_UNKNOWN
	}
}

func main() {
	flag.Parse()

	var err error
	var conn *grpc.ClientConn

	if *insecure {
		conn, err = grpc.Dial(*grpcEndpoint, grpc.WithInsecure())
	} else {
		tlsCert, _, err2 := getcert.FromTLSServer(*grpcEndpoint, false)
		if err2 != nil {
			log.Fatalf("error while retrieving public certificate: %v\n", err2)
		}
		conn, err = grpc.Dial(*grpcEndpoint, grpc.WithTransportCredentials(credentials.NewServerTLSFromCert(&tlsCert)))
	}
	if err != nil {
		log.Fatalf("error while connecting to gRPC endpoint: %v\n", err)
	}
	defer conn.Close()

	client := pb.NewKillSwitchServiceClient(conn)
	md := metadata.Pairs("authorization", *jwt)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	var jsonFeatures Features
	err = json.NewDecoder(os.Stdin).Decode(&jsonFeatures)
	if err != nil {
		log.Fatalf("can't decode JSON file: %v\n", err)
	}

	features := make([]*pb.Feature, 0)
	for codename, jsonFeature := range jsonFeatures {
		if jsonFeature.Context == "internal" && jsonFeature.KillSwitchType != "internalKillSwitch" {
			continue
		}

		feature := &pb.Feature{
			Codename: codename,
			Type:     convertStringTypeToPb(jsonFeature.KillSwitchType),
		}
		features = append(features, feature)
	}

	request := &pb.SyncFeaturesRequest{
		Features: features,
	}
	_, err = client.SyncFeatures(ctx, request)
	if err != nil {
		log.Fatalf("error syncing features: %v\n", err)
	}

	log.Println("Synced features successfully.")
}
