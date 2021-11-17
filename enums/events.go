package enums

import (
	"fmt"
	"strings"
)

type Event int

const (
	WaterLevels Event = iota
	Rush
	Fog
	GoldieSeeking
	Griller
	CohockCharge
	Mothership
)

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

func stringToEvent(eventStr string) (*Event, error) {
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

func GetEventArgs(eventStr string) ([]Event, error) {
	events := []Event{}
	eventsStrArr := strings.Split(eventStr, " ")
	for i := range eventsStrArr {
		eventVal, err := stringToEvent(eventsStrArr[i])
		if err != nil {
			return nil, err
		}
		events = append(events, *eventVal)
	}
	return events, nil
}

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
