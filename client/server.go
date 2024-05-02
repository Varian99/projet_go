package client

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "JulienHugo/Projet_RT0805/file_transfer"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

type myFileTransferServer struct {
	pb.UnimplementedFileTransferServer
}

type MongoDBData struct {
	DeviceName string          `bson:"device_name"`
	Operations []OperationData `bson:"operations"`
}

type OperationData struct {
	Type          string `bson:"type"`
	HasSucceed    int    `bson:"has_succeed"`
	HasNotSucceed int    `bson:"has_not_succeed"`
}

func (s myFileTransferServer) Create(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	log.Printf("requete recu: %v", req)
	log.Printf("Stockage dans mangoDB")
	day := req.Day
	for _, device := range req.Devices {
		data := MongoDBData{
			DeviceName: device.DeviceName,
		}
		for _, op := range device.Operations {
			data.Operations = append(data.Operations, OperationData{
				Type:          op.Type,
				HasSucceed:    int(op.HasSucceed),
				HasNotSucceed: int(op.HasNotSucceed),
			})
		}
		s.StoreToMongoDB(data, int(day))
	}

	return &pb.Response{}, nil
}

func (s myFileTransferServer) StoreToMongoDB(data MongoDBData, day int) error {
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI("mongodb://root:root@localhost:27017/"))
	if err != nil {
		return err
	}
	collectionName := fmt.Sprintf("donnee_journee_%d", day)

	coll := client.Database("Donnee_Projet").Collection(collectionName)

	_, err = coll.InsertOne(context.Background(), data, nil)
	if err != nil {
		return err
	}

	return nil
}

func RunServer() {
	lis, err := net.Listen("tcp", ":8089")
	if err != nil {
		log.Fatal("cannot create listener : %v", err)
	}

	serverRegistar := grpc.NewServer()
	service := &myFileTransferServer{}

	pb.RegisterFileTransferServer(serverRegistar, service)
	serverRegistar.Serve(lis)
	if err != nil {
		log.Fatal("Impossible to serve : %v", err)
	}
}
