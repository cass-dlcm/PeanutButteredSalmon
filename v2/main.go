package main

import (
	"flag"
	"fmt"
	"github.com/cass-dlcm/PeanutButteredSalmon/v2/lib"
	"github.com/cass-dlcm/PeanutButteredSalmon/v2/salmonstats"
	"github.com/cass-dlcm/PeanutButteredSalmon/v2/splatnet"
	"github.com/cass-dlcm/PeanutButteredSalmon/v2/statink"
	"github.com/cass-dlcm/PeanutButteredSalmon/v2/types"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func setLanguage() {
	log.Println("Please enter your locale (see readme for list).")

	var locale string
	// Taking input from user
	if _, err := fmt.Scanln(&locale); err != nil {
		log.Panic(err)
	}
	languageList := map[string]string{
		"en-US": "en-US",
		"es-MX": "es-MX",
		"fr-CA": "fr-CA",
		"ja-JP": "ja-JP",
		"en-GB": "en-GB",
		"es-ES": "es-ES",
		"fr-FR": "fr-FR",
		"de-DE": "de-DE",
		"it-IT": "it-IT",
		"nl-NL": "nl-NL",
		"ru-RU": "ru-RU",
	}
	_, exists := languageList[locale]
	for !exists {
		log.Println("Invalid language code. Please try entering it again.")

		if _, err := fmt.Scanln(&locale); err != nil {
			panic(err)
		}

		_, exists = languageList[locale]
	}
	viper.Set("user_lang", locale)

	if err := viper.WriteConfig(); err != nil {
		panic(err)
	}
}

func getFlags() ([]types.Stage, []types.Event, []types.Tide, []types.WeaponSchedule, []types.Server, bool, []types.Server, int) {
	stagesStr := flag.String("stage", "spawning_grounds marooners_bay lost_outpost salmonid_smokeyard ruins_of_ark_polaris", "To set a specific set of stages.")
	hasEventsStr := flag.String("event", "water_levels rush fog goldie_seeking griller cohock_charge mothership", "To set a specific set of events.")
	hasTides := flag.String("tide", "LT NT HT", "To set a specific set of tides.")
	hasWeapons := flag.String("weapon", "set single_random four_random random_gold", "To restrict to a specific set of weapon types.")
	statInk := flag.String("statink", "", "To read data from stat.ink. Use \"official\" for the server at stat.ink.")
	useSplatnet := flag.Bool("splatnet", false, "To read data from splatnet.")
	salmonStats := flag.String("salmonstats", "", "To read data from salmon-stats. Use \"official\" for the server at salmon-stats-api.yuki.games")
	m := flag.Int("monitor", -1, "To monitor for new personal bests.")
	flag.Parse()

	stages, err := types.GetStageArgs(*stagesStr)
	if err != nil {
		log.Panicln(err)
	}
	hasEvents, err := types.GetEventArgs(*hasEventsStr)
	if err != nil {
		log.Panicln(err)
	}
	weapons, err := types.GetWeaponArgs(*hasWeapons)
	if err != nil {
		log.Panicln(err)
	}

	tides, err := types.GetTideArgs(*hasTides)
	if err != nil {
		panic(err)
	}

	statInkUrlNicks := strings.Split(*statInk, " ")
	var statInkUrlConf []types.Server
	if err := viper.UnmarshalKey("statink_servers", &statInkUrlConf); err != nil {
		log.Panicln(err)
	}
	statInkServers := []types.Server{}
	for i := range statInkUrlNicks {
		for j := range statInkUrlConf {
			if statInkUrlConf[j].ShortName == statInkUrlNicks[i] {
				statInkServers = append(statInkServers, statInkUrlConf[j])
			}
		}
	}

	salmonStatsUrlNicks := strings.Split(*salmonStats, " ")
	var salmonStatsUrlConf []types.Server
	if err := viper.UnmarshalKey("salmonstats_servers", &salmonStatsUrlConf); err != nil {
		log.Panicln(err)
	}
	salmonStatsServers := []types.Server{}
	for i := range salmonStatsUrlNicks {
		for j := range salmonStatsUrlConf {
			if salmonStatsUrlConf[j].ShortName == salmonStatsUrlNicks[i] {
				salmonStatsServers = append(salmonStatsServers, salmonStatsUrlConf[j])
			}
		}
	}

	return stages, hasEvents, tides, weapons, statInkServers, *useSplatnet, salmonStatsServers, *m
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("json")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("No config file found. One will be created.")
			viper.Set("cookie", "")
			viper.Set("session_token", "")
			viper.Set("user_lang", "")
			viper.Set("user_id", "")
			viper.Set("statink_servers", []types.Server{{
				ShortName: "official",
				ApiKey:    "",
				Address:   "https://stat.ink/api/v2/",
			}})
			viper.Set("salmonstats_servers", []types.Server{{
				ShortName: "official",
				Address:   "https://salmon-stats-api.yuki.games/api/",
			}})
			if err := viper.WriteConfigAs("./config.json"); err != nil {
				panic(err)
			}
		} else {
			// Config file was found but another error was produced
			log.Printf("Error reading the config file. Error is %v", err)
			os.Exit(1)
		}
	}
	viper.SetDefault("cookie", "")
	viper.SetDefault("session_token", "")
	viper.SetDefault("user_lang", "")
	viper.SetDefault("user_id", "")
	viper.SetDefault("statink_servers", []types.Server{{
		ShortName: "official",
		ApiKey:    "",
		Address:   "https://stat.ink/api/v2/",
	}})
	viper.SetDefault("salmonstats_servers", []types.Server{{
		ShortName: "official",
		Address:   "https://salmon-stats-api.yuki.games/api/",
	}})
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	types.CheckForUpdate(client)
	if !(viper.IsSet("user_lang")) || viper.GetString("user_lang") == "" {
		setLanguage()
	}
	stages, hasEvents, tides, weapons, statInkServers, _, salmonStatsServers, _ := getFlags()
	_, timezone := time.Now().Zone()
	timezone = -timezone / 60
	appHead := http.Header{
		"Host":              []string{"app.splatoon2.nintendo.net"},
		"x-unique-id":       []string{"32449507786579989235"},
		"x-requested-with":  []string{"XMLHttpRequest"},
		"x-timezone-offset": []string{fmt.Sprint(timezone)},
		"User-Agent":        []string{"Mozilla/5.0 (Linux; Android 7.1.2; Pixel Build/NJH47D; wv) AppleWebKit/537.36 (KHTML, like Gecko) version/4.0 Chrome/59.0.3071.125 Mobile Safari/537.36"},
		"Accept":            []string{"*/*"},
		"Referer":           []string{"https://app.splatoon2.nintendo.net/home"},
		"Accept-Encoding":   []string{"gzip deflate"},
		"Accept-Language":   []string{viper.GetString("user_lang")},
	}
	iterators := []lib.ShiftIterator{}
	for i := range salmonStatsServers {
		err := salmonstats.GetAllShifts(salmonStatsServers[i], client)
		if err != nil {
			return
		}
		iterators = append(iterators, salmonstats.LoadFromFileIterator(salmonStatsServers[i]))
	}
	for i := range statInkServers {
		statink.GetAllShifts(statInkServers[i], client)
		iterators = append(iterators, statink.LoadFromFileIterator(salmonStatsServers[i]))
	}
	splatnet.GetAllShifts(appHead, client)
	iterators = append(iterators, splatnet.LoadFromFileIterator())
	lib.FindRecords(iterators, stages, hasEvents, tides, weapons, client)
}
