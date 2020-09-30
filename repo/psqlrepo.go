package repo

import (
	"context"
	"database/sql"
	pb "mail-srvc/pkg/api"
)

//PsqlRepository - repo for psql database
type PsqlRepository struct {
	DB *sql.DB
}

//NewPsqlREpository - create new Psql repo
func NewPsqlREpository(database *sql.DB) Repository {
	return &PsqlRepository{DB: database}
}

// func (r *PsqlRepository) VirifyAccount(ctx context.Context, req *pb.CreatedUser, token string) (*pb.CreatedUserResponse, error) {
// 	return nil, nil
// }

func (r *PsqlRepository) SaveEmailVerification(ctx context.Context, req *pb.CreatedUser, token string) error {
	return nil
}

func (r *PsqlRepository) VerifyIfExist(ctx context.Context, req *pb.ConfirmUserRequest) bool {
	return false
}
