package logic

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	uuid "github.com/satori/go.uuid"

	pb "mail-srvc/pkg/api"
	rp "mail-srvc/repo"

	mail "gopkg.in/mail.v2"

	"github.com/golang/protobuf/ptypes/empty"
)

//CONSTANTS SHULD BE IN OS ENV
const NUMBER_OF_MAILS_SEND_AT_THE_SAME_TIME = 4

type MailServer struct {
	repo       rp.Repository
	mailDialer mail.Dialer
	queue      chan MailTask
	closeChan  chan os.Signal
}

func (m *MailServer) taskHandler() {
	sem := make(chan bool, NUMBER_OF_MAILS_SEND_AT_THE_SAME_TIME)
	var wg sync.WaitGroup
	for {
		select {
		case task := <-m.queue:
			sem <- true
			wg.Add(1)
			go func() {
				defer wg.Done()
				task.Send(&m.mailDialer)
				<-sem
			}()
		case <-m.closeChan:
			log.Println("\r MailTaskHendler stoped")
			wg.Wait()
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
	// ifExist := m.repo.VerifyIfExist(ctx, &pb.ConfirmUserRequest{Id: req.GetId(), Token: token})

	// if !ifExist {
	// 	return nil, errors.New("Email never has sent")
	// }

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

	m.queue <- MailTask{from: "vlad.homam@gmail.com", to: req.GetEmail(), content: fmt.Sprintf("Hello User, this is your uuid: %s", token)}

	return nil
}

func (m *MailServer) VerifyEmail(ctx context.Context, req *pb.ConfirmUserRequest) (*pb.ConfirmUserResponse, error) {
	mailResp := m.repo.VerifyIfExist(ctx, req)

	return &pb.ConfirmUserResponse{Confirmed: mailResp}, nil
}
