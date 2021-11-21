package statink

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cass-dlcm/peanutbutteredsalmon/types"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// GetAllShifts downloads every shiftStatInk from the provided stat.ink server.
// It returns an array of the downloaded shifts.
// It has the option of saving the shifts to individual JSON files.
//
// Warning:
//  • Panics on error.
//
// Breaking changes v1->v2:
//  • Automatically saves all the files. (no `save bool` param).
//  • Downloads from where it stopped last time by looking at the files. (no `id int` param).
//  • Returns an error instead of []ShiftStatInk.
func GetAllShifts(id int, statInkServer types.Server, client *http.Client, save bool) []ShiftStatInk {
	var data []ShiftStatInk
	getShift := func(id int) []ShiftStatInk {
		url := fmt.Sprintf("%suser-salmon", statInkServer.Address)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			panic(err)
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", statInkServer.ApiKey))
		query := req.URL.Query()
		query.Set("newer_than", fmt.Sprint(id))
		query.Set("order", "asc")
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
		var data []ShiftStatInk
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			log.Panicln(err)
		}
		if save {
			for i := range data {
				fileText, err := json.MarshalIndent(data[i], "", " ")
				if err != nil {
					log.Panicln(err)
				}

				if _, err := os.Stat("statink_shifts"); errors.Is(err, os.ErrNotExist) {
					err := os.Mkdir("statink_shifts", os.ModePerm)
					if err != nil {
						log.Println(err)
					}
				}

				if _, err := os.Stat(fmt.Sprintf("statink_shifts/%s/", statInkServer.ShortName)); errors.Is(err, os.ErrNotExist) {
					err := os.Mkdir(fmt.Sprintf("statink_shifts/%s", statInkServer.ShortName), os.ModePerm)
					if err != nil {
						log.Println(err)
					}
				}

				if err := ioutil.WriteFile(fmt.Sprintf("statink_shifts/%s/%010d.json", statInkServer.ShortName, data[i].SplatnetNumber), fileText, 0600); err != nil {
					log.Panicln(err)
				}
			}
		}
		return data
	}
	for true {
		tempData := getShift(id)
		if len(tempData) == 0 || (len(data) != 0 && tempData[len(tempData)-1].SplatnetNumber == data[len(data)-1].ID) {
			return data
		}
		data = append(data, tempData...)
		id = data[len(data)-1].ID
	}
	return nil
}

// LoadFromFile reads all the files in at once and loads them into a slice of shifts.
//
// Warning:
//  • Panics on error.
func LoadFromFile(statInkServer types.Server) []ShiftStatInk {
	f, err := os.Open(fmt.Sprintf("statink_shifts/%s", statInkServer.ShortName))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return []ShiftStatInk{}
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
	data := []ShiftStatInk{}
	for i := range files {
		data = append(data, func(fileName string) ShiftStatInk {
			f, err := os.Open(fmt.Sprintf("statink_shifts/%s/%s", statInkServer.ShortName, fileName))
			if err != nil {
				log.Panicln(err)
			}
			data := ShiftStatInk{}
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
