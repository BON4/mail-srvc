package main

import (
	"context"
	"log"
	pb "mail-srvc/pkg/api"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func Verify(conn *grpc.ClientConn, id string, email string) {
	client := pb.NewMailServiceClient(conn)
	request := &pb.CreatedUser{Id: id, Email: email}
	response, err := client.VirifyAccount(context.Background(), request)

	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	log.Println(response)
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

	start := time.Now()

	for i := 0; i < 5; i++ {
		Verify(conn, "1", "vlad.homam@gmail.com")
	}

	log.Println(time.Since(start))

	// c := make(chan string)

	// closeChan := make(chan os.Signal)
	// signal.Notify(closeChan, os.Interrupt, syscall.SIGTERM)

	// go testFunc(c, closeChan)

	// for i := 0; i < 10000000; i++ {
	// 	time.Sleep(time.Millisecond * 2)
	// 	select {
	// 	case <-closeChan:
	// 		return
	// 	default:
	// 		c <- fmt.Sprintf("Hello %d", i)
	// 	}
	// }
}
