package redis

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/rlarkin212/url-shortener/shortener"
)

type redisRepository struct {
	client *redis.Client
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)
	_, err = client.Ping().Result()

	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRedisRepository")
	}

	return client, nil
}

// NewRedisRepository generates new redirect repository
func NewRedisRepository(redisURL string) (shortener.RedirectRespository, error) {
	repo := &redisRepository{}
	client, err := newRedisClient(redisURL)

	if err != nil {
		return nil, errors.Wrap(err, "repository.NewRedisRepository")
	}

	repo.client = client
	return repo, nil
}

func (r *redisRepository) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

func (r *redisRepository) Find(code string) (*shortener.Redirect, error) {
	key := r.generateKey(code)

	data, err := r.client.HGetAll(key).Result()

	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	if len(data) == 0 {
		return nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
	}

	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "repository.Redirect.Find")
	}

	redirect := &shortener.Redirect{
		Code:      data["code"],
		URL:       data["url"],
		CreatedAt: createdAt,
	}

	return redirect, nil
}

func (r *redisRepository) Store(redirect *shortener.Redirect) error {
	key := r.generateKey(redirect.Code)

	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}

	_, err := r.client.HMSet(key, data).Result()

	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}

	return nil
}
