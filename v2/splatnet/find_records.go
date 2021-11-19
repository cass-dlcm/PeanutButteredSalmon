package splatnet

import (
	"github.com/cass-dlcm/PeanutButteredSalmon/v2/schedules"
	"github.com/cass-dlcm/PeanutButteredSalmon/v2/types"
	"time"
)

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

func (s *ShiftSplatnet) GetEvents() types.EventArr {
	events := types.EventArr{}
	for i := range s.WaveDetails {
		events = append(events, s.WaveDetails[i].EventType.Key.ToEvent())
	}
	return events
}

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

func (s *ShiftSplatnet) GetTides() types.TideArr {
	tides := types.TideArr{}
	for i := range s.WaveDetails {
		tides = append(tides, s.WaveDetails[i].WaterLevel.Key.ToTide())
	}
	return tides
}

func (s *ShiftSplatnet) GetEggsWaves() []int {
	eggs := []int{}
	for i := range s.WaveDetails {
		eggs = append(eggs, s.WaveDetails[i].GoldenEggs)
	}
	return eggs
}

func (s *ShiftSplatnet) GetWaveCount() int {
	return len(s.WaveDetails)
}

func (s *ShiftSplatnet) GetTime() time.Time {
	return time.Unix(s.PlayTime, 0)
}
