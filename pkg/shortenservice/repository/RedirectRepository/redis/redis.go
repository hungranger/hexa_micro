package redis

import (
	"fmt"
	"hexa_micro/pkg/shortenservice/config"
	"hexa_micro/pkg/shortenservice/container/logger"
	"hexa_micro/pkg/shortenservice/model"
	"hexa_micro/pkg/shortenservice/repository"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type redisRepository struct {
	db *redis.Client
}

func newRedisClient(redisURL string, password string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	opts.Password = password
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(opts)
	_, err = client.Ping().Result()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func NewRedisRepository(redisURL, password string) (repository.IRedirectRepository, error) {
	repo := &redisRepository{}
	client, err := newRedisClient(redisURL, password)
	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRedisRepository")
	} else {
		logger.Log.Info("Connect Redis Successfully")
	}
	repo.db = client
	return repo, nil
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect: %s", code)
}

func (r *redisRepository) Find(code string) (*model.Redirect, error) {
	redirect := &model.Redirect{}
	key := r.generateKey(code)
	data, err := r.db.HGetAll(key).Result()
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	if len(data) == 0 {
		return nil, errors.Wrap(config.ErrRedirectNotFound, "repository.Redirect.Find")
	}
	createAt, err := strconv.ParseInt(data["create_at"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}
	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreateAt = createAt
	return redirect, nil
}

func (r *redisRepository) Store(redirect *model.Redirect) error {
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":      redirect.Code,
		"url":       redirect.URL,
		"create_at": redirect.CreateAt,
	}
	_, err := r.db.HMSet(key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}
