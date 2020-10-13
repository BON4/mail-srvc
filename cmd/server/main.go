package main

import (
	"context"
	"crypto/tls"
	"log"
	"mail-srvc/pkg/logic"
	"net"
	"os"
	"time"

	"github.com/go-redis/redis/v8"

	rp "mail-srvc/repo"

	pb "mail-srvc/pkg/api"

	_ "github.com/lib/pq"

	"google.golang.org/grpc"
	mail "gopkg.in/mail.v2"
)

//CONSTANTS SHULD BE IN OS ENV
const EMAIL_CONFIRMATION_DURATION = 10

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	addr := "0.0.0.0:8080"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
		//panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		//Change to os.Getenv
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	status := rdb.Ping(context.Background())
	if status.Err() != nil {
		log.Fatal(status.Err(), "Fail to dial to redis")
	}

	defer rdb.Close()

	redisRepo := rp.NewRedisRepository(rdb, time.Duration(EMAIL_CONFIRMATION_DURATION*time.Second))

	log.Println(os.Getenv("MAIL_EMAIL"), os.Getenv("MAIL_PASSWORD"))

	dailer := mail.NewDialer("smtp.gmail.com", 587, os.Getenv("MAIL_EMAIL"), os.Getenv("MAIL_PASSWORD"))

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	dailer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	serverOptions := []grpc.ServerOption{}
	s := grpc.NewServer(serverOptions...)

	pb.RegisterMailServiceServer(s, logic.NewMailServer(redisRepo, dailer))

	log.Println("Serving gRPC on https://", addr)
	log.Fatal(s.Serve(lis))
}
