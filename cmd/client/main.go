package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	pb "mail-srvc/pkg/api"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func Send(conn *grpc.ClientConn, id string, email string) {
	client := pb.NewMailServiceClient(conn)
	request := &pb.CreatedUser{Id: id, Email: email}
	response, err := client.SendEmail(context.Background(), request)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	log.Println(response)
}

func Verify(conn *grpc.ClientConn, id string, token string) {
	client := pb.NewMailServiceClient(conn)
	request := &pb.ConfirmUserRequest{Id: id, Token: strings.TrimSuffix(token, "\n")}
	response, err := client.VerifyEmail(context.Background(), request)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	log.Println(response.Confirmed)
}

func main() {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial("0.0.0.0:8080", opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for i := 0; i < 1; i++ {
		Send(conn, "1", "vlad.homam@gmail.com")
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter token: ")
	token, _ := reader.ReadString('\n')

	Verify(conn, "1", token)

}
