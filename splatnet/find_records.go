package splatnet

import (
	"github.com/cass-dlcm/peanutbutteredsalmon/schedules"
	"github.com/cass-dlcm/peanutbutteredsalmon/types"
	"time"
)

// GetWeaponSet implements lib.Shift{}.GetWeaponSet(schedules.Schedule) types.WeaponSchedule.
// This function returns what kind of weapon set was used in the shift.
func (s *ShiftSplatnet) GetWeaponSet(_ schedules.Schedule) types.WeaponSchedule {
	if s.Schedule.Weapons[0].CoopSpecialWeapon != nil &&
		s.Schedule.Weapons[1].CoopSpecialWeapon != nil &&
		s.Schedule.Weapons[2].CoopSpecialWeapon != nil &&
		s.Schedule.Weapons[3].CoopSpecialWeapon != nil &&
		s.Schedule.Weapons[0].CoopSpecialWeapon.Name == RandomGrizzcoSchedule &&
		s.Schedule.Weapons[1].CoopSpecialWeapon.Name == RandomGrizzcoSchedule &&
		s.Schedule.Weapons[2].CoopSpecialWeapon.Name == RandomGrizzcoSchedule &&
		s.Schedule.Weapons[3].CoopSpecialWeapon.Name == RandomGrizzcoSchedule {
		return types.RandommGrizzco
	}
	if s.Schedule.Weapons[0].CoopSpecialWeapon != nil &&
		s.Schedule.Weapons[1].CoopSpecialWeapon != nil &&
		s.Schedule.Weapons[2].CoopSpecialWeapon != nil &&
		s.Schedule.Weapons[3].CoopSpecialWeapon != nil &&
		s.Schedule.Weapons[0].CoopSpecialWeapon.Name == RandomSchedule &&
		s.Schedule.Weapons[1].CoopSpecialWeapon.Name == RandomSchedule &&
		s.Schedule.Weapons[2].CoopSpecialWeapon.Name == RandomSchedule &&
		s.Schedule.Weapons[3].CoopSpecialWeapon.Name == RandomSchedule {
		return types.FourRandom
	}
	if s.Schedule.Weapons[0].Weapon != nil &&
		s.Schedule.Weapons[1].Weapon != nil &&
		s.Schedule.Weapons[2].Weapon != nil &&
		s.Schedule.Weapons[3].CoopSpecialWeapon != nil &&
		s.Schedule.Weapons[3].CoopSpecialWeapon.Name == RandomSchedule {
		return types.SingleRandom
	}
	if s.Schedule.Weapons[0].Weapon != nil &&
		s.Schedule.Weapons[1].Weapon != nil &&
		s.Schedule.Weapons[2].Weapon != nil &&
		s.Schedule.Weapons[3].Weapon != nil {
		return types.Set
	}
	return ""
}

// GetEvents implements lib.Shift{}.GetEvents() types.EventArr.
// This function returns a named type of slice of types.Event consisting of the events played in each wave of the shift.
func (s *ShiftSplatnet) GetEvents() types.EventArr {
	events := types.EventArr{}
	for i := range s.WaveDetails {
		events = append(events, s.WaveDetails[i].EventType.Key.ToEvent())
	}
	return events
}

// ScheduleInfo is the structure of a schedule retrieved from SplatNet.
// This isn't used anywhere in the library.
type ScheduleInfo struct {
	Schedules []struct {
		StartTime int `json:"start_time"`
		EndTime   int `json:"end_time"`
	} `json:"schedules"`
	Details []struct {
		StartTime int `json:"start_time"`
		Stage     struct {
			Name  string `json:"name"`
			Image string `json:"image"`
		} `json:"stage"`
		EndTime int `json:"end_time"`
		Weapons []struct {
			ID     string `json:"id"`
			Weapon struct {
				Thumbnail string `json:"thumbnail"`
				Name      string `json:"name"`
				Image     string `json:"image"`
				ID        string `json:"id"`
			} `json:"weapon"`
		} `json:"weapons"`
	} `json:"details"`
}

// GetStage implements lib.Shift{}.GetStage(schedules.Schedule) types.Stage.
// This function returns which stage the shift was played on.
func (s *ShiftSplatnet) GetStage(_ schedules.Schedule) types.Stage {
	switch s.Schedule.Stage.Name {
	case polaris:
		return types.RuinsOfArkPolaris
	case outpost:
		return types.LostOutpost
	case bay:
		return types.MaroonersBay
	case smokeyard:
		return types.SalmonidSmokeyard
	case grounds:
		return types.SpawningGrounds
	}
	return -1
}

// GetTides implements lib.Shift{}.GetTides() types.TideArr.
// This function returns a named type of slice of types.Tide consisting of the tides played in each wave of the shift.
func (s *ShiftSplatnet) GetTides() types.TideArr {
	tides := types.TideArr{}
	for i := range s.WaveDetails {
		tides = append(tides, s.WaveDetails[i].WaterLevel.Key.ToTide())
	}
	return tides
}

// GetEggsWaves implements lib.Shift{}.GetEggsWaves() []int.
// This function returns a slice of integers consisting of the amount of golden eggs delivered each wave.
func (s *ShiftSplatnet) GetEggsWaves() []int {
	eggs := []int{}
	for i := range s.WaveDetails {
		eggs = append(eggs, s.WaveDetails[i].GoldenEggs)
	}
	return eggs
}

// GetWaveCount implements lib.Shift{}.GetWaveCount() int.
// This function returns the number of waves in the shift.
func (s *ShiftSplatnet) GetWaveCount() int {
	return len(s.WaveDetails)
}

// GetTime implements lib.Shift{}.GetTime() time.Time
// This function returns the time at which the shift started.
func (s *ShiftSplatnet) GetTime() time.Time {
	return time.Unix(s.PlayTime, 0)
}
