package database

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/winstonjr/goexpert-desafio-rate-limiter/internal/entity"
	"time"
)

var ctx = context.Background()

type InteractionDTO struct {
	NumberOfInteractions uint  `json:"numberOfInteractions"`
	AllowedInteractions  uint  `json:"allowedInteractions"`
	AllowedInterval      int64 `json:"allowedInterval"`
	BlockInterval        int64 `json:"blockInterval"`
	Blocked              bool  `json:"blocked"`
}

type FilterStoreRedis struct {
	rdb           *redis.Client
	initialConfig map[string]*entity.TokenBucketConfig
}

func NewFilterStoreRedis(initialConfig map[string]*entity.TokenBucketConfig, client *redis.Client) (*FilterStoreRedis, error) {
	if initialConfig == nil || len(initialConfig) == 0 {
		return nil, errors.New("you need to provide at least one initial config")
	}

	return &FilterStoreRedis{
		rdb:           client,
		initialConfig: initialConfig,
	}, nil
}

func (f *FilterStoreRedis) InsideLimit(key string) bool {
	interaction := f.returnInteraction(key)
	return entity.ValidateRules(key, interaction, f.createEmptyInteraction)
}

func (f *FilterStoreRedis) returnInteraction(key string) *entity.Interaction {
	if interaction, err := f.rdb.Get(ctx, "key").Result(); err == nil {
		return f.stringToInteraction(interaction)
	} else if interaction, err = f.rdb.Get(ctx, "*").Result(); err == nil {
		return f.stringToInteraction(interaction)
	}

	return f.createEmptyInteraction(key)
}

func (f *FilterStoreRedis) createEmptyInteraction(key string) *entity.Interaction {
	actualConfig := &entity.TokenBucketConfig{
		MaxRequests:    10_000_000_000,
		LimitInSeconds: 60 * 60,
		BlockInSeconds: 0,
	}
	if specificConfig, ok := f.initialConfig[key]; ok {
		actualConfig = specificConfig
	}

	now := time.Now().Unix()

	blockingSeconds := actualConfig.LimitInSeconds + actualConfig.BlockInSeconds
	newInteraction := &entity.Interaction{
		NumberOfInteractions: 0,
		AllowedInteractions:  actualConfig.MaxRequests,
		AllowedInterval:      now + actualConfig.LimitInSeconds,
		BlockInterval:        now + blockingSeconds,
		Blocked:              false,
	}
	expiration := time.Duration(blockingSeconds) * time.Second
	f.rdb.Set(ctx, key, f.interactionToString(newInteraction), expiration)

	return newInteraction
}

func (f *FilterStoreRedis) stringToInteraction(jsonString string) *entity.Interaction {
	var result InteractionDTO
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		panic(err)
	}
	return &entity.Interaction{
		NumberOfInteractions: result.NumberOfInteractions,
		AllowedInteractions:  result.AllowedInteractions,
		AllowedInterval:      result.AllowedInterval,
		BlockInterval:        result.BlockInterval,
		Blocked:              result.Blocked,
	}
}

func (f *FilterStoreRedis) interactionToString(interaction *entity.Interaction) string {
	input := &InteractionDTO{
		NumberOfInteractions: interaction.NumberOfInteractions,
		AllowedInteractions:  interaction.AllowedInteractions,
		AllowedInterval:      interaction.AllowedInterval,
		BlockInterval:        interaction.BlockInterval,
		Blocked:              interaction.Blocked,
	}
	result, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}
	return string(result)
}
