package configs

import (
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"github.com/winstonjr/goexpert-desafio-rate-limiter/internal/entity"
)

type TokenBucketConfigDTO struct {
	MaxRequests    uint  `json:"maxRequests"`
	LimitInSeconds int64 `json:"limitInSeconds"`
	BlockInSeconds int64 `json:"blockInSeconds"`
}

type Conf struct {
	RateLimiterRulesJSON string `mapstructure:"RATE_LIMITER_RULES"`
	RateLimiterRules     map[string]*entity.TokenBucketConfig
	RedisAddress         string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword        string `mapstructure:"REDIS_PASSWORD"`
	RedisDb              int    `mapstructure:"REDIS_DB"`
	StoreKind            string `mapstructure:"STORE_KIND"`
}

func LoadConfig(path string) (*Conf, error) {
	var cfg *Conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err == nil {
		err = viper.Unmarshal(&cfg)
		if err != nil {
			panic(err)
		}
	} else {
		rateLimiterRulesJSON := viper.GetString("RATE_LIMITER_RULES")
		storeKind := viper.GetString("STORE_KIND")
		if rateLimiterRulesJSON == "" {
			return nil, errors.New("no rate limiter rules configured")
		}
		cfg = &Conf{
			RateLimiterRulesJSON: rateLimiterRulesJSON,
			StoreKind:            storeKind,
		}
		if storeKind == "redis" {
			redisAddress := viper.GetString("REDIS_ADDRESS")
			redisPassword := viper.GetString("REDIS_PASSWORD")
			redisDb := viper.GetInt("REDIS_DB")

			if redisAddress == "" {
				return nil, errors.New("no redis address configured")
			}

			cfg.RedisAddress = redisAddress
			cfg.RedisPassword = redisPassword
			cfg.RedisDb = redisDb
		}
	}

	rulesInObject, err := cfg.getRateLimiterRules(cfg.RateLimiterRulesJSON)
	if err != nil {
		return nil, err
	}
	cfg.RateLimiterRules = rulesInObject

	return cfg, err
}

func (c *Conf) getRateLimiterRules(jsonString string) (map[string]*entity.TokenBucketConfig, error) {
	var result map[string]TokenBucketConfigDTO
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, err
	}
	retVal := make(map[string]*entity.TokenBucketConfig)
	for k, v := range result {
		retVal[k] = v.toEntity()
	}
	return retVal, nil
}

func (t *TokenBucketConfigDTO) toEntity() *entity.TokenBucketConfig {
	return &entity.TokenBucketConfig{
		MaxRequests:    t.MaxRequests,
		LimitInSeconds: t.LimitInSeconds,
		BlockInSeconds: t.BlockInSeconds,
	}
}
