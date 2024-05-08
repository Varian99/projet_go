package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	pb "JulienHugo/Projet_RT0805/file_transfer"

	"google.golang.org/grpc"
)

type Client struct {
	Name string `bson:"name"`
}

type Operation struct {
	Type          string `json:"type"`
	Has_succeeded bool   `json:"has_succeeded"`
}

type Data struct {
	Device_name string      `json:"device_name"`
	Operations  []Operation `json:"operations"`
}

func (c Client) PrettyPrint() {
	fmt.Printf("Client %s démarre \n", c.Name)
}

// lis les fichiers json et decode le contenu json dans une liste de type Data (structure ci-dessus)
func LireFichier(name string) ([]Data, error) {

	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	var data []Data
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		log.Fatal(err)
	}
	return data, nil
}

// On crée les requêtes à envoyer (les strucutres dans le fichiers .proto) en utilisant les données de la liste de type data
func Extract_data(data []Data, num_jour int) *pb.Request {

	request := &pb.Request{}
	request.Day = int32(num_jour)

	for _, d := range data {

		device := &pb.Device{}

		//pour chaque device, on a 3 opérations create, delete et update

		operation_create := &pb.Operation{}
		operation_delete := &pb.Operation{}
		operation_update := &pb.Operation{}

		operation_create.Type = "CREATE"
		operation_delete.Type = "DELETE"
		operation_update.Type = "UPDATE"

		//On traite les données sur le client, pour chaques opérations d'un device, on va envoyé le nombres de réussite et d'échecs

		create_has_succeed := 0
		create_has_not_succeed := 0
		delete_has_succeed := 0
		delete_has_not_succeed := 0
		update_has_succeed := 0
		update_has_not_succeed := 0

		device.DeviceName = d.Device_name

		for _, op := range d.Operations {

			if op.Type == "CREATE" {
				if op.Has_succeeded {
					create_has_succeed++
				} else {
					create_has_not_succeed++
				}
			}
			if op.Type == "DELETE" {
				if op.Has_succeeded {
					delete_has_succeed++
				} else {
					delete_has_not_succeed++
				}
			}
			if op.Type == "UPDATE" {
				if op.Has_succeeded {
					update_has_succeed++
				} else {
					update_has_not_succeed++
				}
			}
		}

		operation_create.HasSucceed = int32(create_has_succeed)
		operation_create.HasNotSucceed = int32(create_has_not_succeed)
		operation_delete.HasSucceed = int32(delete_has_succeed)
		operation_delete.HasNotSucceed = int32(delete_has_not_succeed)
		operation_update.HasSucceed = int32(update_has_succeed)
		operation_update.HasNotSucceed = int32(update_has_not_succeed)

		device.Operations = append(device.Operations, operation_create, operation_delete, operation_update)
		request.Devices = append(request.Devices, device)
	}
	return request
}

// fonction appellé dans le main.go
// démarre la connexion avec le serveur sur le port 8089
// Lis tous les fichiers journee.json dans le dossier donnees
func RunClient() {
	fmt.Println("Démarrage du client")
	conn, err := grpc.Dial(":8089", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial server: %v", err)
	}
	defer conn.Close()

	c := pb.NewFileTransferClient(conn)

	for i := 1; ; i++ {
		nomFichier := fmt.Sprintf("donnees/journee_%d.json", i)
		_, err := os.Stat(nomFichier)

		//quand le fichier n'existe pas alors cela veut dire qu'on a lu tous les fichiers disponibles
		//on sort donc de la boucle et on arrête la communication avec le serveur
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("Le fichier %s n'existe pas. Sortie de la boucle.\n", nomFichier)
				break
			}
		}

		data, err := LireFichier(nomFichier)
		if err != nil {
			log.Fatalf("failed to dial server: %v", err)
		}

		request := Extract_data(data, i)
		response, err := c.Create(context.Background(), request)
		if err != nil {
			log.Fatalf("failed to call CreateClient RPC method: %v", err)
		}
		log.Printf("Response from server: %v", response)
	}

	log.Printf("Fin de l'envoie")
}
