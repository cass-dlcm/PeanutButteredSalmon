package types

import (
	"fmt"
	"strings"
)

// Event is an integer enum for denoting the event of a wave.
type Event int

// The seven Salmon Run events.
const (
	WaterLevels Event = iota
	Rush
	Fog
	GoldieSeeking
	Griller
	CohockCharge
	Mothership
)

// ToString returns the name of the Event, currently hardcoded as the en-US locale.
func (e Event) ToString() string {
	switch e {
	case WaterLevels:
		return "Water Levels"
	case Rush:
		return "Rush"
	case Fog:
		return "Fog"
	case GoldieSeeking:
		return "Goldie Seeking"
	case Griller:
		return "Griller"
	case CohockCharge:
		return "Cohock Charge"
	case Mothership:
		return "Mothership"
	}
	return ""
}

// StringToEvent returns a pointer to an Event if the Event matches the inputted string, otherwise it returns an error.
func StringToEvent(eventStr string) (*Event, error) {
	var eventRes Event
	switch eventStr {
	case "water_levels":
		eventRes = WaterLevels
	case "rush":
		eventRes = Rush
	case "fog":
		eventRes = Fog
	case "goldie_seeking":
		eventRes = GoldieSeeking
	case "griller":
		eventRes = Griller
	case "cohock_charge":
		eventRes = CohockCharge
	case "mothership":
		eventRes = Mothership
	default:
		return nil, fmt.Errorf("event not found: %s", eventStr)
	}
	return &eventRes, nil
}

// GetEventArgs turns a string of space seperated string representations of events into a slice of Event.
func GetEventArgs(eventStr string) ([]Event, error) {
	events := []Event{}
	eventsStrArr := strings.Split(eventStr, " ")
	for i := range eventsStrArr {
		eventVal, err := StringToEvent(eventsStrArr[i])
		if err != nil {
			return nil, err
		}
		events = append(events, *eventVal)
	}
	return events, nil
}

// GetAllEvents returns a slice containing every Event constant.
func GetAllEvents() []Event {
	return []Event{
		WaterLevels,
		Rush,
		Fog,
		GoldieSeeking,
		Griller,
		CohockCharge,
		Mothership,
	}
}

// EventArr is a wrapper around an Event slice for the purpose of using the IsAllElementExist function.
type EventArr []Event

// IsAllElementExist finds whether the given Event slice contains every element in the EventArr.
func (e *EventArr) IsAllElementExist(arr []Event) bool {
	for _, i := range *e {
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
