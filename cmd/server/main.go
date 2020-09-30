package main

import (
	"crypto/tls"
	"database/sql"
	"log"
	"mail-srvc/pkg/logic"
	"net"

	rp "mail-srvc/repo"

	pb "mail-srvc/pkg/api"

	_ "github.com/lib/pq"

	"google.golang.org/grpc"
	mail "gopkg.in/mail.v2"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	addr := "0.0.0.0:8080"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
		//panic(err)
	}

	db, err := sql.Open("postgres", "host=psqlhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")

	if err != nil {
		log.Fatal(err)
		//panic(err)
	}

	repo := rp.NewPsqlREpository(db)

	dailer := mail.NewDialer("smtp.gmail.com", 587, "email", "password")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	dailer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	serverOptions := []grpc.ServerOption{}
	s := grpc.NewServer(serverOptions...)

	pb.RegisterMailServiceServer(s, logic.NewMailServer(repo, dailer))

	log.Println("Serving gRPC on https://", addr)
	log.Fatal(s.Serve(lis))
}
