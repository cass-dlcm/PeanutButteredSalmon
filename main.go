package main

import (
	"PeanutButteredSalmon/enums"
	"PeanutButteredSalmon/splatnet"
	"errors"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"time"
)

func setLanguage() {
	log.Println("Please enter your locale (see readme for list).")

	var locale string
	// Taking input from user
	if _, err := fmt.Scanln(&locale); err != nil {
		if !errors.Is(err, errors.New("unexpected newline")) {
			panic(err)
		}
		locale = ""
	}

	if locale == "" {
		viper.Set("user_lang", "en-US")
	} else {
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
	}

	if err := viper.WriteConfig(); err != nil {
		panic(err)
	}
}

func getFlags() ([]enums.Stage, []enums.Event, []enums.Tide, []enums.WeaponSchedule, bool, bool, bool, bool, int) {
	a := flag.Bool("all", false, "To find all personal bests.")
	stagesStr := flag.String("stage", "spawning_grounds marooners_bay lost_outpost salmonid_smokeyard ruins_of_ark_polaris", "To set a specific set of stages.")
	hasEventsStr := flag.String("events", "water_levels rush fog goldie_seeking griller cohock_charge mothership", "To set a specific set of events.")
	hasTides := flag.String("tides", "LT NT HT", "To set a specific set of tides.")
	hasWeapons := flag.String("weapons", "set single_random four_random random_gold", "To restrict to a specific set of weapon types.")
	save := flag.Bool("save", false, "To save data to json files.")
	load := flag.Bool("load", false, "To load data from json files.")
	statInk := flag.Bool("statink", false, "To read data from stat.ink.")
	useSplatnet := flag.Bool("splatnet", false, "To read data from splatnet.")
	m := flag.Int("monitor", -1, "To monitor for new personal bests.")
	flag.Parse()

	if *m != -1 && *useSplatnet {
		log.Panicln("cannot use monitoring mode with stat.ink")
	}

	if *a && (*stagesStr != "spawning_grounds marooners_bay lost_outpost salmonid_smokeyard ruins_of_ark_polaris" || *hasEventsStr != "water_levels rush fog goldie_seeking griller cohock_charge mothership" || *hasWeapons != "set single_random four_random random_gold" || *hasTides != "LT NT HT") {
		log.Panicln("cannot specify filters when using all")
	}

	if *useSplatnet && *statInk {
		log.Panicln("cannot use both stat.ink and splatnet")
	}

	if *save && *load {
		log.Panicln("cannot save and load simultaneously")
	}

	if *load && *useSplatnet {
		log.Panicln("cannot load and use splatnet")
	}

	if *load && *statInk {
		log.Panicln("cannot load and use stat.ink")
	}
	stages, err := enums.GetStageArgs(*stagesStr)
	if err != nil {
		log.Panicln(err)
	}
	hasEvents, err := enums.GetEventArgs(*hasEventsStr)
	if err != nil {
		log.Panicln(err)
	}
	weapons, err := enums.GetWeaponArgs(*hasWeapons)
	if err != nil {
		log.Panicln(err)
	}

	tides, err := enums.GetTideArgs(*hasTides)
	if err != nil {
		panic(err)
	}

	return stages, hasEvents, tides, weapons, *save, *load, *statInk, *useSplatnet, *m
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
	stages, hasEvents, tides, weapons, save, load, useStatInk, useSplatnet, m := getFlags()
	if useSplatnet {
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
		if m != -1 {

		}
		splatnet.FindRecords(stages, hasEvents, tides, weapons, save, appHead, client)
	}
	if useStatInk {

	}
	if load {

	}
}
