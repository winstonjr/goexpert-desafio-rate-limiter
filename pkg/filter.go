package pkg

import "github.com/winstonjr/goexpert-desafio-rate-limiter/internal/entity"

type Filter struct {
	store       entity.FilterStoreInterface
	configToken map[string]entity.TokenBucketConfig
}

func NewFilter(store entity.FilterStoreInterface) *Filter {
	return &Filter{
		store:       store,
		configToken: make(map[string]entity.TokenBucketConfig),
	}
}
