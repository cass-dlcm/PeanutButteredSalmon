package enums

import (
	"fmt"
	"strings"
)

type Tide string

const (
	Ht Tide = "HT"
	Nt Tide = "NT"
	Lt Tide = "LT"
)

func GetTideArgs(tideStr string) ([]Tide, error) {
	tides := []Tide{}
	eventsStrArr := strings.Split(tideStr, " ")
	for i := range eventsStrArr {
		tideVal, err := stringToTide(eventsStrArr[i])
		if err != nil {
			return nil, err
		}
		tides = append(tides, *tideVal)
	}
	return tides, nil
}

func stringToTide(in string) (*Tide, error) {
	inTide := Tide(in)
	switch in {
	case "HT", "NT", "LT":
		return &inTide, nil
	}
	return nil, fmt.Errorf("tide not found: %s", in)
}

func GetAllTides() []Tide {
	return []Tide{
		Ht,
		Nt,
		Lt,
	}
}
