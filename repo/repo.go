package repo

import (
	"context"

	pb "mail-srvc/pkg/api"
)

type Repository interface {
	SaveEmailVerification(context.Context, *pb.CreatedUser, string) error
	//VirifyAccount(context.Context, *pb.CreatedUser, string) (*pb.CreatedUserResponse, error)
	VerifyIfExist(context.Context, *pb.ConfirmUserRequest) bool
}
