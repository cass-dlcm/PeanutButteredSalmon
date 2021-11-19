package salmonstats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cass-dlcm/PeanutButteredSalmon/v2/lib"
	"github.com/cass-dlcm/PeanutButteredSalmon/v2/types"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func GetAllShifts(server types.Server, client *http.Client) error {
	if _, err := fmt.Println("Pulling Salmon Run data from online..."); err != nil {
		panic(err)
	}

	getShifts := func(page int) bool {

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
				log.Panicln(err)
			}
		}()

		var data ShiftPage

		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Panicln(err)
		}

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
		return len(data.Results) > 0
	}
	f, err := os.Open(fmt.Sprintf("salmonstats_shifts/%s", server.ShortName))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return err
		}
		log.Panicln(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Panicln(err)
		}
	}(f)
	files, err := f.Readdirnames(-1)
	if err != nil {
		log.Panicln(err)
	}
	page := len(files)/200 + 1
	for getShifts(page) {
		page++
	}
	return nil
}

type ShiftSalmonStatsIterator struct {
	files      []string
	index      int
	serverName string
}

func (s *ShiftSalmonStatsIterator) Next() (lib.Shift, error) {
	if s.index == len(s.files) {
		return nil, errors.New("no more shifts")
	}
	f, err := os.Open(fmt.Sprintf("salmonstats_shifts/%s/%s", s.serverName, s.files[s.index]))
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
	s.index++
	data.Page = s.index
	return data, nil
}

func LoadFromFileIterator(server types.Server) lib.ShiftIterator {
	f, err := os.Open(fmt.Sprintf("salmonstats_shifts/%s", server.ShortName))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		log.Panicln(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)
	returnVal := ShiftSalmonStatsIterator{serverName: server.ShortName}
	returnVal.files, err = f.Readdirnames(-1)
	if err != nil {
		log.Panicln(err)
	}
	return &returnVal
}
