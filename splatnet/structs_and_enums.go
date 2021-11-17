package splatnet

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

import "PeanutButteredSalmon/enums"

type ShiftList struct {
	Code    *string `json:"code"`
	Summary struct {
		Card struct {
			GoldenIkuraTotal int `json:"golden_ikura_total"`
			HelpTotal        int `json:"help_total"`
			KumaPointTotal   int `json:"kuma_point_total"`
			IkuraTotal       int `json:"ikura_total"`
			KumaPoint        int `json:"kuma_point"`
			JobNum           int `json:"job_num"`
		} `json:"card"`
		Stats []struct {
			DeadTotal            int   `json:"dead_total"`
			MyGoldenIkuraTotal   int   `json:"my_golden_ikura_total"`
			GradePoint           int   `json:"grade_point"`
			TeamGoldenIkuraTotal int   `json:"team_golden_ikura_total"`
			HelpTotal            int   `json:"help_total"`
			TeamIkuraTotal       int   `json:"team_ikura_total"`
			StartTime            int   `json:"start_time"`
			MyIkuraTotal         int   `json:"my_ikura_total"`
			FailureCounts        []int `json:"failure_counts"`
			Schedule             struct {
				Stage struct {
					Image string `json:"image"`
					Name  string `json:"name"`
				} `json:"stage"`
				EndTime   int `json:"end_time"`
				StartTime int `json:"start_time"`
				Weapons   []struct {
					Weapon struct {
						ID        string `json:"id"`
						Image     string `json:"image"`
						Name      string `json:"name"`
						Thumbnail string `json:"thumbnail"`
					} `json:"weapon"`
					ID string `json:"id"`
				} `json:"weapons"`
			} `json:"schedule"`
			JobNum         int `json:"job_num"`
			KumaPointTotal int `json:"kuma_point_total"`
			EndTime        int `json:"end_time"`
			ClearNum       int `json:"clear_num"`
			Grade          struct {
				Name string `json:"name"`
				ID   string `json:"id"`
			} `json:"grade"`
		} `json:"stats"`
	} `json:"summary"`
	RewardGear struct {
		Thumbnail string `json:"thumbnail"`
		Kind      string `json:"kind"`
		ID        string `json:"id"`
		Name      string `json:"name"`
		Brand     struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Image string `json:"image"`
		} `json:"brand"`
		Rarity int    `json:"rarity"`
		Image  string `json:"image"`
	} `json:"reward_gear"`
	Results []ShiftSplatnet `json:"results"`
}

func (s *ShiftSplatnet) TotalEggs() int {
	sum := 0
	for i := range s.WaveDetails {
		sum += s.WaveDetails[i].GoldenEggs
	}
	return sum
}

type ShiftSplatnet struct {
	JobId           int64                   `json:"job_id"`
	DangerRate      float64                 `json:"danger_rate"`
	JobResult       ShiftSplatnetJobResult  `json:"job_result"`
	JobScore        int                     `json:"job_score"`
	JobRate         int                     `json:"job_rate"`
	GradePoint      int                     `json:"grade_point"`
	GradePointDelta int                     `json:"grade_point_delta"`
	OtherResults    []ShiftSplatnetPlayer   `json:"other_results"`
	KumaPoint       int                     `json:"kuma_point"`
	StartTime       int64                   `json:"start_time"`
	PlayerType      SplatnetPlayerType      `json:"player_type"`
	PlayTime        int64                   `json:"play_time"`
	BossCounts      ShiftSplatnetBossCounts `json:"boss_counts"`
	EndTime         int64                   `json:"end_time"`
	MyResult        ShiftSplatnetPlayer     `json:"my_result"`
	WaveDetails     []ShiftSplatnetWave     `json:"wave_details"`
	Grade           ShiftSplatnetGrade      `json:"grade"`
	Schedule        ShiftSplatnetSchedule   `json:"schedule"`
}

type ShiftSplatnetJobResult struct {
	IsClear       bool                 `json:"is_clear,omitempty"`
	FailureReason *enums.FailureReason `json:"failure_reason,omitempty"`
	FailureWave   *int                 `json:"failure_wave,omitempty"`
}

