package entity

type FilterStoreInterface interface {
	LimitExceeded(key string) bool
}
