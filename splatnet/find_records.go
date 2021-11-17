package splatnet

import (
	"PeanutButteredSalmon/enums"
	"fmt"
	"log"
	"net/http"
	"time"
)

func filterStages(stages []enums.Stage, data ShiftList) []*ShiftSplatnet {
	filteredData := []*ShiftSplatnet{}
	splatnetStages := []stageName{}
	for i := range stages {
		splatnetStages = append(splatnetStages, getSplatnetStage(stages[i]))
	}
	for i := range data.Results {
		if data.Results[i].Schedule.Stage.Name.isElementExist(splatnetStages) {
			filteredData = append(filteredData, &data.Results[i])
		}
	}
	return filteredData
}

func filterEvents(events []enums.Event, shifts []*ShiftSplatnet) []*ShiftSplatnet {
	splatnetEvents := []event{}
	for i := range events {
		splatnetEvents = append(splatnetEvents, getSplatnetEvent(events[i]))
	}
	filteredShifts := []*ShiftSplatnet{}
	for i := range shifts {
		if shifts[i].WaveDetails[0].EventType.Key.isElementExist(splatnetEvents) || (len(shifts[i].WaveDetails) > 1 &&
			(shifts[i].WaveDetails[1].EventType.Key.isElementExist(splatnetEvents) || (len(shifts[i].WaveDetails) > 2 &&
				shifts[i].WaveDetails[2].EventType.Key.isElementExist(splatnetEvents)))) {
			filteredShifts = append(filteredShifts, shifts[i])
		}
	}
	return filteredShifts
}

func filterTides(tides []enums.Tide, shifts []*ShiftSplatnet) []*ShiftSplatnet {
	splatnetTides := []tide{}
	for i := range tides {
		splatnetTides = append(splatnetTides, getSplatnetTide(tides[i]))
	}
	filteredShifts := []*ShiftSplatnet{}
	for i := range shifts {
		if shifts[i].WaveDetails[0].WaterLevel.Key.isElementExist(splatnetTides) || (len(shifts[i].WaveDetails) > 1 &&
			(shifts[i].WaveDetails[1].WaterLevel.Key.isElementExist(splatnetTides) || (len(shifts[i].WaveDetails) > 2 &&
				shifts[i].WaveDetails[2].WaterLevel.Key.isElementExist(splatnetTides)))) {
			filteredShifts = append(filteredShifts, shifts[i])
		}
	}
	return filteredShifts
}

func filterWeapons(weapons []enums.WeaponSchedule, shifts []*ShiftSplatnet) []*ShiftSplatnet {
	goldenRandom := false
	singleRandom := false
	quadRandom := false
	normal := false
	for i := range weapons {
		goldenRandom = goldenRandom || weapons[i] == enums.RandommGrizzco
		singleRandom = singleRandom || weapons[i] == enums.SingleRandom
		quadRandom = quadRandom || weapons[i] == enums.FourRandom
		normal = normal || weapons[i] == enums.Set
	}
	filteredShifts := []*ShiftSplatnet{}
	for i := range shifts {
		if (shifts[i].Schedule.Weapons[0].CoopSpecialWeapon != nil &&
			shifts[i].Schedule.Weapons[1].CoopSpecialWeapon != nil &&
			shifts[i].Schedule.Weapons[2].CoopSpecialWeapon != nil &&
			shifts[i].Schedule.Weapons[3].CoopSpecialWeapon != nil &&
			((goldenRandom &&
				shifts[i].Schedule.Weapons[0].CoopSpecialWeapon.Name == RandomGrizzcoSchedule &&
				shifts[i].Schedule.Weapons[1].CoopSpecialWeapon.Name == RandomGrizzcoSchedule &&
				shifts[i].Schedule.Weapons[2].CoopSpecialWeapon.Name == RandomGrizzcoSchedule &&
				shifts[i].Schedule.Weapons[3].CoopSpecialWeapon.Name == RandomGrizzcoSchedule) ||
				(quadRandom &&
					shifts[i].Schedule.Weapons[0].CoopSpecialWeapon.Name == RandomSchedule &&
					shifts[i].Schedule.Weapons[1].CoopSpecialWeapon.Name == RandomSchedule &&
					shifts[i].Schedule.Weapons[2].CoopSpecialWeapon.Name == RandomSchedule &&
					shifts[i].Schedule.Weapons[3].CoopSpecialWeapon.Name == RandomSchedule))) ||
			(shifts[i].Schedule.Weapons[0].Weapon != nil &&
				shifts[i].Schedule.Weapons[1].Weapon != nil &&
				shifts[i].Schedule.Weapons[2].Weapon != nil &&
				((singleRandom &&
					shifts[i].Schedule.Weapons[3].CoopSpecialWeapon != nil &&
					shifts[i].Schedule.Weapons[3].CoopSpecialWeapon.Name == RandomSchedule) ||
					(normal &&
						shifts[i].Schedule.Weapons[3].Weapon != nil))) {
			filteredShifts = append(filteredShifts, shifts[i])
		}
	}
	return filteredShifts
}

