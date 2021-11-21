package types

import (
	"fmt"
	"strings"
)

// Stage is an int enum of the stage for the rotation.
type Stage int

// The five Salmon Run stages.
const (
	SpawningGrounds Stage = iota
	MaroonersBay
	LostOutpost
	SalmonidSmokeyard
	RuinsOfArkPolaris
)

// ToString returns the name of the Stage, currently hardcoded as the en-US locale.
func (s Stage) ToString() string {
	switch s {
	case SpawningGrounds:
		return "Spawning Grounds"
	case MaroonersBay:
		return "Marooner's Bay"
	case LostOutpost:
		return "Lost Outpost"
	case SalmonidSmokeyard:
		return "Salmonid Smokeyard"
	case RuinsOfArkPolaris:
		return "Ruins of Ark Polaris"
	}
	return ""
}

// GetStageArgs turns a string of space seperated string representations of stages into a slice of Stage.
func GetStageArgs(stagesStr string) ([]Stage, error) {
	stages := []Stage{}
	stagesStrArr := strings.Split(stagesStr, " ")
	for i := range stagesStrArr {
		var stageRes Stage
		switch stagesStrArr[i] {
		case "spawning_grounds":
			stageRes = SpawningGrounds
		case "marooners_bay":
			stageRes = MaroonersBay
		case "lost_outpost":
			stageRes = LostOutpost
		case "salmonid_smokeyard":
			stageRes = SalmonidSmokeyard
		case "ruins_of_ark_polaris":
			stageRes = RuinsOfArkPolaris
		default:
			return nil, fmt.Errorf("stage not found: %s", stagesStrArr[i])
		}
		stages = append(stages, stageRes)
	}
	return stages, nil
}

// GetAllStages returns a slice containing every Stage constant.
func GetAllStages() []Stage {
	return []Stage{
		SpawningGrounds,
		MaroonersBay,
		LostOutpost,
		SalmonidSmokeyard,
		RuinsOfArkPolaris,
	}
}

// IsElementExists finds whether the given Stage is in the Stage slice.
func (s *Stage) IsElementExists(arr []Stage) bool {
	for _, v := range arr {
		if v == *s {
			return true
		}
	}
	return false
}
