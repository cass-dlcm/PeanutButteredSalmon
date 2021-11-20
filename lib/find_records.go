package lib

import (
	"fmt"
	"github.com/cass-dlcm/peanutbutteredsalmon/salmonstats"
	"github.com/cass-dlcm/peanutbutteredsalmon/schedules"
	"github.com/cass-dlcm/peanutbutteredsalmon/splatnet"
	"github.com/cass-dlcm/peanutbutteredsalmon/statink"
	"github.com/cass-dlcm/peanutbutteredsalmon/types"
	"net/http"
)

func filterStages(stages []types.Stage, data []Shift, schedules schedules.Schedule) []Shift {
	filteredData := []Shift{}
	for i := range data {
		stage := data[i].GetStage(schedules)
		if stage.IsElementExists(stages) {
			filteredData = append(filteredData, data[i])
		}
	}
	return filteredData
}

func filterEvents(events []types.Event, data []Shift) []Shift {
	filteredData := []Shift{}
	for i := range data {
		getEvents := data[i].GetEvents()
		if getEvents.IsAllElementExist(events) {
			filteredData = append(filteredData, data[i])
		}
	}
	return filteredData
}

func filterTides(tides []types.Tide, data []Shift) []Shift {
	filteredData := []Shift{}
	for i := range data {
		getTides := data[i].GetTides()
		if getTides.IsAllElementExist(tides) {
			filteredData = append(filteredData, data[i])
		}
	}
	return filteredData
}

func filterWeapons(weapons []types.WeaponSchedule, data []Shift, schedules schedules.Schedule) []Shift {
	filteredData := []Shift{}
	for i := range data {
		set := data[i].GetWeaponSet(schedules)
		if set.IsElementExists(weapons) {
			filteredData = append(filteredData, data[i])
		}
	}
	return filteredData
}

func FindRecords(useSplatnet, load bool, statInkServers, salmonStats []types.Server, stages []types.Stage, hasEvents []types.Event, tides []types.Tide, weapons []types.WeaponSchedule, save bool, appHead http.Header, client *http.Client) {
	data := []Shift{}
	if useSplatnet {
		if load {
			tempData := splatnet.LoadFromFile()
			for i := range tempData {
				data = append(data, &tempData[i])
			}
		}
		tempData := splatnet.GetAllShifts(appHead, client, save)
		for i := range tempData.Results {
			data = append(data, &tempData.Results[i])
		}
	}
	for i := range statInkServers {
		id := 1
		if load {
			tempData := statink.LoadFromFile(statInkServers[i])
			for j := range tempData {
				data = append(data, &tempData[j])
			}
			if len(tempData) > 0 {
				id = tempData[len(tempData)-1].ID
			}
		}
		tempData := statink.GetAllShifts(id, statInkServers[i], client, save)
		for j := range tempData {
			data = append(data, &tempData[j])
		}
	}
	for i := range salmonStats {
		page := 1
		if load {
			tempData := salmonstats.LoadFromFile(salmonStats[i])
			for j := range tempData {
				data = append(data, &tempData[j])
			}
			if len(tempData) > 0 {
				page = tempData[len(tempData)-1].Page/200 + 1
			}
		}
		tempData := salmonstats.GetAllShifts(page, salmonStats[i], client, save)
		for j := range tempData {
			data = append(data, &tempData[j])
		}
	}
	scheduleList := schedules.GetSchedules(client)
	stageFiltered := filterStages(stages, data, scheduleList)
	eventsFiltered := filterEvents(hasEvents, stageFiltered)
	tidesFiltered := filterTides(tides, eventsFiltered)
	finishedFiltering := filterWeapons(weapons, tidesFiltered, scheduleList)
	records := getAllRecords()
	for i := range finishedFiltering {
		totalEggs := finishedFiltering[i].GetTotalEggs()
		weaponsType := finishedFiltering[i].GetWeaponSet(scheduleList)
		stage := finishedFiltering[i].GetStage(scheduleList)
		if records[totalGoldenEggs] == nil {
			records[totalGoldenEggs] = &map[types.Stage]*map[types.WeaponSchedule]*Record{}
		}
		if (*records[totalGoldenEggs])[stage] == nil {
			(*records[totalGoldenEggs])[stage] = &map[types.WeaponSchedule]*Record{}
		}
		if (*(*records[totalGoldenEggs])[stage])[weaponsType] == nil || (*(*records[totalGoldenEggs])[stage])[weaponsType].RecordAmount < totalEggs {
			(*(*records[totalGoldenEggs])[stage])[weaponsType] = &Record{
				Time:         finishedFiltering[i].GetTime(),
				RecordAmount: totalEggs,
				ShiftData:    finishedFiltering[i],
				Identifier:   finishedFiltering[i].GetIdentifier(),
			}
		}
		nightCount := 0
		waveCount := finishedFiltering[i].GetWaveCount()
		waveEggs := finishedFiltering[i].GetEggsWaves()
		waveEvents := finishedFiltering[i].GetEvents()
		waveWaterLevel := finishedFiltering[i].GetTides()
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
								Time:         finishedFiltering[i].GetTime(),
								RecordAmount: waveEggs[l],
								ShiftData:    finishedFiltering[i],
								Identifier:   finishedFiltering[i].GetIdentifier(),
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
					Time:         finishedFiltering[i].GetTime(),
					RecordAmount: totalEggs,
					ShiftData:    finishedFiltering[i],
					Identifier:   finishedFiltering[i].GetIdentifier(),
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
					Time:         finishedFiltering[i].GetTime(),
					RecordAmount: totalEggs,
					ShiftData:    finishedFiltering[i],
					Identifier:   finishedFiltering[i].GetIdentifier(),
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
					Time:         finishedFiltering[i].GetTime(),
					RecordAmount: totalEggs,
					ShiftData:    finishedFiltering[i],
					Identifier:   finishedFiltering[i].GetIdentifier(),
				}
			}
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
