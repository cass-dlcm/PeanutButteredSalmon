package types

import (
	"fmt"
	"strings"
)

// Tide is a string enum for denoting the water level of the wave.
type Tide string

// The three tides.
const (
	Ht Tide = "HT"
	Nt Tide = "NT"
	Lt Tide = "LT"
)

// GetTideArgs turns a string of space seperated Tide strings into a slice of Tide.
func GetTideArgs(tideStr string) ([]Tide, error) {
	tides := []Tide{}
	tidesStrArr := strings.Split(tideStr, " ")
	for i := range tidesStrArr {
		switch tidesStrArr[i] {
		case "HT", "NT", "LT":
			tides = append(tides, Tide(tidesStrArr[i]))
		default:
			return nil, fmt.Errorf("tide not found: %s", tidesStrArr[i])
		}
	}
	return tides, nil
}

// GetAllTides returns a slice containing every Tide constant.
func GetAllTides() []Tide {
	return []Tide{
		Ht,
		Nt,
		Lt,
	}
}

// TideArr is a wrapper around a Tide slice for the purpose of using the IsAllElementExist function.
type TideArr []Tide

// IsAllElementExist finds whether the given Tide slice contains every element in the TideArr.
func (t *TideArr) IsAllElementExist(arr []Tide) bool {
	for _, i := range *t {
		found := false
		for _, j := range arr {
			if i == j {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}