type ShiftSplatnetPlayer struct {
	SpecialCounts  []int                           `json:"special_counts"`
	Special        SplatnetQuad                    `json:"special"`
	Pid            string                          `json:"pid"`
	PlayerType     SplatnetPlayerType              `json:"player_type"`
	WeaponList     []ShiftSplatnetPlayerWeaponList `json:"weapon_list"`
	Name           string                          `json:"name"`
	DeadCount      int                             `json:"dead_count"`
	GoldenEggs     int                             `json:"golden_ikura_num"`
	BossKillCounts ShiftSplatnetBossCounts         `json:"boss_kill_counts"`
	PowerEggs      int                             `json:"ikura_num"`
	HelpCount      int                             `json:"help_count"`
}

type ShiftSplatnetPlayerWeaponList struct {
	Id     string                              `json:"id"`
	Weapon ShiftSplatnetPlayerWeaponListWeapon `json:"weapon"`
}

type ShiftSplatnetPlayerWeaponListWeapon struct {
	Id        weapon `json:"id"`
	Image     string `json:"image"`
	Name      string `json:"name"`
	Thumbnail string `json:"thumbnail"`
}

type ShiftSplatnetBossCounts struct {
	Goldie    ShiftSplatnetBossCountsBoss `json:"3"`
	Steelhead ShiftSplatnetBossCountsBoss `json:"6"`
	Flyfish   ShiftSplatnetBossCountsBoss `json:"9"`
	Scrapper  ShiftSplatnetBossCountsBoss `json:"12"`
	SteelEel  ShiftSplatnetBossCountsBoss `json:"13"`
	Stinger   ShiftSplatnetBossCountsBoss `json:"14"`
	Maws      ShiftSplatnetBossCountsBoss `json:"15"`
	Griller   ShiftSplatnetBossCountsBoss `json:"16"`
	Drizzler  ShiftSplatnetBossCountsBoss `json:"21"`
}

type ShiftSplatnetBossCountsBoss struct {
	Boss  SplatnetDouble `json:"boss"`
	Count int            `json:"count"`
}

type ShiftSplatnetWave struct {
	WaterLevel   WaterLevels `json:"water_level"`
	EventType    eventStruct `json:"event_type"`
	GoldenEggs   int         `json:"golden_ikura_num"`
	GoldenAppear int         `json:"golden_ikura_pop_num"`
	PowerEggs    int         `json:"ikura_num"`
	QuotaNum     int         `json:"quota_num"`
}

type ShiftSplatnetGrade struct {
	Id        string `json:"id,omitempty"`
	ShortName string `json:"short_name,omitempty"`
	LongName  string `json:"long_name,omitempty"`
	Name      string `json:"name,omitempty"`
}

type ShiftSplatnetSchedule struct {
	StartTime int64                         `json:"start_time"`
	Weapons   []ShiftSplatnetScheduleWeapon `json:"weapons"`
	EndTime   int64                         `json:"end_time"`
	Stage     ShiftSplatnetScheduleStage    `json:"stage"`
}

type ShiftSplatnetScheduleWeapon struct {
	Id                string                                    `json:"id"`
	Weapon            *ShiftSplatnetScheduleWeaponWeapon        `json:"weapon"`
	CoopSpecialWeapon *ShiftSplatnetScheduleWeaponSpecialWeapon `json:"coop_special_weapon"`
}

type ShiftSplatnetScheduleWeaponWeapon struct {
	Id        string         `json:"id"`
	Image     string         `json:"image"`
	Name      weaponSchedule `json:"name"`
	Thumbnail string         `json:"thumbnail"`
}

type ShiftSplatnetScheduleWeaponSpecialWeapon struct {
	Image string         `json:"image"`
	Name  weaponSchedule `json:"name"`
}

type ShiftSplatnetScheduleStage struct {
	Image splatnetScheduleStageImage `json:"image"`
	Name  stageName                  `json:"name"`
}

type SplatnetTriple struct {
	Id    string `json:"id"`
	Image string `json:"image"`
	Name  string `json:"name"`
}

type SplatnetDouble struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

type WaterLevels struct {
	Key  tide   `json:"key"`
	Name string `json:"name"`
}

type eventStruct struct {
	Key  event  `json:"key"`
	Name string `json:"name"`
}

type SplatnetPlayerType struct {
	Gender  gender  `json:"style,omitempty"`
	Species species `json:"species,omitempty"`
}

type SplatnetQuad struct {
	Id     string `json:"id"`
	ImageA string `json:"image_a"`
	ImageB string `json:"image_b"`
	Name   string `json:"name"`
}