func getWeaponsScheduleShift(shift *ShiftSplatnet) enums.WeaponSchedule {
	if shift.Schedule.Weapons[0].CoopSpecialWeapon != nil &&
		shift.Schedule.Weapons[1].CoopSpecialWeapon != nil &&
		shift.Schedule.Weapons[2].CoopSpecialWeapon != nil &&
		shift.Schedule.Weapons[3].CoopSpecialWeapon != nil &&
		shift.Schedule.Weapons[0].CoopSpecialWeapon.Name == RandomGrizzcoSchedule &&
		shift.Schedule.Weapons[1].CoopSpecialWeapon.Name == RandomGrizzcoSchedule &&
		shift.Schedule.Weapons[2].CoopSpecialWeapon.Name == RandomGrizzcoSchedule &&
		shift.Schedule.Weapons[3].CoopSpecialWeapon.Name == RandomGrizzcoSchedule {
		return enums.RandommGrizzco
	}
	if shift.Schedule.Weapons[0].CoopSpecialWeapon != nil &&
		shift.Schedule.Weapons[1].CoopSpecialWeapon != nil &&
		shift.Schedule.Weapons[2].CoopSpecialWeapon != nil &&
		shift.Schedule.Weapons[3].CoopSpecialWeapon != nil &&
		shift.Schedule.Weapons[0].CoopSpecialWeapon.Name == RandomSchedule &&
		shift.Schedule.Weapons[1].CoopSpecialWeapon.Name == RandomSchedule &&
		shift.Schedule.Weapons[2].CoopSpecialWeapon.Name == RandomSchedule &&
		shift.Schedule.Weapons[3].CoopSpecialWeapon.Name == RandomSchedule {
		return enums.FourRandom
	}
	if shift.Schedule.Weapons[0].Weapon != nil &&
		shift.Schedule.Weapons[1].Weapon != nil &&
		shift.Schedule.Weapons[2].Weapon != nil &&
		shift.Schedule.Weapons[3].CoopSpecialWeapon != nil &&
		shift.Schedule.Weapons[3].CoopSpecialWeapon.Name == RandomSchedule {
		return enums.SingleRandom
	}
	if shift.Schedule.Weapons[0].Weapon != nil &&
		shift.Schedule.Weapons[1].Weapon != nil &&
		shift.Schedule.Weapons[2].Weapon != nil &&
		shift.Schedule.Weapons[3].Weapon != nil {
		return enums.Set
	}
	return ""
}

