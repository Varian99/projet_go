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

// Cette fonction est appellé quand le serveur reçoit une requête de la part du client
// il convertit la requête au format bson en utilisant la structure mongoDBData
// et utilise la fonction StoreToMongoDB pour stocker les données dans mongodb
func (s myFileTransferServer) Create(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	log.Printf("requete recu: %v", req)
	log.Printf("Stockage dans mangoDB")
	day := req.Day
	for _, device := range req.Devices {
		s.StoreToMongoDB(device, int(day))
	}

	return &pb.Response{}, nil
}

// fonction qui stocke les données data dans la base de donnée Donnee_Projet et dans la collections donnee_journee corespondante
// On crée une collection pour chaque journée
func (s myFileTransferServer) StoreToMongoDB(data *pb.Device, day int) error {
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

// fonction appellé dans le main.go
// démarre la connexion avec le client sur le port 8089
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
