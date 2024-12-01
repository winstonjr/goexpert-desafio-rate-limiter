package entity

import "time"

type Interaction struct {
	NumberOfInteractions uint
	AllowedInteractions  uint
	AllowedInterval      int64
	BlockInterval        int64
	Blocked              bool
	Expiration           time.Duration
}

func ValidateRules(key string, interaction *Interaction, createEmptyInteraction func(key string) *Interaction, updateInteraction func(key string, interaction *Interaction)) bool {
	now := time.Now().Unix()

	if interaction.Blocked && interaction.BlockInterval < now {
		interaction = createEmptyInteraction(key)
		updateInteraction(key, interaction)
		return true
	} else if interaction.Blocked {
		return false
	} else if !interaction.Blocked && interaction.AllowedInterval < now {
		interaction = createEmptyInteraction(key)
		updateInteraction(key, interaction)
		return true
	} else {
		nextInteraction := interaction.NumberOfInteractions + 1
		if nextInteraction <= interaction.AllowedInteractions {
			updateInteraction(key, interaction)
			return true
		} else {
			interaction.Blocked = true
			return false
		}
	}
}
