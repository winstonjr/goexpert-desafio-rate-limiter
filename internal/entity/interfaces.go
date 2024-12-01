package entity

type FilterStoreInterface interface {
	InsideLimit(key string) bool
}
