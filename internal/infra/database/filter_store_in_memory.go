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

func (f *FilterStoreInMemory) LimitExceeded(key string) bool {
	interaction := f.returnInteraction(key)
	return f.validateRules(key, interaction)
}

func (f *FilterStoreInMemory) validateRules(key string, interaction *entity.Interaction) bool {
	now := time.Now().Unix()

	if interaction.Blocked && interaction.BlockInterval < now {
		interaction = f.createEmptyInteraction(key)
		interaction.NumberOfInteractions = 1
		return true
	} else if interaction.Blocked {
		return false
	} else if !interaction.Blocked && interaction.AllowedInterval < now {
		interaction = f.createEmptyInteraction(key)
		interaction.NumberOfInteractions = 1
		return true
	} else {
		nextInteraction := interaction.NumberOfInteractions + 1
		if nextInteraction <= interaction.AllowedInteractions {
			interaction.NumberOfInteractions = nextInteraction
			return true
		} else {
			interaction.Blocked = true
			return false
		}
	}
}

func (f *FilterStoreInMemory) returnInteraction(key string) *entity.Interaction {
	if interaction, ok := f.interactions[key]; ok {
		return interaction
	}

	return f.createEmptyInteraction(key)
}

func (f *FilterStoreInMemory) createEmptyInteraction(key string) *entity.Interaction {
	actualConfig := &entity.TokenBucketConfig{
		MaxRequests:    1_000_000,
		LimitInSeconds: 365 * 24 * 60 * 60,
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
