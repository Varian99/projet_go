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

func LireFichier(name string) ([]Data, error) {
	// fmt.Printf("ouverture du fichier %s \n", name)
	file, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	// Décodez le contenu JSON dans une structure de données
	var data []Data
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		log.Fatal(err)
	}
	return data, nil
}

func Extract_data(data []Data,num_jour int) (*pb.Request){

	request := &pb.Request{}
	request.Day = int32(num_jour)

	for _, d := range data {

		device := &pb.Device{}

		operation_create := &pb.Operation{}
		operation_delete := &pb.Operation{}
		operation_update := &pb.Operation{}

		operation_create.Type = "CREATE"
		operation_delete.Type = "DELETE"
		operation_update.Type = "UPDATE"

		create_has_succeed := 0
		create_has_not_succeed := 0
		delete_has_succeed := 0
		delete_has_not_succeed := 0
		update_has_succeed := 0
		update_has_not_succeed := 0
		//fmt.Printf("Nom du périphérique : %s\n", d.Device_name)
		device.DeviceName = d.Device_name
		// fmt.Printf("Verif client nom du device %v \n", d.Device_name)
		//fmt.Println("Opérations :")
		for _, op := range d.Operations {
			// fmt.Printf("Verif du client du nom operation %v \n", op.Type)
			// fmt.Printf("Verif client du has suceed %v \n", op.Has_succeeded)
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

			//fmt.Printf("  - Type : %s, Réussi : %t\n", op.Type, op.Has_succeeded)
			//operation := &pb.Operation{Type: op.Type, HasSucceed: op.Has_succeeded}
			//operations = append(operations, operation)

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
		
		request := Extract_data(data,i)
		response, err := c.Create(context.Background(), request)
		if err != nil {
			log.Fatalf("failed to call CreateClient RPC method: %v", err)
		}
		log.Printf("Response from server: %v", response)
	}

	log.Printf("Fin de l'envoie")
}
