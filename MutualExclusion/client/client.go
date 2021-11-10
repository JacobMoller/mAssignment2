package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"mAssignment2/MutualExclusion/protobuf"

	"google.golang.org/grpc"
)

func main() {
	log.Print("Welcome Client. You need to provide a name for the server to remember you:")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	name := strings.Replace(text, "\n", "", 1)

	conn, err := grpc.Dial(":8080", grpc.WithInsecure(), grpc.WithBlock()) //maybe it has to be: localhost:8080
	if err != nil {                                                        //error can not establish connection
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := protobuf.NewMutualExclusionClient(conn)
	message, err2 := client.NewParticipant(context.Background(), &protobuf.NewClientRequest{ClientName: name})
	if err2 != nil {
		//Error handling
		if message == nil {
			fmt.Println("Username is already in use")
		}
	} else {
		for {
			client.ClientManipulation(context.Background(), &protobuf.ClientManipulationRequest{ClientName: name})
			time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		}
	}
}
