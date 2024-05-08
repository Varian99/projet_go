package main

import (
	pb "JulienHugo/Projet_RT0805/client"
	pc "JulienHugo/Projet_RT0805/file_transfer"
	"log"
	"reflect"
	"testing"
)

// Test pour vérifier que la lecture et le renvoi des données d'un fichier se passe correctement
// On utilise un fichier préparé test.json et on compare le retour avec celui attendu
func TestLireFichier(t *testing.T) {
	result, err := pb.LireFichier("test.json")

	if err != nil {
		t.Errorf("Erreur lors de la lecture du fichier : %v", err)
	}
	expected := []pb.Data{
		{
			Device_name: "c1153f7a-b060-4215-bf22-601e8f8e704c",
			Operations: []pb.Operation{
				{Type: "DELETE", Has_succeeded: false},
				{Type: "CREATE", Has_succeeded: true},
				{Type: "DELETE", Has_succeeded: true},
				{Type: "CREATE", Has_succeeded: true},
				{Type: "DELETE", Has_succeeded: true},
				{Type: "UPDATE", Has_succeeded: false},
				{Type: "CREATE", Has_succeeded: true},
				{Type: "UPDATE", Has_succeeded: true},
				{Type: "DELETE", Has_succeeded: true},
				{Type: "UPDATE", Has_succeeded: true},
				{Type: "DELETE", Has_succeeded: false},
			},
		},
		{
			Device_name: "971645e6-6870-4db6-9e6b-817227d8f338",
			Operations: []pb.Operation{
				{Type: "DELETE", Has_succeeded: false},
				{Type: "CREATE", Has_succeeded: true},
				{Type: "DELETE", Has_succeeded: true},
				{Type: "CREATE", Has_succeeded: true},
				{Type: "DELETE", Has_succeeded: true},
				{Type: "UPDATE", Has_succeeded: false},
				{Type: "CREATE", Has_succeeded: true},
			},
		},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Les résultats ne correspondent pas. \nAttendu : %v \nObtenu : %v", expected, result)
	}
}

// Test pour vérifier que les données ont bien été formatées selon le bon format pour pouvoir être envoyées
// Ce format est celui des structures contenues dans le .proto
// On utilise le même fichier de test que pour le premier test, si il fonctionne bien on peut tester cette fonction avec

func TestExctractData(t *testing.T) {
	result, err := pb.LireFichier("test.json")
	if err != nil {
		log.Fatalf("Erreur lors de la lecture du fichier : %v", err)
	}

	request := pb.Extract_data(result, 1)
	expected := &pc.Request{}
	expected.Day = 1

	device1 := &pc.Device{}
	operation_create := &pc.Operation{}
	operation_delete := &pc.Operation{}
	operation_update := &pc.Operation{}
	operation_create.Type = "CREATE"
	operation_delete.Type = "DELETE"
	operation_update.Type = "UPDATE"
	operation_create.HasSucceed = 3
	operation_create.HasNotSucceed = 0
	operation_delete.HasSucceed = 3
	operation_delete.HasNotSucceed = 2
	operation_update.HasSucceed = 2
	operation_update.HasNotSucceed = 1
	device1.Operations = append(device1.Operations, operation_create, operation_delete, operation_update)
	device1.DeviceName = "c1153f7a-b060-4215-bf22-601e8f8e704c"
	expected.Devices = append(expected.Devices, device1)

	device2 := &pc.Device{}
	operation_create2 := &pc.Operation{}
	operation_delete2 := &pc.Operation{}
	operation_update2 := &pc.Operation{}
	operation_create2.Type = "CREATE"
	operation_delete2.Type = "DELETE"
	operation_update2.Type = "UPDATE"
	operation_create2.HasSucceed = 3
	operation_create2.HasNotSucceed = 0
	operation_delete2.HasSucceed = 2
	operation_delete2.HasNotSucceed = 1
	operation_update2.HasSucceed = 0
	operation_update2.HasNotSucceed = 1
	device2.Operations = append(device2.Operations, operation_create2, operation_delete2, operation_update2)
	device2.DeviceName = "971645e6-6870-4db6-9e6b-817227d8f338"
	expected.Devices = append(expected.Devices, device2)

	// On voulait le faire de cette façon mais il y avait un problème de pointeurs
	// device1 := &pc.Device{
	//     DeviceName: "c1153f7a-b060-4215-bf22-601e8f8e704c",
	//     Operations: []*pc.Operation{
	//         {Type: "DELETE", HasSucceed: 3, HasNotSucceed: 2},
	//         {Type: "CREATE", HasSucceed: 3, HasNotSucceed: 0},
	//         {Type: "UPDATE", HasSucceed: 2, HasNotSucceed: 1},
	//     },
	// }
	// devices = append(devices, device1)
	// device2 := pc.Device{
	//     DeviceName: "971645e6-6870-4db6-9e6b-817227d8f338",
	//     Operations: []*pc.Operation{
	//         {Type: "DELETE", HasSucceed: 2, HasNotSucceed: 1},
	//         {Type: "CREATE", HasSucceed: 3, HasNotSucceed: 0},
	//         {Type: "UPDATE", HasSucceed: 0, HasNotSucceed: 1},
	//     },
	// }
	// devices = append(devices, device2)

	// if !reflect.DeepEqual(request, expected) {
	//     t.Errorf("\nLes résultats ne correspondent pas. \nAttendu :\n%v \nObtenu :\n%v", expected, request)
	// }
	if len(request.Devices) != len(expected.Devices) {
		t.Errorf("Les longueurs des slices devices ne correspondent pas. Attendu : %d, Obtenu : %d", len(expected.Devices), len(request.Devices))
	} else {
		for i := range expected.Devices {
			if request.Devices[i].DeviceName != expected.Devices[i].DeviceName {
				t.Errorf("Le champ device_name pour le device %d ne correspond pas. Attendu : %s, Obtenu : %s", i, expected.Devices[i].DeviceName, request.Devices[i].DeviceName)
			}
			if len(request.Devices[i].Operations) != len(expected.Devices[i].Operations) {
				t.Errorf("Les longueurs des slices operations pour le device %d ne correspondent pas. Attendu : %d, Obtenu : %d", i, len(expected.Devices[i].Operations), len(request.Devices[i].Operations))
			} else {
				for j := range expected.Devices[i].Operations {
					if request.Devices[i].Operations[j].Type != expected.Devices[i].Operations[j].Type {
						t.Errorf("Le champ type pour l'opération %d du device %d ne correspond pas. Attendu : %s, Obtenu : %s", j, i, expected.Devices[i].Operations[j].Type, request.Devices[i].Operations[j].Type)
					}
					if request.Devices[i].Operations[j].HasSucceed != expected.Devices[i].Operations[j].HasSucceed {
						t.Errorf("Le champ has_succeed pour l'opération %d du device %d ne correspond pas. Attendu : %d, Obtenu : %d", j, i, expected.Devices[i].Operations[j].HasSucceed, request.Devices[i].Operations[j].HasSucceed)
					}
					if request.Devices[i].Operations[j].HasNotSucceed != expected.Devices[i].Operations[j].HasNotSucceed {
						t.Errorf("Le champ has_not_succeed pour l'opération %d du device %d ne correspond pas. Attendu : %d, Obtenu : %d", j, i, expected.Devices[i].Operations[j].HasNotSucceed, request.Devices[i].Operations[j].HasNotSucceed)
					}
				}
			}
		}
	}
}