func FindRecords(stages []enums.Stage, hasEvents []enums.Event, tides []enums.Tide, weapons []enums.WeaponSchedule, save bool, appHead http.Header, client *http.Client) {
	data := getAllShifts(appHead, client, save)
	stageFiltered := filterStages(stages, data)
	eventsFiltered := filterEvents(hasEvents, stageFiltered)
	tidesFiltered := filterTides(tides, eventsFiltered)
	finishedFiltering := filterWeapons(weapons, tidesFiltered)
	log.Println(len(stageFiltered))
	records := getAllRecords()
	for i := range finishedFiltering {
		totalEggs := finishedFiltering[i].TotalEggs()
		weaponsType := getWeaponsScheduleShift(finishedFiltering[i])
		if records[totalGoldenEggs][finishedFiltering[i].Schedule.Stage.Name][weaponsType].RecordAmount < totalEggs {
			records[totalGoldenEggs][finishedFiltering[i].Schedule.Stage.Name][weaponsType] = &Record{
				Time:         time.Unix(finishedFiltering[i].PlayTime, 0),
				RecordAmount: totalEggs,
				Shift:        *finishedFiltering[i],
			}
		}
		nightCount := 0
		for l := range finishedFiltering[i].WaveDetails {
			for j := range hasEvents {
				for k := range tides {
					if hasEvents[j] == enums.WaterLevels || (tides[k] == enums.Lt && hasEvents[j] == enums.GoldieSeeking) || (tides[k] != enums.Lt && (hasEvents[j] == enums.CohockCharge || hasEvents[j] == enums.Griller)) {
						continue
					}
					if finishedFiltering[i].WaveDetails[l].EventType.Key == getSplatnetEvent(hasEvents[j]) &&
						finishedFiltering[i].WaveDetails[l].WaterLevel.Key == getSplatnetTide(tides[k]) &&
						finishedFiltering[i].WaveDetails[l].GoldenEggs > records[RecordName(string(tides[k])+" "+hasEvents[j].ToString())][finishedFiltering[i].Schedule.Stage.Name][weaponsType].RecordAmount {
						records[RecordName(string(tides[k])+" "+hasEvents[j].ToString())][finishedFiltering[i].Schedule.Stage.Name][weaponsType] = &Record{
							Time:         time.Unix(finishedFiltering[i].PlayTime, 0),
							RecordAmount: finishedFiltering[i].WaveDetails[l].GoldenEggs,
							Shift:        *finishedFiltering[i],
						}
					}
				}
			}
			if finishedFiltering[i].WaveDetails[l].EventType.Key != waterLevels {
				nightCount++
			}
		}
		if nightCount <= 2 && records[totalGoldenEggsTwoNight][finishedFiltering[i].Schedule.Stage.Name][weaponsType].RecordAmount < totalEggs {
			records[totalGoldenEggsTwoNight][finishedFiltering[i].Schedule.Stage.Name][weaponsType] = &Record{
				Time:         time.Unix(finishedFiltering[i].PlayTime, 0),
				RecordAmount: totalEggs,
				Shift:        *finishedFiltering[i],
			}
		}
		if nightCount <= 1 && records[totalGoldenEggsOneNight][finishedFiltering[i].Schedule.Stage.Name][weaponsType].RecordAmount < totalEggs {
			records[totalGoldenEggsOneNight][finishedFiltering[i].Schedule.Stage.Name][weaponsType] = &Record{
				Time:         time.Unix(finishedFiltering[i].PlayTime, 0),
				RecordAmount: totalEggs,
				Shift:        *finishedFiltering[i],
			}
		}
		if nightCount == 0 && records[totalGoldenEggsNoNight][finishedFiltering[i].Schedule.Stage.Name][weaponsType].RecordAmount < totalEggs {
			records[totalGoldenEggsNoNight][finishedFiltering[i].Schedule.Stage.Name][weaponsType] = &Record{
				Time:         time.Unix(finishedFiltering[i].PlayTime, 0),
				RecordAmount: totalEggs,
				Shift:        *finishedFiltering[i],
			}
		}
	}
	fmt.Println("{")
	recordNames := getRecordNames()
	for i := range recordNames {
		fmt.Printf("\t\"%s\": {\n", recordNames[i])
		for j := range stages {
			fmt.Printf("\t\t\"%s\": {\n", getSplatnetStage(stages[j]))
			for k := range weapons {
				fmt.Printf("\t\t\t\"%v\": {\n\t\t\t\t\"Golden Eggs\": %d\n\t\t\t\t\"Time\": \"%s\"\n\t\t\t}", weapons[k], records[recordNames[i]][getSplatnetStage(stages[j])][weapons[k]].RecordAmount, records[recordNames[i]][getSplatnetStage(stages[j])][weapons[k]].Time.Format("2006-01-02 15-04-05"))
				if k < len(weapons)-1 {
					fmt.Print(",\n")
				}
			}
			fmt.Print("\n\t\t}")
			if j < len(stages)-1 {
				fmt.Print(",\n")
			}
		}
		fmt.Print("\n\t}")
		if i < len(recordNames)-1 {
			fmt.Print(",\n")
		}
	}
	fmt.Println("\n}")
}
