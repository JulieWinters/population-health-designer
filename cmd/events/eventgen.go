package events

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	// "strings"
	"time"

	"github.com/JulieWinters/population-health-designer/internal/config"
	"github.com/JulieWinters/population-health-designer/internal/modeling"
	"github.com/JulieWinters/population-health-designer/internal/utils"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
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

		patEvents, err := makeEvents(&patient, &eventGen, startTime, endTime)
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

func makeEvents(patient *modeling.Patient, eventGen *modeling.EventGen, start time.Time, end time.Time) ([]Event, error) {

	patientNode, err := config.StructToNode(patient)
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int)
	events := make([]Event, 0)
	for _, definition := range eventGen.EventDefinitions {

		if count, ok := counts[definition.Name]; ok {
			_, high, err := utils.ParseCardinality(definition.Cardinality)
			if err != nil {
				return nil, err
			}
			if count > high {
				continue
			}
		}

		nodes := make([]*yaml.Node, 0)
		for _, trigger := range definition.Triggers {
			tPath, err := yamlpath.NewPath(trigger)
			if err != nil {
				return nil, err
			}
			tNodes, err := tPath.Find(&patientNode)
			if err != nil {
				return nil, err
			}
			nodes = append(nodes, tNodes...)
		}

		for _, node := range nodes {
			timePath, err := yamlpath.NewPath(definition.When.Field)
			if err != nil {
				return nil, err
			}
			whenNode, err := timePath.Find(node)
			if err != nil {
				return nil, err
			}

			if len(whenNode) != 1 {
				return nil, fmt.Errorf("when field of '%v' did not return exactly 1 data node", definition.When.Field)
			}

			t := config.RandMaskedTime(whenNode[0].Value, definition.When.Mask)
			event := NewEvent(0, t, definition.MessageType, definition.MessageEvent)

			for _, source := range definition.DataSources {
				// var sourceRootNode *yaml.Node
				// var sourceDestinationNode yaml.Node
				if source.Type == modeling.PatientSource {
					event.Patient = *patient
					// sourceRootNode = &patientNode
					// sourceDestinationNode, err = config.StructToNode(event.Patient)
					// if err != nil {
					// 	return nil, err
					// }
				} else if source.Type == modeling.ProviderSource {
					return nil, fmt.Errorf("passing provider data is unsupported at this time")
				} else {
					return nil, fmt.Errorf("passing %v data is unsupported at this time", source.Type)
				}

				// sourcePath, err := yamlpath.NewPath(source.Filter)
				// if err != nil {
				// 	return nil, err
				// }

				// sourceNodes, err := sourcePath.Find(sourceRootNode)
				// if err != nil {
				// 	return nil, err
				// }

				// for _, sourceNode := range sourceNodes {
				// 	err = assignSourceData(source.Filter, &sourceDestinationNode, sourceNode)
				// 	if err != nil {
				// 		return nil, err
				// 	}
				// }

				// if source.Type == modeling.PatientSource {
				// 	sourceDestinationNode.Decode(&event.Patient)
				// } else if source.Type == modeling.ProviderSource {
				// 	return nil, fmt.Errorf("passing provider data is unsupported at this time")
				// } else {
				// 	return nil, fmt.Errorf("passing %v data is unsupported at this time", source.Type)
				// }
			}
			events = append(events, event)
		}
	}

	return events, nil
}

func assignSourceData(sourcePath string, destination *yaml.Node, data *yaml.Node) error {
	destinationNodePath := sourcePath
	open := strings.IndexRune(destinationNodePath, '[')
	for ; open != -1; open = strings.IndexRune(destinationNodePath, '[') {
		close := strings.IndexRune(destinationNodePath, ']')
		temp := destinationNodePath[:open]
		if close+1 < len(destinationNodePath) {
			temp += destinationNodePath[close+1:]
		}
		destinationNodePath = temp
	}
	parts := strings.Split(destinationNodePath, ".")

	content := &destination.Content
	for _, part := range parts {
		foundPart := false
		for i := 0; i < len(*content); i += 2 {
			k := (*content)[i].Kind
			if k == yaml.ScalarNode && (*content)[i].Value == part {
				foundPart = true
				content = &(*content)[i+1].Content
				break
			}
			if k == yaml.MappingNode {
				content = &(*content)[i].Content
				i = -2
			}
		}
		if !foundPart {
			return errors.New("unable to assign data source")
		}
	}
	_ = append((*content), data)
	return nil

}