type stageName string

const (
	smokeyard stageName = "Salmonid Smokeyard"
	polaris   stageName = "Ruins of Ark Polaris"
	grounds   stageName = "Spawning Grounds"
	bay       stageName = "Marooner's Bay"
	outpost   stageName = "Lost Outpost"
)

type splatnetScheduleStageImage string

const (
	smokeyardSplatnetImg splatnetScheduleStageImage = "/images/coop_stage/e9f7c7b35e6d46778cd3cbc0d89bd7e1bc3be493.png"
	polarisSplatnetImg   splatnetScheduleStageImage = "/images/coop_stage/50064ec6e97aac91e70df5fc2cfecf61ad8615fd.png"
	groundsSplatnetImg   splatnetScheduleStageImage = "/images/coop_stage/65c68c6f0641cc5654434b78a6f10b0ad32ccdee.png"
	baySplatnetImg       splatnetScheduleStageImage = "/images/coop_stage/e07d73b7d9f0c64e552b34a2e6c29b8564c63388.png"
	outpostSplatnetImg   splatnetScheduleStageImage = "/images/coop_stage/6d68f5baa75f3a94e5e9bfb89b82e7377e3ecd2c.png"
)

func (sN *stageName) UnmarshalJSON(b []byte) error {
	// Define a secondary type to avoid ending up with a recursive call to json.Unmarshal
	type SSE stageName
	r := (*SSE)(sN)
	err := json.Unmarshal(b, &r)
	if err != nil {
		panic(err)
	}
	switch *sN {
	case smokeyard, polaris, grounds, bay, outpost:
		return nil
	}
	return errors.New("invalid type")
}

func (ssssie *splatnetScheduleStageImage) UnmarshalJSON(b []byte) error {
	// Define a secondary type to avoid ending up with a recursive call to json.Unmarshal
	type SSSSIE splatnetScheduleStageImage
	r := (*SSSSIE)(ssssie)
	err := json.Unmarshal(b, &r)
	if err != nil {
		panic(err)
	}
	switch *ssssie {
	case smokeyardSplatnetImg, polarisSplatnetImg, groundsSplatnetImg, baySplatnetImg, outpostSplatnetImg:
		return nil
	}
	return errors.New("invalid type")
}

type gender string

type species string

func getSplatnetStage(s enums.Stage) stageName {
	switch s {
	case enums.LostOutpost:
		return outpost
	case enums.MaroonersBay:
		return bay
	case enums.SpawningGrounds:
		return grounds
	case enums.SalmonidSmokeyard:
		return smokeyard
	case enums.RuinsOfArkPolaris:
		return polaris
	}
	return ""
}

func getSplatnetEvent(e enums.Event) event {
	switch e {
	case enums.Griller:
		return griller
	case enums.Fog:
		return fog
	case enums.CohockCharge:
		return cohockCharge
	case enums.GoldieSeeking:
		return goldieSeeking
	case enums.Mothership:
		return mothership
	case enums.WaterLevels:
		return waterLevels
	case enums.Rush:
		return rush
	}
	return ""
}

func (sN *stageName) isElementExist(arr []stageName) bool {
	for _, v := range arr {
		if v == *sN {
			return true
		}
	}
	return false
}

type event string

const (
	griller       event = "griller"
	fog           event = "fog"
	cohockCharge  event = "cohock-charge"
	goldieSeeking event = "goldie-seeking"
	mothership    event = "the-mothership"
	waterLevels   event = "water-levels"
	rush          event = "rush"
)

func (e *event) UnmarshalJSON(b []byte) error {
	// Define a secondary type to avoid ending up with a recursive call to json.Unmarshal
	type E event
	r := (*E)(e)
	err := json.Unmarshal(b, &r)
	if err != nil {
		panic(err)
	}
	switch *e {
	case griller, fog, cohockCharge, goldieSeeking, mothership, waterLevels, rush:
		return nil
	}
	return errors.New("Invalid event. Got: " + fmt.Sprint(e))
}

func (e *event) isElementExist(arr []event) bool {
	for _, v := range arr {
		if v == *e {
			return true
		}
	}
	return false
}

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

