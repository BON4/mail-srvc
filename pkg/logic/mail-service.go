package logic

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	pb "mail-srvc/pkg/api"
	rp "mail-srvc/repo"

	mail "gopkg.in/mail.v2"

	"github.com/golang/protobuf/ptypes/empty"
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
			//TODO handle gorutines errors
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

func (m *MailServer) SendEmail(ctx context.Context, req *pb.CreatedUser) (*empty.Empty, error) {
	token := fmt.Sprintf("%d", rand.Int31())

	//Check if id and token not in db
	ifExist := m.repo.VerifyIfExist(ctx, &pb.ConfirmUserRequest{Id: req.GetId(), Token: token})

	if ifExist {
		return nil, errors.New("Email has been already sent")
	}

	err := m.repo.SaveEmailVerification(ctx, req, token)

	if err != nil {
		return nil, err
	}

	m.queue <- MailTask{from: "Vlad", to: req.GetEmail(), content: fmt.Sprintf("Hello User, this is your uuid: %s", token)}

	return &empty.Empty{}, nil
}

func (m *MailServer) VerifyEmail(ctx context.Context, req *pb.ConfirmUserRequest) (*pb.ConfirmUserResponse, error) {
	mailResp := m.repo.VerifyIfExist(ctx, req)

	return &pb.ConfirmUserResponse{Confirmed: mailResp}, nil
}
