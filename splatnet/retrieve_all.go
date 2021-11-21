package splatnet

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cass-dlcm/peanutbutteredsalmon/splatnet/iksm"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// GetAllShifts downloads every shiftSplatnet from the SplatNet server and returns them all.
//
// Warnings:
//  • Panics on error.
//
// Breaking changes:
//  • v1->v2:
//    • No return value.
func GetAllShifts(appHead http.Header, client *http.Client, save bool) ShiftList {
	if _, err := fmt.Println("Pulling Salmon Run data from online..."); err != nil {
		panic(err)
	}

	url := "https://app.splatoon2.nintendo.net/api/coop_results"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		panic(err)
	}

	req.Header = appHead

	if viper.GetString("cookie") == "" {
		iksm.GenNewCookie("blank", "1.6.0", client)
	}

	req.AddCookie(&http.Cookie{Name: "iksm_session", Value: viper.GetString("cookie")})

	resp, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	var data ShiftList

	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		iksm.GenNewCookie("auth", "1.6.0", client)
		return GetAllShifts(appHead, client, save)
	}

	if data.Code != nil {
		iksm.GenNewCookie("auth", "1.6.0", client)
		return GetAllShifts(appHead, client, save)
	}

	if save {
		for i := range data.Results {
			fileText, err := json.MarshalIndent(data.Results[i], "", " ")
			if err != nil {
				log.Panicln(err)
			}

			if _, err := os.Stat("shifts"); errors.Is(err, os.ErrNotExist) {
				err := os.Mkdir("shifts", os.ModePerm)
				if err != nil {
					log.Println(err)
				}
			}

			if err := ioutil.WriteFile(fmt.Sprintf("shifts/%d.json", data.Results[i].JobId), fileText, 0600); err != nil {
				log.Panicln(err)
			}
		}
	}

	return data
}

// LoadFromFile reads in files of individual shifts and loads them into a []ShiftSplatnet, returning it.
//
// Warnings:
//  • Panics on error.
func LoadFromFile() []ShiftSplatnet {
	f, err := os.Open("shifts")
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []ShiftSplatnet{}
		}
		log.Panicln(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	files, err := f.Readdirnames(-1)
	if err != nil {
		log.Panicln(err)
	}
	data := []ShiftSplatnet{}
	for i := range files {
		data = append(data, func(fileName string) ShiftSplatnet {
			f, err := os.Open(fmt.Sprintf("shifts/%s", fileName))
			if err != nil {
				log.Panicln(err)
			}
			data := ShiftSplatnet{}
			if err := json.NewDecoder(f).Decode(&data); err != nil {
				if err2 := f.Close(); err2 != nil {
					log.Panicln(fmt.Errorf("%v\n%v", err2, err))
				}
				log.Panicln(err)
			}
			return data
		}(files[i]))
	}
	return data
}
