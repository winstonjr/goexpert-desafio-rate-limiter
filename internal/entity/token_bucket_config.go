package entity

type TokenBucketConfig struct {
	MaxRequests    uint
	LimitInSeconds int64
	BlockInSeconds int64
}
