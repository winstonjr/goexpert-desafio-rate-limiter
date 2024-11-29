package database

import (
	"github.com/stretchr/testify/assert"
	"github.com/winstonjr/goexpert-desafio-rate-limiter/internal/entity"
	"testing"
	"time"
)

func TestFilterStoreInMemory_LimitPass(t *testing.T) {
	limiterConfig := make(map[string]*entity.TokenBucketConfig)
	limiterConfig["127.0.0.1"] = &entity.TokenBucketConfig{
		MaxRequests:    1,
		LimitInSeconds: 5,
		BlockInSeconds: 10,
	}
	fsi, err := NewFilterStoreInMemory(limiterConfig)
	assert.Nil(t, err)

	approved := fsi.LimitExceeded("127.0.0.1")
	assert.True(t, approved)
}

func TestFilterStoreInMemory_LimitFailed(t *testing.T) {
	limiterConfig := make(map[string]*entity.TokenBucketConfig)
	limiterConfig["127.0.0.1"] = &entity.TokenBucketConfig{
		MaxRequests:    2,
		LimitInSeconds: 1,
		BlockInSeconds: 2,
	}
	fsi, err := NewFilterStoreInMemory(limiterConfig)
	assert.Nil(t, err)

	assert.True(t, fsi.LimitExceeded("127.0.0.1"))
	assert.True(t, fsi.LimitExceeded("127.0.0.1"))
	assert.False(t, fsi.LimitExceeded("127.0.0.1"))
	time.Sleep(3 * time.Second)
	assert.False(t, fsi.LimitExceeded("127.0.0.1"))
}

func TestFilterStoreInMemory_LimitFailsAndThenPass(t *testing.T) {
	limiterConfig := make(map[string]*entity.TokenBucketConfig)
	limiterConfig["127.0.0.1"] = &entity.TokenBucketConfig{
		MaxRequests:    2,
		LimitInSeconds: 1,
		BlockInSeconds: 2,
	}
	fsi, err := NewFilterStoreInMemory(limiterConfig)
	assert.Nil(t, err)

	assert.True(t, fsi.LimitExceeded("127.0.0.1"))
	assert.True(t, fsi.LimitExceeded("127.0.0.1"))
	assert.False(t, fsi.LimitExceeded("127.0.0.1"))
	time.Sleep(4 * time.Second)
	assert.True(t, fsi.LimitExceeded("127.0.0.1"))
}

func TestFilterStoreInMemory_LimitBiggerRequests(t *testing.T) {
	limiterConfig := make(map[string]*entity.TokenBucketConfig)
	limiterConfig["abc123"] = &entity.TokenBucketConfig{
		MaxRequests:    10,
		LimitInSeconds: 1,
		BlockInSeconds: 2,
	}
	fsi, err := NewFilterStoreInMemory(limiterConfig)
	assert.Nil(t, err)

	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.False(t, fsi.LimitExceeded("abc123"))
	time.Sleep(3 * time.Second)
	assert.False(t, fsi.LimitExceeded("abc123"))
	time.Sleep(1 * time.Second)
	assert.True(t, fsi.LimitExceeded("abc123"))

	approved := fsi.LimitExceeded("abc123")
	assert.True(t, approved)
}

func TestFilterStoreInMemory_LimitBiggerRequestsNoBlockAfterLimit(t *testing.T) {
	limiterConfig := make(map[string]*entity.TokenBucketConfig)
	limiterConfig["abc123"] = &entity.TokenBucketConfig{
		MaxRequests:    10,
		LimitInSeconds: 2,
		BlockInSeconds: 0,
	}
	fsi, err := NewFilterStoreInMemory(limiterConfig)
	assert.Nil(t, err)

	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.True(t, fsi.LimitExceeded("abc123"))
	assert.False(t, fsi.LimitExceeded("abc123"))
	time.Sleep(2 * time.Second)
	assert.False(t, fsi.LimitExceeded("abc123"))
	time.Sleep(1 * time.Second)
	assert.True(t, fsi.LimitExceeded("abc123"))

	approved := fsi.LimitExceeded("abc123")
	assert.True(t, approved)
}
