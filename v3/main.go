package main

import (
	"flag"
	"fmt"
	"github.com/cass-dlcm/PeanutButteredSalmon/v3/lib"
	"github.com/cass-dlcm/PeanutButteredSalmon/v3/salmonstats"
	"github.com/cass-dlcm/PeanutButteredSalmon/v3/splatnet"
	"github.com/cass-dlcm/PeanutButteredSalmon/v3/statink"
	"github.com/cass-dlcm/PeanutButteredSalmon/v3/types"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strings"
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
			log.Panicf("Error reading the config file. Error is %v\n", err)
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
	if errs := types.CheckForUpdate(client); len(errs) > 0 {
		for i := range errs {
			log.Println(errs[i])
		}
		log.Panicln(nil)
	}
	if !(viper.IsSet("user_lang")) || viper.GetString("user_lang") == "" {
		setLanguage()
	}
	stages, hasEvents, tides, weapons, statInkServers, useSplatnet, salmonStatsServers, _ := getFlags()
	iterators := []lib.ShiftIterator{}
	if useSplatnet {
		errs := splatnet.GetAllShifts(client)
		if errs != nil {
			log.Panicln(errs)
		}
		iter, err := splatnet.LoadFromFileIterator()
		if err != nil {
			log.Panicln(err)
		}
		iterators = append(iterators, iter)
	}
	for i := range salmonStatsServers {
		errs := salmonstats.GetAllShifts(salmonStatsServers[i], client)
		if errs != nil {
			log.Panicln(errs)
		}
		iter, err := salmonstats.LoadFromFileIterator(salmonStatsServers[i])
		if err != nil {
			log.Panicln(err)
		}
		iterators = append(iterators, iter)
	}
	for i := range statInkServers {
		errs := statink.GetAllShifts(statInkServers[i], client)
		if errs != nil {
			log.Panicln(errs)
		}
		iter, err := statink.LoadFromFileIterator(salmonStatsServers[i])
		if err != nil {
			log.Panicln(err)
		}
		iterators = append(iterators, iter)
	}
	log.Println(iterators)
	lib.FindRecords(iterators, stages, hasEvents, tides, weapons, client)
}
