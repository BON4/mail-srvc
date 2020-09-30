package repo

import (
	"context"

	pb "mail-srvc/pkg/api"
)

type Repository interface {
	VirifyAccount(context.Context, *pb.CreatedUser) (*pb.CreatedUserResponse, error)
}
