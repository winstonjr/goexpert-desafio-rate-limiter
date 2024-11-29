package entity

type Interaction struct {
	NumberOfInteractions uint
	AllowedInteractions  uint
	AllowedInterval      int64
	BlockInterval        int64
	Blocked              bool
}
