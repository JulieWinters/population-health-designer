package events

import (
	"fmt"

	"github.com/JulieWinters/population-health-designer/internal/modeling"
)

type Event struct {
	Id       int              `yaml:"id"`
	Time     string           `yaml:"time"`
	Type     string           `yaml:"type"`
	Event    string           `yaml:"event"`
	Patient  modeling.Patient `yaml:"patient"`
	Provider modeling.Person  `yaml:"provider"`
}

func NewEvent(id int, time string, typ string, evnt string) Event {
	event := Event{}
	event.Id = id
	event.Time = time
	event.Type = typ
	event.Event = evnt
	event.Patient = modeling.Patient{}
	event.Patient.Conditions = make([]modeling.Condition, 0)
	event.Provider = modeling.Person{}

	return event
}

// func (event *Event) AddData(key modeling.DataSource, node *yaml.Node) {
// 	if event.Data[key] == nil {
// 		event.Data[key] = make([]*yaml.Node, 1)
// 		event.Data[key][0] = node
// 		return
// 	}
// 	event.Data[key] = append(event.Data[key], node)
// }

func (event *Event) GetMessageKey() string {
	return fmt.Sprintf("%s^%s", event.Type, event.Event)
}

type ByEventTime []Event

func (bet ByEventTime) Len() int {
	return len(bet)
}

func (bet ByEventTime) Less(i int, j int) bool {
	// valuei, err := strconv.ParseInt(bet[i].EventTime[0:4], 0, 64)
	// if err != nil {
	// 	return false, err
	// }
	// valuej, err := strconv.ParseInt(bet[j].EventTime[0:4], 0, 64)
	// if err != nil {
	// 	return false, err
	// }
	// if valuei != valuej {
	// 	return valuei < valuej, nil
	// }

	// valuei, err = strconv.ParseInt(bet[i].EventTime[5:6], 0, 64)
	// if err != nil {
	// 	return false, err
	// }
	// valuej, err = strconv.ParseInt(bet[j].EventTime[5:6], 0, 64)
	// if err != nil {
	// 	return false, err
	// }
	// if valuei != valuej {
	// 	return valuei < valuej, nil
	// }

	return bet[i].Time < bet[j].Time
}

func (bet ByEventTime) Swap(i, j int) {
	bet[i], bet[j] = bet[j], bet[i]
}
