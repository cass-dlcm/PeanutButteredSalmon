package types

import (
	"fmt"
	"strings"
)

type Stage int

const (
	SpawningGrounds Stage = iota
	MaroonersBay
	LostOutpost
	SalmonidSmokeyard
	RuinsOfArkPolaris
)

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

func stringToStage(stageStr string) (*Stage, error) {
	var stageRes Stage
	switch stageStr {
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
		return nil, fmt.Errorf("stage not found: %s", stageStr)
	}
	return &stageRes, nil
}

func GetStageArgs(stagesStr string) ([]Stage, error) {
	stages := []Stage{}
	stagesStrArr := strings.Split(stagesStr, " ")
	for i := range stagesStrArr {
		stageVal, err := stringToStage(stagesStrArr[i])
		if err != nil {
			return nil, err
		}
		stages = append(stages, *stageVal)
	}
	return stages, nil
}

func GetAllStages() []Stage {
	return []Stage{
		SpawningGrounds,
		MaroonersBay,
		LostOutpost,
		SalmonidSmokeyard,
		RuinsOfArkPolaris,
	}
}

func (s *Stage) IsElementExists(arr []Stage) bool {
	for _, v := range arr {
		if v == *s {
			return true
		}
	}
	return false
}
