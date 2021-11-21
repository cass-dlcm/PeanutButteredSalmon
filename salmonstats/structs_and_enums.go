package salmonstats

import (
	"fmt"
	"github.com/cass-dlcm/peanutbutteredsalmon/schedules"
	"github.com/cass-dlcm/peanutbutteredsalmon/types"
	"log"
	"time"
)

// ShiftPage is the payload as downloaded from SplatNet, containing all the shifts plus some additional data.
type ShiftPage struct {
	CurrentPage  int    `json:"current_page"`
	FirstPageURL string `json:"first_page_url"`
	From         int    `json:"from"`
	LastPage     int    `json:"last_page"`
	LastPageURL  string `json:"last_page_url"`
	Links        []struct {
		URL    interface{} `json:"url"`
		Label  interface{} `json:"label"`
		Active bool        `json:"active"`
	} `json:"links"`
	NextPageURL string             `json:"next_page_url"`
	Path        string             `json:"path"`
	PerPage     int                `json:"per_page"`
	PrevPageURL interface{}        `json:"prev_page_url"`
	To          int                `json:"to"`
	Total       int                `json:"total"`
	Results     []ShiftSalmonStats `json:"results"`
}

// Bosses is used in the ShiftSalmonStats struct but has no methods of its own.
type Bosses struct {
	Num3  int `json:"3"`
	Num6  int `json:"6"`
	Num9  int `json:"9"`
	Num12 int `json:"12"`
	Num13 int `json:"13"`
	Num14 int `json:"14"`
	Num15 int `json:"15"`
	Num16 int `json:"16"`
	Num21 int `json:"21"`
}

// ShiftSalmonStats is the data of a single shift from a salmon-stats/api instance.
// It implements the lib.Shift interface.
type ShiftSalmonStats struct {
	ID                         int    `json:"id"`
	ScheduleID                 string `json:"schedule_id"`
	PlayerID                   string
	Page                       int
	StartAt                    string    `json:"start_at"`
	Members                    []string  `json:"members"`
	BossAppearances            Bosses    `json:"boss_appearances"`
	UploaderUserID             int       `json:"uploader_user_id"`
	ClearWaves                 int       `json:"clear_waves"`
	FailReasonID               int       `json:"fail_reason_id"`
	DangerRate                 string    `json:"danger_rate"`
	CreatedAt                  time.Time `json:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at"`
	GoldenEggDelivered         int       `json:"golden_egg_delivered"`
	PowerEggCollected          int       `json:"power_egg_collected"`
	BossAppearanceCount        int       `json:"boss_appearance_count"`
	BossEliminationCount       int       `json:"boss_elimination_count"`
	IsEligibleForNoNightRecord bool      `json:"is_eligible_for_no_night_record"`
	PlayerResults              []struct {
		PlayerID             string      `json:"player_id"`
		GoldenEggs           int         `json:"golden_eggs"`
		PowerEggs            int         `json:"power_eggs"`
		Rescue               int         `json:"rescue"`
		Death                int         `json:"death"`
		SpecialID            int         `json:"special_id"`
		BossEliminationCount int         `json:"boss_elimination_count"`
		GradePoint           interface{} `json:"grade_point"`
		BossEliminations     struct {
			Counts Bosses `json:"counts"`
		} `json:"boss_eliminations"`
		SpecialUses []struct {
			Count int `json:"count"`
		} `json:"special_uses"`
		Weapons []struct {
			WeaponID int `json:"weapon_id"`
		} `json:"weapons"`
	} `json:"player_results"`
	Waves []struct {
		Wave                 int `json:"wave"`
		EventID              int `json:"event_id"`
		WaterID              int `json:"water_id"`
		GoldenEggQuota       int `json:"golden_egg_quota"`
		GoldenEggAppearances int `json:"golden_egg_appearances"`
		GoldenEggDelivered   int `json:"golden_egg_delivered"`
		PowerEggCollected    int `json:"power_egg_collected"`
	} `json:"waves"`
}

// GetTotalEggs implements lib.Shift{}.GetTotalEggs() int.
// This function returns the total eggs obtained in the shift.
func (s ShiftSalmonStats) GetTotalEggs() int {
	return s.GoldenEggDelivered
}

// GetStage implements lib.Shift{}.GetStage(schedules.Schedule) types.Stage.
// This function returns which stage the shift was played on.
func (s ShiftSalmonStats) GetStage(schedule schedules.Schedule) types.Stage {
	scheduleTime, err := time.Parse("2006-01-02 15:04:05", s.ScheduleID)
	if err != nil {
		log.Panicln(err)
	}
	for i := range schedule.Result {
		if schedule.Result[i].StartUtc.Equal(scheduleTime) {
			switch schedule.Result[i].Stage.Name {
			case "シェケナダム":
				return types.SpawningGrounds
			case "難破船ドン・ブラコ":
				return types.MaroonersBay
			case "トキシラズいぶし工房":
				return types.SalmonidSmokeyard
			case "朽ちた箱舟 ポラリス":
				return types.RuinsOfArkPolaris
			case "海上集落シャケト場":
				return types.LostOutpost
			}
		}
	}
	return -1
}

