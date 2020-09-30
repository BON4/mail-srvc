package logic

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	pb "mail-srvc/pkg/api"
	rp "mail-srvc/repo"

	mail "gopkg.in/mail.v2"
)

type MailServer struct {
	repo       rp.Repository
	mailDialer mail.Dialer
	queue      chan MailTask
	closeChan  chan os.Signal
}

func (m *MailServer) taskHandler() {
	for {
		select {
		case task := <-m.queue:
			go task.Send(&m.mailDialer)
		case <-m.closeChan:
			fmt.Println("\r MailTaskHendler stoped")
			signal.Stop(m.closeChan)
			syscall.Kill(syscall.Getpid(), syscall.SIGINT)
			close(m.queue)
			return
		}
	}
}

func NewMailServer(repo rp.Repository, dialer *mail.Dialer) pb.MailServiceServer {
	closeChan := make(chan os.Signal, 1)
	signal.Notify(closeChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	mailServer := MailServer{repo: repo, mailDialer: *dialer, queue: make(chan MailTask), closeChan: closeChan}

	go mailServer.taskHandler()

	return &mailServer
}

func (m *MailServer) VirifyAccount(ctx context.Context, req *pb.CreatedUser) (*pb.CreatedUserResponse, error) {
	testUUID := fmt.Sprintf("%d", rand.Int31())
	m.queue <- MailTask{from: "Vlad", to: req.GetEmail(), content: fmt.Sprintf("Hello User, this is your uuid: %s", testUUID)}
	return &pb.CreatedUserResponse{VerifyID: testUUID}, nil
}
