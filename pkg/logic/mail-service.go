package logic

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	uuid "github.com/satori/go.uuid"

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

func (m *MailServer) SendEmailOnce(ctx context.Context, req *pb.CreatedUser) (*empty.Empty, error) {
	token := fmt.Sprintf("%s", uuid.NewV4())

	//Check if id and token not in db
	ifExist := m.repo.VerifyIfExist(ctx, &pb.ConfirmUserRequest{Id: req.GetId(), Token: token})

	if ifExist {
		return nil, errors.New("Email already has sent")
	}

	if err := m.SendEmail(ctx, req, token); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (m *MailServer) ResendEmail(ctx context.Context, req *pb.CreatedUser) (*empty.Empty, error) {
	token := fmt.Sprintf("%s", uuid.NewV4())

	//Check if id and token not in db
	//TODO resend email checking token too, but this method creates new token
	ifExist := m.repo.VerifyIfExist(ctx, &pb.ConfirmUserRequest{Id: req.GetId(), Token: token})

	if !ifExist {
		return nil, errors.New("Email never has sent")
	}

	if err := m.SendEmail(ctx, req, token); err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (m *MailServer) SendEmail(ctx context.Context, req *pb.CreatedUser, token string) error {
	err := m.repo.SaveEmailVerification(ctx, req, token)

	if err != nil {
		return err
	}

	m.queue <- MailTask{from: "Vlad", to: req.GetEmail(), content: fmt.Sprintf("Hello User, this is your uuid: %s", token)}

	return nil
}

func (m *MailServer) VerifyEmail(ctx context.Context, req *pb.ConfirmUserRequest) (*pb.ConfirmUserResponse, error) {
	mailResp := m.repo.VerifyIfExist(ctx, req)

	return &pb.ConfirmUserResponse{Confirmed: mailResp}, nil
}
