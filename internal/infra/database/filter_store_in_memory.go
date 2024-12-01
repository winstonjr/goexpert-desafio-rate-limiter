package database

import (
	"errors"
	"github.com/winstonjr/goexpert-desafio-rate-limiter/internal/entity"
	"time"
)

type FilterStoreInMemory struct {
	interactions  map[string]*entity.Interaction
	initialConfig map[string]*entity.TokenBucketConfig
}

func NewFilterStoreInMemory(initialConfig map[string]*entity.TokenBucketConfig) (*FilterStoreInMemory, error) {
	if initialConfig == nil || len(initialConfig) == 0 {
		return nil, errors.New("you need to provide at least one initial config")
	}

	return &FilterStoreInMemory{
		interactions:  make(map[string]*entity.Interaction),
		initialConfig: initialConfig,
	}, nil
}

func (f *FilterStoreInMemory) InsideLimit(key string) bool {
	interaction := f.returnInteraction(key)
	return entity.ValidateRules(key, interaction, f.createEmptyInteraction, f.updateInteraction)
}

func (f *FilterStoreInMemory) returnInteraction(key string) *entity.Interaction {
	if interaction, ok := f.interactions[key]; ok {
		return interaction
	} else if interaction, ok = f.interactions["*"]; ok {
		return interaction
	}

	return f.createEmptyInteraction(key)
}

func (f *FilterStoreInMemory) createEmptyInteraction(key string) *entity.Interaction {
	actualConfig := &entity.TokenBucketConfig{
		MaxRequests:    10_000_000_000,
		LimitInSeconds: 60 * 60,
		BlockInSeconds: 0,
	}
	if specificConfig, ok := f.initialConfig[key]; ok {
		actualConfig = specificConfig
	}

	now := time.Now().Unix()

	newInteraction := &entity.Interaction{
		NumberOfInteractions: 0,
		AllowedInteractions:  actualConfig.MaxRequests,
		AllowedInterval:      now + actualConfig.LimitInSeconds,
		BlockInterval:        now + actualConfig.LimitInSeconds + actualConfig.BlockInSeconds,
		Blocked:              false,
	}
	f.interactions[key] = newInteraction

	return newInteraction
}

func (f *FilterStoreInMemory) updateInteraction(key string, interaction *entity.Interaction) {
	interaction.NumberOfInteractions = interaction.NumberOfInteractions + 1
}
