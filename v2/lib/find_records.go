package lib

import (
	"fmt"
	"github.com/cass-dlcm/PeanutButteredSalmon/v2/schedules"
	"github.com/cass-dlcm/PeanutButteredSalmon/v2/types"
	"net/http"
)

func filterStages(stages []types.Stage, data Shift, schedules schedules.Schedule) Shift {
	stage := data.GetStage(schedules)
	if stage.IsElementExists(stages) {
		return data
	}
	return nil
}

func filterEvents(events []types.Event, data Shift) Shift {
	getEvents := data.GetEvents()
	if getEvents.IsAllElementExist(events) {
		return data
	}
	return nil
}

func filterTides(tides []types.Tide, data Shift) Shift {
	getTides := data.GetTides()
	if getTides.IsAllElementExist(tides) {
		return data
	}
	return nil
}

func filterWeapons(weapons []types.WeaponSchedule, data Shift, schedules schedules.Schedule) Shift {
	set := data.GetWeaponSet(schedules)
	if set.IsElementExists(weapons) {
		return data
	}
	return nil
}

func FindRecords(iterators []ShiftIterator, stages []types.Stage, hasEvents []types.Event, tides []types.Tide, weapons []types.WeaponSchedule, client *http.Client) {
	scheduleList := schedules.GetSchedules(client)
	records := getAllRecords()
	for i := range iterators {
		shift, err := iterators[i].Next()
		if err != nil {
			continue
		}
		for err == nil {
			shift = filterStages(stages, shift, scheduleList)
			if shift == nil {
				shift, err = iterators[i].Next()
				continue
			}
			shift = filterEvents(hasEvents, shift)
			if shift == nil {
				shift, err = iterators[i].Next()
				continue
			}
			shift = filterTides(tides, shift)
			if shift == nil {
				shift, err = iterators[i].Next()
				continue
			}
			shift = filterWeapons(weapons, shift, scheduleList)
			if shift == nil {
				shift, err = iterators[i].Next()
				continue
			}
			totalEggs := shift.GetTotalEggs()
			weaponsType := shift.GetWeaponSet(scheduleList)
			stage := shift.GetStage(scheduleList)
			if records[totalGoldenEggs] == nil {
				records[totalGoldenEggs] = &map[types.Stage]*map[types.WeaponSchedule]*Record{}
			}
			if (*records[totalGoldenEggs])[stage] == nil {
				(*records[totalGoldenEggs])[stage] = &map[types.WeaponSchedule]*Record{}
			}
			if (*(*records[totalGoldenEggs])[stage])[weaponsType] == nil || (*(*records[totalGoldenEggs])[stage])[weaponsType].RecordAmount < totalEggs {
				(*(*records[totalGoldenEggs])[stage])[weaponsType] = &Record{
					Time:         shift.GetTime(),
					RecordAmount: totalEggs,
					ShiftData:    shift,
					Identifier:   shift.GetIdentifier(),
				}
			}
			nightCount := 0
			waveCount := shift.GetWaveCount()
			waveEggs := shift.GetEggsWaves()
			waveEvents := shift.GetEvents()
			waveWaterLevel := shift.GetTides()
			for l := 0; l < waveCount; l++ {
				for j := range hasEvents {
					for k := range tides {
						if hasEvents[j] == types.WaterLevels || (tides[k] == types.Lt && hasEvents[j] == types.GoldieSeeking) || (tides[k] != types.Lt && (hasEvents[j] == types.CohockCharge || hasEvents[j] == types.Griller)) {
							continue
						}
						if waveEvents[l] == hasEvents[j] &&
							waveWaterLevel[l] == tides[k] {
							if records[RecordName(string(tides[k])+" "+hasEvents[j].ToString())] == nil {
								records[RecordName(string(tides[k])+" "+hasEvents[j].ToString())] = &map[types.Stage]*map[types.WeaponSchedule]*Record{}
							}
							if (*records[RecordName(string(tides[k])+" "+hasEvents[j].ToString())])[stage] == nil {
								(*records[RecordName(string(tides[k])+" "+hasEvents[j].ToString())])[stage] = &map[types.WeaponSchedule]*Record{}
							}
							if (*(*records[RecordName(string(tides[k])+" "+hasEvents[j].ToString())])[stage])[weaponsType] == nil || waveEggs[l] > (*(*records[RecordName(string(tides[k])+" "+hasEvents[j].ToString())])[stage])[weaponsType].RecordAmount {
								(*(*records[RecordName(string(tides[k])+" "+hasEvents[j].ToString())])[stage])[weaponsType] = &Record{
									Time:         shift.GetTime(),
									RecordAmount: waveEggs[l],
									ShiftData:    shift,
									Identifier:   shift.GetIdentifier(),
								}
							}
						}
					}
				}
				if waveEvents[l] != types.WaterLevels {
					nightCount++
				}
			}
			if nightCount <= 2 {
				if records[totalGoldenEggsTwoNight] == nil {
					records[totalGoldenEggsTwoNight] = &map[types.Stage]*map[types.WeaponSchedule]*Record{}
				}
				if (*records[totalGoldenEggsTwoNight])[stage] == nil {
					(*records[totalGoldenEggsTwoNight])[stage] = &map[types.WeaponSchedule]*Record{}
				}
				if (*(*records[totalGoldenEggsTwoNight])[stage])[weaponsType] == nil || (*(*records[totalGoldenEggsTwoNight])[stage])[weaponsType].RecordAmount < totalEggs {
					(*(*records[totalGoldenEggsTwoNight])[stage])[weaponsType] = &Record{
						Time:         shift.GetTime(),
						RecordAmount: totalEggs,
						ShiftData:    shift,
						Identifier:   shift.GetIdentifier(),
					}
				}
			}
			if nightCount <= 1 {
				if records[totalGoldenEggsOneNight] == nil {
					records[totalGoldenEggsOneNight] = &map[types.Stage]*map[types.WeaponSchedule]*Record{}
				}
				if (*records[totalGoldenEggsOneNight])[stage] == nil {
					(*records[totalGoldenEggsOneNight])[stage] = &map[types.WeaponSchedule]*Record{}
				}
				if (*(*records[totalGoldenEggsOneNight])[stage])[weaponsType] == nil || (*(*records[totalGoldenEggsOneNight])[stage])[weaponsType].RecordAmount < totalEggs {
					(*(*records[totalGoldenEggsOneNight])[stage])[weaponsType] = &Record{
						Time:         shift.GetTime(),
						RecordAmount: totalEggs,
						ShiftData:    shift,
						Identifier:   shift.GetIdentifier(),
					}
				}
			}
			if nightCount == 0 {
				if records[totalGoldenEggsNoNight] == nil {
					records[totalGoldenEggsNoNight] = &map[types.Stage]*map[types.WeaponSchedule]*Record{}
				}
				if (*records[totalGoldenEggsNoNight])[stage] == nil {
					(*records[totalGoldenEggsNoNight])[stage] = &map[types.WeaponSchedule]*Record{}
				}
				if (*(*records[totalGoldenEggsNoNight])[stage])[weaponsType] == nil || (*(*records[totalGoldenEggsNoNight])[stage])[weaponsType].RecordAmount < totalEggs {
					(*(*records[totalGoldenEggsNoNight])[stage])[weaponsType] = &Record{
						Time:         shift.GetTime(),
						RecordAmount: totalEggs,
						ShiftData:    shift,
						Identifier:   shift.GetIdentifier(),
					}
				}
			}
			shift, err = iterators[i].Next()
		}
	}
	recordNames := getRecordNames()
	countsStages := map[RecordName]map[types.Stage]int{}
	countsRecords := map[RecordName]int{}
	for i := range recordNames {
		if records[recordNames[i]] == nil {
			continue
		}
		countsStages[recordNames[i]] = map[types.Stage]int{}
		for j := range stages {
			if (*records[recordNames[i]])[stages[j]] == nil {
				continue
			}
			for k := range weapons {
				if (*(*records[recordNames[i]])[stages[j]])[weapons[k]] == nil {
					continue
				}
				countsStages[recordNames[i]][stages[j]] += 1
			}
			countsRecords[recordNames[i]] += 1
		}
	}
	fmt.Println("{")
	for i := range recordNames {
		if records[recordNames[i]] == nil {
			continue
		}
		fmt.Printf("\t\"%s\": {\n", recordNames[i])
		stagesCount := 0
		for j := range stages {
			if (*records[recordNames[i]])[stages[j]] == nil {
				continue
			}
			fmt.Printf("\t\t\"%s\": {\n", stages[j].ToString())
			weaponsCount := 0
			for k := range weapons {
				if (*(*records[recordNames[i]])[stages[j]])[weapons[k]] == nil {
					continue
				}
				fmt.Printf("\t\t\t\"%v\": {\n\t\t\t\t\"Golden Eggs\": %d,\n\t\t\t\t\"Time\": \"%s\"\n\t\t\t\t\"URL\": \"%s\"\n\t\t\t}", weapons[k], (*(*records[recordNames[i]])[stages[j]])[weapons[k]].RecordAmount, (*(*records[recordNames[i]])[stages[j]])[weapons[k]].Time.Format("2006-01-02 15-04-05"), (*(*records[recordNames[i]])[stages[j]])[weapons[k]].Identifier)
				if weaponsCount < countsStages[recordNames[i]][stages[j]]-1 {
					fmt.Print(",\n")
				}
				weaponsCount++
			}
			fmt.Print("\n\t\t}")
			if stagesCount < countsRecords[recordNames[i]]-1 {
				fmt.Print(",\n")
			}
			stagesCount++
		}
		fmt.Print("\n\t}")
		if i < len(recordNames)-1 {
			fmt.Print(",\n")
		}
	}
	fmt.Println("\n}")
}
