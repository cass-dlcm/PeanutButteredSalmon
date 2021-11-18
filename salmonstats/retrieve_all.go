package salmonstats

import (
	"PeanutButteredSalmon/types"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func GetAllShifts(page int, server types.Server, client *http.Client, save bool) []ShiftSalmonStats {
	if _, err := fmt.Println("Pulling Salmon Run data from online..."); err != nil {
		panic(err)
	}

	data := []ShiftSalmonStats{}

	getShifts := func (page int) ShiftPage {

		url := fmt.Sprintf("%splayers/%s/results", server.Address, viper.GetString("user_id"))

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			panic(err)
		}

		query := req.URL.Query()
		query.Set("raw", "1")
		query.Set("count", "200")
		query.Set("page", fmt.Sprint(page))
		req.URL.RawQuery = query.Encode()

		log.Println(req.URL)



		resp, err := client.Do(req)
		if err != nil {
			log.Panicln(err)
		}

		defer func() {
			if err := resp.Body.Close(); err != nil {
				panic(err)
			}
		}()

		var data ShiftPage

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Panicln(err)
		}

		if save {
			for i := range data.Results {
				fileText, err := json.MarshalIndent(data.Results[i], "", " ")
				if err != nil {
					log.Panicln(err)
				}

				if _, err := os.Stat("salmonstats_shifts"); errors.Is(err, os.ErrNotExist) {
					err := os.Mkdir("salmonstats_shifts", os.ModePerm)
					if err != nil {
						log.Println(err)
					}
				}

				if _, err := os.Stat(fmt.Sprintf("salmonstats_shifts/%s", server.ShortName)); errors.Is(err, os.ErrNotExist) {
					err := os.Mkdir(fmt.Sprintf("salmonstats_shifts/%s", server.ShortName), os.ModePerm)
					if err != nil {
						log.Println(err)
					}
				}

				if err := ioutil.WriteFile(fmt.Sprintf("salmonstats_shifts/%s/%010d.json", server.ShortName, data.Results[i].ID), fileText, 0600); err != nil {
					log.Panicln(err)
				}
			}
		}
		return data
	}
	for true {
		tempData := getShifts(page)
		if len(tempData.Results) == 0 {
			return data
		}
		for i := range tempData.Results {
			tempData.Results[i].Page = page * 200 + i
			tempData.Results[i].PlayerID = viper.GetString("user_id")
			data = append(data, tempData.Results[i])
		}
		page++
	}
	return nil
}


func LoadFromFile(server types.Server) []ShiftSalmonStats {
	f, err := os.Open(fmt.Sprintf("salmonstats_shifts/%s", server.ShortName))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []ShiftSalmonStats{}
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
	data := []ShiftSalmonStats{}
	for i := range files {
		data = append(data, func(fileName string) ShiftSalmonStats{
			f, err := os.Open(fmt.Sprintf("salmonstats_shifts/%s/%s", server.ShortName, fileName))
			if err != nil {
				log.Panicln(err)
			}
			data := ShiftSalmonStats{}
			if err := json.NewDecoder(f).Decode(&data); err != nil {
				if err2 := f.Close(); err2 != nil {
					log.Panicln(fmt.Errorf("%v\n%v", err2, err))
				}
				log.Panicln(err)
			}
			data.PlayerID = viper.GetString("user_id")
			data.Page = i
			return data
		}(files[i]))
	}
	return data
}
