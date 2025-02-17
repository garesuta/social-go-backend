package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/backend-production-go-1/internal/store"
	"github.com/go-redis/redis/v8"
)

type UserStore struct {
	rdb *redis.Client
}

const UserExpTime = time.Minute

func (s *UserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	//test recheck rdb
	if s.rdb == nil {
		return nil, fmt.Errorf("redis client is not initialized")
	}
	//check if user is nil
	if userID == 0 {
		return nil, fmt.Errorf("can't find userID")
	}

	cacheKey := fmt.Sprintf("user-%v", userID)
	//{"user":42}
	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	var user store.User
	if data != "" {
		err := json.Unmarshal([]byte(data), &user)
		if err != nil {
			return nil, err
		}
	}
	return &user, nil
}

func (s *UserStore) Set(ctx context.Context, user *store.User) error {
	//TTL
	cacheKey := fmt.Sprintf("user-%v", user.ID)
	json, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return s.rdb.SetEX(ctx, cacheKey, json, UserExpTime).Err()
}
