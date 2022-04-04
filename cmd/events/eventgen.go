package events

import (
	"fmt"
	"sort"
	"time"

	"github.com/JulieWinters/population-health-designer/internal/config"
	"github.com/JulieWinters/population-health-designer/internal/modeling"
)

// Executes the event generation process
func Execute(file string) error {
	fmt.Print("Parsing Event Generation Configuration")

	var eventGen modeling.EventGen
	err := config.Parse(file, &eventGen)
	if err != nil {
		return err
	}

	fmt.Printf("Reading patient list from %s\n", eventGen.Stats)
	patients := make([]modeling.Patient, 0)
	err = config.Parse(eventGen.Stats, &patients)
	if err != nil {
		return err
	}

	startMax := time.Now().AddDate(-1*config.MAX_AGE, 0, 0)
	endTime := time.Now()
	events := make([]Event, 0)
	for _, patient := range patients {
		startTime, err := time.Parse("2006-01-02", patient.Demographics.Birthdate)
		if err != nil {
			startTime = startMax
		}

		patEvents, err := makeEvents(&patient, startTime, endTime)
		if err != nil {
			return err
		}
		events = append(events, patEvents...)
	}

	builder := MessageBuilder{}
	err = builder.BuildMessageMap(eventGen.Messages)
	if err != nil {
		return err
	}
	err = builder.BuildSegmentMap(eventGen.Segments)
	if err != nil {
		return err
	}

	sort.Sort(ByEventTime(events))
	for _, event := range events {
		eventGen.ControlId.CurrentId += eventGen.ControlId.Increment
		event.Id = eventGen.ControlId.CurrentId
		message, err := builder.BuildMessage(event)
		if err != nil {
			return err
		}
		config.WriteString(message.Render(), fmt.Sprintf("./%s", message.Name))
	}

	return nil
}

func makeEvents(patient *modeling.Patient, start time.Time, end time.Time) ([]Event, error) {

	reg := Event{}
	reg.Time = start.Format("2006-01-02T15:04:05.000")
	reg.Event = "A01"
	reg.Type = "ADT"
	reg.Patient = patient

	events := make([]Event, 1)
	events[0] = reg

	day := start
	for day.Before(end) {

		day = day.AddDate(0, 0, 1)
	}

	return events, nil
}
