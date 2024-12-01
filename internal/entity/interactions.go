package entity

import "time"

type Interaction struct {
	NumberOfInteractions uint
	AllowedInteractions  uint
	AllowedInterval      int64
	BlockInterval        int64
	Blocked              bool
}

func ValidateRules(key string, interaction *Interaction, cei func(key string) *Interaction) bool {
	now := time.Now().Unix()

	if interaction.Blocked && interaction.BlockInterval < now {
		interaction = cei(key)
		interaction.NumberOfInteractions = 1
		return true
	} else if interaction.Blocked {
		return false
	} else if !interaction.Blocked && interaction.AllowedInterval < now {
		interaction = cei(key)
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