// GetWeaponSet implements lib.Shift{}.GetWeaponSet(schedules.Schedule) types.WeaponSchedule.
// This function returns what kind of weapon set was used in the shift.
func (s ShiftSalmonStats) GetWeaponSet(weaponSets schedules.Schedule) types.WeaponSchedule {
	scheduleTime, err := time.Parse("2006-01-02 15:04:05", s.ScheduleID)
	if err != nil {
		log.Panicln(err)
	}
	for i := range weaponSets.Result {
		if weaponSets.Result[i].StartUtc.Equal(scheduleTime) {
			if weaponSets.Result[i].Weapons[0].ID == -2 && weaponSets.Result[i].Weapons[1].ID == -2 && weaponSets.Result[i].Weapons[2].ID == -2 && weaponSets.Result[i].Weapons[3].ID == -2 {
				return types.RandommGrizzco
			}
			if weaponSets.Result[i].Weapons[0].ID == -1 && weaponSets.Result[i].Weapons[1].ID == -1 && weaponSets.Result[i].Weapons[2].ID == -1 && weaponSets.Result[i].Weapons[3].ID == -1 {
				return types.FourRandom
			}
			if weaponSets.Result[i].Weapons[0].ID >= 0 && weaponSets.Result[i].Weapons[1].ID >= 0 && weaponSets.Result[i].Weapons[2].ID >= 0 && weaponSets.Result[i].Weapons[3].ID == -1 {
				return types.SingleRandom
			}
			if weaponSets.Result[i].Weapons[0].ID >= 0 && weaponSets.Result[i].Weapons[1].ID >= 0 && weaponSets.Result[i].Weapons[2].ID >= 0 && weaponSets.Result[i].Weapons[3].ID >= 0 {
				return types.Set
			}
		}
	}
	return ""
}

// GetEvents implements lib.Shift{}.GetEvents() types.EventArr.
// This function returns a named type of slice of types.Event consisting of the events played in each wave of the shift.
func (s ShiftSalmonStats) GetEvents() types.EventArr {
	events := types.EventArr{}
	for i := range s.Waves {
		switch s.Waves[i].EventID {
		case 0:
			events = append(events, types.WaterLevels)
		case 1:
			events = append(events, types.CohockCharge)
		case 2:
			events = append(events, types.Fog)
		case 3:
			events = append(events, types.GoldieSeeking)
		case 4:
			events = append(events, types.Griller)
		case 5:
			events = append(events, types.Mothership)
		case 6:
			events = append(events, types.Rush)
		}
	}
	return events
}

// GetTides implements lib.Shift{}.GetTides() types.TideArr.
// This function returns a named type of slice of types.Tide consisting of the tides played in each wave of the shift.
func (s ShiftSalmonStats) GetTides() types.TideArr {
	tides := types.TideArr{}
	for i := range s.Waves {
		switch s.Waves[i].WaterID {
		case 1:
			tides = append(tides, types.Lt)
		case 2:
			tides = append(tides, types.Nt)
		case 3:
			tides = append(tides, types.Ht)
		}
	}
	return tides
}

// GetEggsWaves implements lib.Shift{}.GetEggsWaves() []int.
// This function returns a slice of integers consisting of the amount of golden eggs delivered each wave.
func (s ShiftSalmonStats) GetEggsWaves() []int {
	eggs := []int{}
	for i := range s.Waves {
		eggs = append(eggs, s.Waves[i].GoldenEggDelivered)
	}
	return eggs
}

// GetWaveCount implements lib.Shift{}.GetWaveCount() int.
// This function returns the number of waves in the shift.
func (s ShiftSalmonStats) GetWaveCount() int {
	return len(s.Waves)
}

// GetTime implements lib.Shift{}.GetTime() time.Time
// This function returns the time at which the shift started.
func (s ShiftSalmonStats) GetTime() time.Time {
	startTime, err := time.Parse("2006-01-02 15:04:05", s.StartAt)
	if err != nil {
		log.Panicln(err)
	}
	return startTime
}

// GetIdentifier implements lib.Shift{}.GetIdentifier() string.
// This function returns a unique URL pointing to this exact shift on salmon-stats-api.yuki-games.
func (s ShiftSalmonStats) GetIdentifier() string {
	return fmt.Sprintf("https://salmon-stats-api.yuki.games/api/results/%d/", s.ID)
}