type Record struct {
	Time         time.Time
	RecordAmount int
	Shift        ShiftSplatnet
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

func getStageNames() []stageName {
	return []stageName{
		smokeyard,
		polaris,
		grounds,
		bay,
		outpost,
	}
}

func getAllRecords() map[RecordName]map[stageName]map[enums.WeaponSchedule]*Record {
	records := map[RecordName]map[stageName]map[enums.WeaponSchedule]*Record{}
	recordNames := getRecordNames()
	stageNames := getStageNames()
	weapons := enums.GetAllWeapons()
	for i := range recordNames {
		records[recordNames[i]] = map[stageName]map[enums.WeaponSchedule]*Record{}
		for j := range stageNames {
			records[recordNames[i]][stageNames[j]] = map[enums.WeaponSchedule]*Record{}
			for k := range weapons {
				records[recordNames[i]][stageNames[j]][weapons[k]] = &Record{}
			}
		}
	}
	return records
}

type tide string

const (
	ht tide = "high"
	lt tide = "low"
	nt tide = "normal"
)

func (t *tide) UnmarshalJSON(b []byte) error {
	// Define a secondary type to avoid ending up with a recursive call to json.Unmarshal
	type T tide
	r := (*T)(t)
	err := json.Unmarshal(b, &r)
	if err != nil {
		panic(err)
	}
	switch *t {
	case ht, lt, nt:
		return nil
	}
	return errors.New("Invalid tide. Got: " + fmt.Sprint(t))
}

func (t *tide) isElementExist(arr []tide) bool {
	for _, v := range arr {
		if v == *t {
			return true
		}
	}
	return false
}

func getSplatnetTide(t enums.Tide) tide {
	switch t {
	case enums.Lt:
		return lt
	case enums.Nt:
		return nt
	case enums.Ht:
		return ht
	}
	return ""
}

type weapon string

// List of weapon
const (
	SalmonSplooshOMatic      weapon = "0"
	SalmonSplattershotJr     weapon = "10"
	SalmonSplashOMatic       weapon = "20"
	SalmonAerosprayMg        weapon = "30"
	SalmonSplattershot       weapon = "40"
	Salmon52gal              weapon = "50"
	SalmonNZap85             weapon = "60"
	SalmonSplattershotPro    weapon = "70"
	Salmon96gal              weapon = "80"
	SalmonJetSquelcher       weapon = "90"
	SalmonLunaBlaster        weapon = "200"
	SalmonBlaster            weapon = "210"
	SalmonRangeBlaster       weapon = "220"
	SalmonClashBlaster       weapon = "230"
	SalmonRapidBlaster       weapon = "240"
	SalmonRapidBlasterPro    weapon = "250"
	SalmonL3Nozzlenose       weapon = "300"
	SalmonH3Nozzlenose       weapon = "310"
	SalmonSqueezer           weapon = "400"
	SalmonCarbonRoller       weapon = "1000"
	SalmonSplatRoller        weapon = "1010"
	SalmonDynamoRoller       weapon = "1020"
	SalmonFlingzaRoller      weapon = "1030"
	SalmonInkbrush           weapon = "1100"
	SalmonOctobrush          weapon = "1110"
	SalmonClassicSquiffer    weapon = "2000"
	SalmonSplatCharger       weapon = "2010"
	SalmonSplatterscope      weapon = "2020"
	SalmonELiter4K           weapon = "2030"
	SalmonELiter4KScope      weapon = "2040"
	SalmonBamboozler14MkI    weapon = "2050"
	SalmonGooTuber           weapon = "2060"
	SalmonSlosher            weapon = "3000"
	SalmonTriSlosher         weapon = "3010"
	SalmonSloshingMachine    weapon = "3020"
	SalmonBloblobber         weapon = "3030"
	SalmonExplosher          weapon = "3040"
	SalmonMiniSplatling      weapon = "4000"
	SalmonHeavySplatling     weapon = "4010"
	SalmonHydraSplatling     weapon = "4020"
	SalmonBallpointSplatling weapon = "4030"
	SalmonNautilus47         weapon = "4040"
	SalmonDappleDualies      weapon = "5000"
	SalmonSplatDualies       weapon = "5010"
	SalmonGloogaDualies      weapon = "5020"
	SalmonDualieSquelchers   weapon = "5030"
	SalmonDarkTetraDualies   weapon = "5040"
	SalmonSplatBrella        weapon = "6000"
	SalmonTentaBrella        weapon = "6010"
	SalmonUndercoverBrella   weapon = "6020"
	SalmonGrizzcoBlaster     weapon = "20000"
	SalmonGrizzcoBrella      weapon = "20010"
	SalmonGrizzcoCharger     weapon = "20020"
	SalmonGrizzcoSlosher     weapon = "20030"
)

func (swe *weapon) UnmarshalJSON(b []byte) error {
	// Define a secondary type to avoid ending up with a recursive call to json.Unmarshal
	type SWE weapon
	r := (*SWE)(swe)
	err := json.Unmarshal(b, &r)
	if err != nil {
		panic(err)
	}
	switch *swe {
	case SalmonSplooshOMatic, SalmonSplattershotJr, SalmonSplashOMatic, SalmonAerosprayMg, SalmonSplattershot,
		Salmon52gal, SalmonNZap85, SalmonSplattershotPro, Salmon96gal, SalmonJetSquelcher, SalmonLunaBlaster, SalmonBlaster,
		SalmonRangeBlaster, SalmonClashBlaster, SalmonRapidBlaster, SalmonRapidBlasterPro, SalmonL3Nozzlenose,
		SalmonH3Nozzlenose, SalmonSqueezer, SalmonCarbonRoller, SalmonSplatRoller, SalmonDynamoRoller, SalmonFlingzaRoller,
		SalmonInkbrush, SalmonOctobrush, SalmonClassicSquiffer, SalmonSplatCharger, SalmonSplatterscope, SalmonELiter4K,
		SalmonELiter4KScope, SalmonBamboozler14MkI, SalmonGooTuber, SalmonSlosher, SalmonTriSlosher,
		SalmonSloshingMachine, SalmonBloblobber, SalmonExplosher, SalmonMiniSplatling, SalmonHeavySplatling,
		SalmonHydraSplatling, SalmonBallpointSplatling, SalmonNautilus47, SalmonDappleDualies, SalmonSplatDualies,
		SalmonGloogaDualies, SalmonDualieSquelchers, SalmonDarkTetraDualies, SalmonSplatBrella, SalmonTentaBrella,
		SalmonUndercoverBrella, SalmonGrizzcoBlaster, SalmonGrizzcoBrella, SalmonGrizzcoCharger, SalmonGrizzcoSlosher:
		return nil
	}
	return errors.New("Invalid weapon. Got: " + fmt.Sprint(*swe))
}

type weaponSchedule string

// List of weaponSchedule
const (
	RandomGrizzcoSchedule      weaponSchedule = "Random Grizzco"
	RandomSchedule             weaponSchedule = "Random"
	SplooshOMaticSchedule      weaponSchedule = "Sploosh-o-matic"
	SplattershotJrSchedule     weaponSchedule = "Splattershot Jr."
	SplashOMaticSchedule       weaponSchedule = "Splash-o-matic"
	AerosprayMgSchedule        weaponSchedule = "Aerospray MG"
	SplattershotSchedule       weaponSchedule = "Splattershot"
	Point52GalSchedule         weaponSchedule = ".52 Gal"
	NZap85Schedule             weaponSchedule = "N-ZAP '85"
	SplattershotProSchedule    weaponSchedule = "Splattershot Pro"
	Point96GalSchedule         weaponSchedule = ".96 Gal"
	JetSquelcherSchedule       weaponSchedule = "Jet Squelcher"
	LunaBlasterSchedule        weaponSchedule = "Luna Blaster"
	BlasterSchedule            weaponSchedule = "Blaster"
	RangeBlasterSchedule       weaponSchedule = "Range Blaster"
	ClashBlasterSchedule       weaponSchedule = "Clash Blaster"
	RapidBlasterSchedule       weaponSchedule = "Rapid Blaster"
	RapidBlasterProSchedule    weaponSchedule = "Rapid Blaster Pro"
	L3NozzlenoseSchedule       weaponSchedule = "L-3 Nozzlenose"
	H3NozzlenoseSchedule       weaponSchedule = "H-3 Nozzlenose"
	SqueezerSchedule           weaponSchedule = "Squeezer"
	CarbonRollerSchedule       weaponSchedule = "Carbon Roller"
	SplatRollerSchedule        weaponSchedule = "Splat Roller"
	DynamoRollerSchedule       weaponSchedule = "Dynamo Roller"
	FlingzaRollerSchedule      weaponSchedule = "Flingza Roller"
	InkbrushSchedule           weaponSchedule = "Inkbrush"
	OctobrushSchedule          weaponSchedule = "Octobrush"
	ClassicSquifferSchedule    weaponSchedule = "Classic Squiffer"
	SplatChargerSchedule       weaponSchedule = "Splat Charger"
	SplatterscopeSchedule      weaponSchedule = "Splatterscope"
	ELiter4KSchedule           weaponSchedule = "E-liter 4K"
	ELiter4KScopeSchedule      weaponSchedule = "E-liter 4K Scope"
	Bamboozler14MkISchedule    weaponSchedule = "Bamboozler 14 Mk I"
	GooTuberSchedule           weaponSchedule = "Goo Tuber"
	SlosherSchedule            weaponSchedule = "Slosher"
	SodaSlosherSchedule        weaponSchedule = "Soda Slosher"
	TriSlosherSchedule         weaponSchedule = "Tri-Slosher"
	SloshingMachineSchedule    weaponSchedule = "Sloshing Machine"
	BloblobberSchedule         weaponSchedule = "Bloblobber"
	ExplosherSchedule          weaponSchedule = "Explosher"
	MiniSplatlingSchedule      weaponSchedule = "Mini Splatling"
	HeavySplatlingSchedule     weaponSchedule = "Heavy Splatling"
	HydraSplatlingSchedule     weaponSchedule = "Hydra Splatling"
	BallpointSplatlingSchedule weaponSchedule = "Ballpoint Splatling"
	Nautilus47Schedule         weaponSchedule = "Nautilus 47"
	DappleDualiesSchedule      weaponSchedule = "Dapple Dualies"
	SplatDualiesSchedule       weaponSchedule = "Splat Dualies"
	GloogaDualiesSchedule      weaponSchedule = "Glooga Dualies"
	DualieSquelchersSchedule   weaponSchedule = "Dualie Squelchers"
	DarkTetraDualiesSchedule   weaponSchedule = "Dark Tetra Dualies"
	SplatBrellaSchedule        weaponSchedule = "Splat Brella"
	TentaBrellaSchedule        weaponSchedule = "Tenta Brella"
	UndercoverBrellaSchedule   weaponSchedule = "Undercover Brella"
)

func (swse *weaponSchedule) UnmarshalJSON(b []byte) error {
	// Define a secondary type to avoid ending up with a recursive call to json.Unmarshal
	type SWSE weaponSchedule
	r := (*SWSE)(swse)
	err := json.Unmarshal(b, &r)
	if err != nil {
		panic(err)
	}
	switch *swse {
	case RandomGrizzcoSchedule, RandomSchedule, SplooshOMaticSchedule, SplattershotJrSchedule, SplashOMaticSchedule,
		AerosprayMgSchedule, SplattershotSchedule, Point52GalSchedule, NZap85Schedule, SplattershotProSchedule,
		Point96GalSchedule, JetSquelcherSchedule, LunaBlasterSchedule, BlasterSchedule, RangeBlasterSchedule,
		ClashBlasterSchedule, RapidBlasterSchedule, RapidBlasterProSchedule, L3NozzlenoseSchedule, H3NozzlenoseSchedule,
		SqueezerSchedule, CarbonRollerSchedule, SplatRollerSchedule, DynamoRollerSchedule, FlingzaRollerSchedule,
		InkbrushSchedule, OctobrushSchedule, ClassicSquifferSchedule, SplatChargerSchedule, SplatterscopeSchedule,
		ELiter4KSchedule, ELiter4KScopeSchedule, Bamboozler14MkISchedule, GooTuberSchedule, SlosherSchedule,
		SodaSlosherSchedule, TriSlosherSchedule, SloshingMachineSchedule, BloblobberSchedule, ExplosherSchedule,
		MiniSplatlingSchedule, HeavySplatlingSchedule, HydraSplatlingSchedule, BallpointSplatlingSchedule,
		Nautilus47Schedule, DappleDualiesSchedule, SplatDualiesSchedule, GloogaDualiesSchedule,
		DualieSquelchersSchedule, DarkTetraDualiesSchedule, SplatBrellaSchedule, TentaBrellaSchedule,
		UndercoverBrellaSchedule:
		return nil
	}
	return errors.New("Invalid weaponSchedule. Got: " + fmt.Sprint(*swse))
}
