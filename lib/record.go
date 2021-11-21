package lib

import (
	"github.com/cass-dlcm/peanutbutteredsalmon/schedules"
	"github.com/cass-dlcm/peanutbutteredsalmon/types"
	"time"
)

// RecordName is an enum string of the name of each record.
type RecordName string

const (
	totalGoldenEggs                RecordName = "Total Golden Eggs"
	totalGoldenEggsTwoNight        RecordName = "Total Golden Eggs (~2 Night)"
	totalGoldenEggsOneNight        RecordName = "Total Golden Eggs (~1 Night)"
	totalGoldenEggsNoNight         RecordName = "Total Golden Eggs (No Night)"
	singlePlayerGoldenEggs         RecordName = "Single Player Golden Eggs"
	singlePlayerGoldenEggsTwoNight RecordName = "Single Player Golden Eggs (~2 Night)"
	singlePlayerGoldenEggsOneNight RecordName = "Single Player Golden Eggs (~1 Night)"
	singlePlayerGoldenEggsNoNight  RecordName = "Single Player Golden Eggs (No Night)"
	ntNormal                       RecordName = "NT Normal"
	htNormal                       RecordName = "HT Normal"
	ltNormal                       RecordName = "LT Normal"
	ntRush                         RecordName = "NT Rush"
	htRush                         RecordName = "HT Rush"
	ltRush                         RecordName = "LT Rush"
	ntFog                          RecordName = "NT Fog"
	htFog                          RecordName = "HT Fog"
	ltFog                          RecordName = "LT Fog"
	ntGoldieSeeking                RecordName = "NT Goldie Seeking"
	htGoldieSeeking                RecordName = "HT Goldie Seeking"
	ntGriller                      RecordName = "NT Griller"
	htGriller                      RecordName = "HT Griller"
	ntMothership                   RecordName = "NT Mothership"
	htMothership                   RecordName = "HT Mothership"
	ltMothershp                    RecordName = "LT Mothership"
	ltCohockCharge                 RecordName = "LT Cohock Charge"
)

// Shift is a generic source of match results, with only the necessary details available, accessible as methods.
type Shift interface {
	GetTotalEggs() int
	GetStage(schedules.Schedule) types.Stage
	GetWeaponSet(schedules.Schedule) types.WeaponSchedule
	GetEvents() types.EventArr
	GetTides() types.TideArr
	GetEggsWaves() []int
	GetWaveCount() int
	GetTime() time.Time
	GetIdentifier() string
}

// Record is a set of data containing information needed to identify a record holding game.
type Record struct {
	Time         time.Time
	RecordAmount int
	ShiftData    Shift
	Identifier   string
}

func getRecordNames() []RecordName {
	return []RecordName{
		totalGoldenEggs,
		totalGoldenEggsTwoNight,
		totalGoldenEggsOneNight,
		totalGoldenEggsNoNight,
		singlePlayerGoldenEggs,
		singlePlayerGoldenEggsTwoNight,
		singlePlayerGoldenEggsOneNight,
		singlePlayerGoldenEggsNoNight,
		ntNormal,
		htNormal,
		ltNormal,
		ntRush,
		htRush,
		ltRush,
		ntFog,
		htFog,
		ltFog,
		ntGoldieSeeking,
		htGoldieSeeking,
		ntGriller,
		htGriller,
		ntMothership,
		htMothership,
		ltMothershp,
		ltCohockCharge,
	}
}

func getAllRecords() map[RecordName]*map[types.Stage]*map[types.WeaponSchedule]*Record {
	records := map[RecordName]*map[types.Stage]*map[types.WeaponSchedule]*Record{}
	recordNames := getRecordNames()
	for i := range recordNames {
		records[recordNames[i]] = nil
	}
	return records
}
