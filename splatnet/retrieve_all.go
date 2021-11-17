package splatnet

import (
	"PeanutButteredSalmon/splatnet/iksm"
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

func getAllShifts(appHead http.Header, client *http.Client, save bool) ShiftList {
	if _, err := fmt.Println("Pulling Salmon Run data from online..."); err != nil {
		panic(err)
	}

	url := "https://app.splatoon2.nintendo.net/api/coop_results"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
		return getAllShifts(appHead, client, save)
	}

	if data.Code != nil {
		iksm.GenNewCookie("auth", "1.6.0", client)
		return getAllShifts(appHead, client, save)
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
