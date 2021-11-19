package statink

import (
	"fmt"
	"github.com/cass-dlcm/PeanutButteredSalmon/schedules"
	"github.com/cass-dlcm/PeanutButteredSalmon/types"
	"log"
	"time"
)

type ShiftStatInk struct {
	ID             int    `json:"id"`
	UUID           string `json:"uuid"`
	SplatnetNumber int    `json:"splatnet_number"`
	URL            string `json:"url"`
	APIEndpoint string        `json:"api_endpoint"`
	Stage       StatinkTriple `json:"stage"`
	IsCleared   *bool          `json:"is_cleared"`
	FailReason *struct {
		Key  string `json:"key"`
		Name Name   `json:"name"`
	} `json:"fail_reason"`
	ClearWaves      *int    `json:"clear_waves"`
	DangerRate      interface{} `json:"danger_rate"`
	Quota           []int  `json:"quota"`
	Title           *Title  `json:"title"`
	TitleExp        *int    `json:"title_exp"`
	TitleAfter      *Title  `json:"title_after"`
	TitleExpAfter   *int    `json:"title_exp_after"`
	BossAppearances []BossCount `json:"boss_appearances"`
	Waves []struct {
		KnownOccurrence *StatinkTriple `json:"known_occurrence"`
		WaterLevel      *StatinkTriple `json:"water_level"`
		GoldenEggQuota  *int           `json:"golden_egg_quota"`
		GoldenEggAppearances *int            `json:"golden_egg_appearances"`
		GoldenEggDelivered   *int            `json:"golden_egg_delivered"`
		PowerEggCollected    *int            `json:"power_egg_collected"`
	} `json:"waves"`
	MyData    *Player `json:"my_data"`
	Teammates []Player `json:"teammates"`
	Agent     struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"agent"`
	Automated    bool    `json:"automated"`
	Note         *string `json:"note"`
	LinkURL      *string `json:"link_url"`
	ShiftStartAt *Time    `json:"shift_start_at"`
	StartAt      *Time    `json:"start_at"`
	EndAt        *Time   `json:"end_at"`
	RegisterAt   Time    `json:"register_at"`
}

type StatinkTriple struct {
	Key      string `json:"key"`
	Splatnet *string    `json:"splatnet"`
	Name     Name   `json:"name"`
}

type Title struct {
	Key         string `json:"key"`
	Splatnet    *int    `json:"splatnet"`
	Name        Name   `json:"name"`
	GenericName Name   `json:"generic_name"`
}

type Time struct {
	Time    int64     `json:"time"`
	Iso8601 time.Time `json:"iso8601"`
}

type Boss struct {
	Key         string `json:"key"`
	Splatnet    *int    `json:"splatnet"`
	SplatnetStr string `json:"splatnet_str"`
	Name        Name   `json:"name"`
}

type BossCount struct {
	Boss  Boss `json:"boss"`
	Count int  `json:"count"`
}

type Player struct {
	SplatnetID string `json:"splatnet_id"`
	Name       string `json:"name"`
	Special    struct {
		Key  string `json:"key"`
		Name Name   `json:"name"`
	} `json:"special"`
	Rescue             int `json:"rescue"`
	Death              int `json:"death"`
	GoldenEggDelivered int `json:"golden_egg_delivered"`
	PowerEggCollected  int `json:"power_egg_collected"`
	Species            struct {
		Key  string `json:"key"`
		Name Name   `json:"name"`
	} `json:"species"`
	Gender struct {
		Key     string `json:"key"`
		Iso5218 int    `json:"iso5218"`
		Name    Name   `json:"name"`
	} `json:"gender"`
	SpecialUses []int `json:"special_uses"`
	Weapons     []struct {
		Key      string `json:"key"`
		Splatnet int    `json:"splatnet"`
		Name     Name   `json:"name"`
	} `json:"weapons"`
	BossKills []BossCount `json:"boss_kills"`
}

type Name struct {
	DeDE string `json:"de_DE"`
	EnGB string `json:"en_GB"`
	EnUS string `json:"en_US"`
	EsES string `json:"es_ES"`
	EsMX string `json:"es_MX"`
	FrCA string `json:"fr_CA"`
	FrFR string `json:"fr_FR"`
	ItIT string `json:"it_IT"`
	JaJP string `json:"ja_JP"`
	NlNL string `json:"nl_NL"`
	RuRU string `json:"ru_RU"`
	ZhCN string `json:"zh_CN"`
	ZhTW string `json:"zh_TW"`
}

func (s *ShiftStatInk) GetTotalEggs() int {
	sum := 0
	for i := range s.Waves {
		sum += *s.Waves[i].GoldenEggDelivered
	}
	return sum
}

func (s *ShiftStatInk) GetStage(_ schedules.Schedule) types.Stage {
	switch s.Stage.Key {
	case "dam":
		return types.SpawningGrounds
	case "donburako":
		return types.MaroonersBay
	case "polaris":
		return types.RuinsOfArkPolaris
	case "shaketoba":
		return types.LostOutpost
	case "tokishirazu":
		return types.SalmonidSmokeyard
	}
	return -1
}

func (s *ShiftStatInk) GetWeaponSet(weaponSets schedules.Schedule) types.WeaponSchedule {
	for i := range weaponSets.Result {
		if weaponSets.Result[i].StartUtc.Equal(s.ShiftStartAt.Iso8601) {
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

func (s *ShiftStatInk) GetEvents() types.EventArr {
	events := types.EventArr{}
	for i := range s.Waves {
		if s.Waves[i].KnownOccurrence == nil {
			events = append(events, types.WaterLevels)
		} else {
			event, err := types.StringToEvent(s.Waves[i].KnownOccurrence.Key)
			if err != nil {
				log.Panic(err)
			}
			events = append(events, *event)
		}
	}
	return events
}

func (s *ShiftStatInk) GetTides() types.TideArr {
	tides := types.TideArr{}
	for i := range s.Waves {
		switch s.Waves[i].WaterLevel.Key {
		case "low":
			tides = append(tides, types.Lt)
		case "normal":
			tides = append(tides, types.Nt)
		case "high":
			tides = append(tides, types.Ht)
		}
	}
	return tides
}

func (s *ShiftStatInk) GetEggsWaves() []int {
	eggs := []int{}
	for i := range s.Waves {
		eggs = append(eggs, *s.Waves[i].GoldenEggDelivered)
	}
	return eggs
}

func (s *ShiftStatInk) GetWaveCount() int {
	return len(s.Waves)
}

func (s *ShiftStatInk) GetTime() time.Time {
	return s.StartAt.Iso8601
}

func (s *ShiftStatInk) GetIdentifier() string {
	return fmt.Sprintf("https://stat.ink/api/v2/salmon/%d", s.ID)
}