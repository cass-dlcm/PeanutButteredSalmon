package schedules

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Schedule struct {
	Result []struct {
		Start    string    `json:"start"`
		StartUtc time.Time `json:"start_utc"`
		StartT   int       `json:"start_t"`
		End      string    `json:"end"`
		EndUtc   time.Time `json:"end_utc"`
		EndT     int       `json:"end_t"`
		Stage    struct {
			Image string `json:"image"`
			Name  string `json:"name"`
		} `json:"stage"`
		Weapons []struct {
			ID    int    `json:"id"`
			Image string `json:"image"`
			Name  string `json:"name"`
		} `json:"weapons"`
	} `json:"result"`
}

func GetSchedules(client *http.Client) Schedule {
	url := "https://spla2.yuu26.com/coop"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	data := Schedule{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Panicln(err)
	}
	return data
}