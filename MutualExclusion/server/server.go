package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mAssignment2/MutualExclusion/protobuf"
	"net"
	"sync"

	"google.golang.org/grpc"
)

type server struct {
	protobuf.UnimplementedMutualExclusionServer
}

var clients []string
var queue []string
var criticalValue int
var mu sync.Mutex

func main() {
	criticalValue = 0
	lis, err := net.Listen("tcp", ":8080")

	if err != nil { //error before listening
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer() //we create a new server
	protobuf.RegisterMutualExclusionServer(s, &server{})

	go serverManipulationRoutine()
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil { //error while listening
		log.Fatalf("failed to serve: %v", err)
	}
}

func serverManipulationRoutine() {
	for {
		clientName := takeNext()
		if clientName != "" {
			fmt.Println("Manipulating for " + clientName)
			//If anyone is waiting, take the next from the queue
			manipulateCriticalValue()
		}
	}
}

func (s *server) NewParticipant(ctx context.Context, in *protobuf.NewClientRequest) (*protobuf.NewClientReply, error) {
	if alreadyExists(in.ClientName) {
		fmt.Println("Client denied: " + in.ClientName)
		return &protobuf.NewClientReply{}, errors.New("USERNAME IS ALREADY IN USE")
	} else {
		fmt.Println("New Client with name: " + in.ClientName)
		clients = append(clients, in.ClientName)
		return &protobuf.NewClientReply{}, nil
	}
}

func (s *server) ClientManipulation(ctx context.Context, in *protobuf.ClientManipulationRequest) (*protobuf.ClientManipulationReply, error) {
	queue = append(queue, in.ClientName)
	fmt.Print("[")
	for i := 0; i < len(queue); i++ {
		fmt.Print(queue[i] + ", ")
	}
	fmt.Print("]")
	fmt.Println()
	return &protobuf.ClientManipulationReply{}, nil
}

func alreadyExists(clientName string) bool {
	var existsInClients = false
	for i := 0; i < len(clients); i++ {
		fmt.Println(clients[i] == clientName)
		if clients[i] == clientName {
			existsInClients = true
		}
	}
	return existsInClients
}

func takeNext() string {
	for len(queue) > 0 {
		value := queue[0]
		queue = queue[1:]
		return value
	}
	return ""
}

func manipulateCriticalValue() {
	mu.Lock()
	criticalValue++
	mu.Unlock()
	fmt.Println(criticalValue)
}

/*
//Client: i want to participate
//Server: C1 is noted

//Client: (c1) i want to manipulate
c1.wantsToManipulate();
server.push(c1);

func takeNext(){
	if(queue.length == 0){
		//Wait for next queue?
	}
	else{
		server.dequeue().runManipulation();
		//Server: Ok i made your manipulation, thanks. I will take the next one
	}
}


manipulation() {
	mutex.lock();
	critical = [...];
	mutex.unlock();
}

takeNext(); //server knows that client is done, who is next?

*/
