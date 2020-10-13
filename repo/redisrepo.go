package repo

import (
	"context"
	"errors"
	"fmt"
	"log"
	pb "mail-srvc/pkg/api"
	"time"

	redis "github.com/go-redis/redis/v8"
)

type RedisRepository struct {
	conn                      *redis.Client
	emailConfirmationDuration time.Duration
}

func NewRedisRepository(cache *redis.Client, dur time.Duration) Repository {
	if dur < time.Second {
		log.Fatal("Too short email duration")
	}
	return &RedisRepository{conn: cache, emailConfirmationDuration: dur}
}

func (r *RedisRepository) SaveEmailVerification(ctx context.Context, req *pb.CreatedUser, token string) error {
	//hset mailusers:2
	//expire mailusers:2 6000

	mailKey := fmt.Sprintf("%s:%s", "mail", req.GetId())
	mailVal := map[string]interface{}{"id": req.GetId(), "t": token}
	intCmd := r.conn.HSet(ctx, mailKey, mailVal)

	if _, err := intCmd.Result(); err != nil {
		if err != nil {
			return err
		}
		//return errors.New("Can not add email in redis db, email must already has sent")
	}

	boolCmd := r.conn.Expire(ctx, mailKey, r.emailConfirmationDuration)

	if i, err := boolCmd.Result(); err != nil || i == false {
		return errors.New("Can not set email time in redis db")
	}

	return nil
}

func (r *RedisRepository) VerifyIfExist(ctx context.Context, req *pb.ConfirmUserRequest) bool {
	//hmget mailusers:1
	mailKey := fmt.Sprintf("%s:%s", "mail", req.GetId())
	sliceRedisResp := r.conn.HMGet(ctx, mailKey, "id", "t")

	if sliceRedisResp.Err() != nil {
		return false
	}

	i, err := sliceRedisResp.Result()

	if err != nil || len(i) == 0 {
		return false
	}

	if i[0] == nil || i[1] == nil {
		return false
	}

	// if i[0].(string) != req.GetId() || i[1].(string) != req.GetToken() {
	// 	return false
	// }

	fmt.Println(i[0].(string) == req.GetId())

	if i[0].(string) == req.GetId() && i[1].(string) != req.GetToken() {
		return true
	}

	return true
}
