package main

import (
	"PeanutButteredSalmon/lib"
	"PeanutButteredSalmon/types"
	"flag"
	"fmt"
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

func getFlags() ([]types.Stage, []types.Event, []types.Tide, []types.WeaponSchedule, bool, bool, []types.Server, bool, []types.Server, int) {
	stagesStr := flag.String("stage", "spawning_grounds marooners_bay lost_outpost salmonid_smokeyard ruins_of_ark_polaris", "To set a specific set of stages.")
	hasEventsStr := flag.String("event", "water_levels rush fog goldie_seeking griller cohock_charge mothership", "To set a specific set of events.")
	hasTides := flag.String("tide", "LT NT HT", "To set a specific set of tides.")
	hasWeapons := flag.String("weapon", "set single_random four_random random_gold", "To restrict to a specific set of weapon types.")
	save := flag.Bool("save", false, "To save data to json files.")
	load := flag.Bool("load", false, "To load data from json files.")
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
	statInkUrlConf := viper.Get("statink_servers").([]map[string]string)
	statInkServers := []types.Server{}
	for i := range statInkUrlNicks {
		for j := range statInkUrlConf {
			if statInkUrlConf[j]["short_name"] == statInkUrlNicks[i] {
				statInkServers = append(statInkServers, types.Server{ShortName: statInkUrlConf[j]["short_name"], ApiKey: statInkUrlConf[j]["api_key"], Address: statInkUrlConf[j]["address"]})
			}
		}
	}


	salmonStatsUrlNicks := strings.Split(*salmonStats, " ")
	salmonStatsUrlConf := viper.Get("splatstats_servers").([]map[string]string)
	salmonStatsServers := []types.Server{}
	for i := range salmonStatsUrlNicks {
		for j := range salmonStatsUrlConf {
			if salmonStatsUrlConf[j]["short_name"] == salmonStatsUrlNicks[i] {
				salmonStatsServers = append(statInkServers, types.Server{ShortName: salmonStatsUrlConf[j]["short_name"], Address: salmonStatsUrlConf[j]["address"]})
			}
		}
	}


	return stages, hasEvents, tides, weapons, *save, *load, statInkServers, *useSplatnet, salmonStatsServers, *m
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("No config file found. One will be created.")
			viper.Set("cookie", "")
			viper.Set("session_token", "")
			viper.Set("user_lang", "")
			viper.Set("statink_api_key", "")
			if err := viper.WriteConfigAs("./config.yaml"); err != nil {
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
	viper.SetDefault("statink_api_key", "")
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	if !(viper.IsSet("user_lang")) || viper.GetString("user_lang") == "" {
		setLanguage()
	}
	stages, hasEvents, tides, weapons, save, load, statInkServers, useSplatnet, salmonStatsServers, _ := getFlags()
	_, timezone := time.Now().Zone()
	timezone = -timezone / 60
	appHead := http.Header{
		"Host":              []string{"app.splatoon2.nintendo.net"},
		"x-unique-id":       []string{"32449507786579989235"},
		"x-requested-with":  []string{"XMLHttpRequest"},
		"x-timezone-offset": []string{fmt.Sprint(timezone)},
		"User-Agent":        []string{"Mozilla/5.0 (Linux; Android 7.1.2; Pixel Build/NJH47D; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/59.0.3071.125 Mobile Safari/537.36"},
		"Accept":            []string{"*/*"},
		"Referer":           []string{"https://app.splatoon2.nintendo.net/home"},
		"Accept-Encoding":   []string{"gzip deflate"},
		"Accept-Language":   []string{viper.GetString("user_lang")},
	}
	lib.FindRecords(useSplatnet, load, statInkServers, salmonStatsServers, stages, hasEvents, tides, weapons, save, appHead, client)
}
