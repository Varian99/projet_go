package main

import (
	pb "JulienHugo/Projet_RT0805/client"
	"fmt"
	"sync"
)

// obliger d'utiliser les go routines pour que le main execute le serveur et le client
// mais il faut utiliser sync aussi sinon lke main n'attends pas et se termine

func main() {
	var wg sync.WaitGroup
	wg.Add(2) //ajoute deux go routines
	go func() {
		defer wg.Done()
		pb.RunServer()
	}()
	go func() {
		defer wg.Done()
		pb.RunClient()
	}()
	wg.Wait()
}

