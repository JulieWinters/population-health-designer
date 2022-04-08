package events

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/JulieWinters/population-health-designer/internal/config"
	"github.com/JulieWinters/population-health-designer/internal/modeling"
	"github.com/JulieWinters/population-health-designer/internal/utils"
	"github.com/vmware-labs/yaml-jsonpath/pkg/yamlpath"
	"gopkg.in/yaml.v3"
	// "gopkg.in/yaml.v3"
)

type MessageBuilder struct {
	Messages map[string]modeling.MessageDefinition
	Segments map[string]modeling.SegmentDefinition
}

func (builder *MessageBuilder) BuildMessageMap(files []string) error {
	if builder.Messages == nil {
		builder.Messages = make(map[string]modeling.MessageDefinition)
	}
	for _, file := range files {
		definition := modeling.MessageDefinition{}
		err := config.Parse(file, &definition)
		if err != nil {
			return err
		}
		key := fmt.Sprintf("%s^%s", definition.Type, definition.Event)
		builder.Messages[key] = definition
	}
	return nil
}

func (builder *MessageBuilder) BuildSegmentMap(files []string) error {
	if builder.Segments == nil {
		builder.Segments = make(map[string]modeling.SegmentDefinition)
	}
	for _, file := range files {

		definition := modeling.SegmentDefinition{}

		err := config.Parse(file, &definition)
		if err != nil {
			return err
		}
		builder.Segments[definition.Id] = definition
	}
	return nil
}

func (builder *MessageBuilder) BuildMessage(event Event) (modeling.Hl7v2Message, error) {

	key := event.GetMessageKey()
	mDef := builder.Messages[key]
	message := modeling.Hl7v2Message{}
	message.Name = fmt.Sprintf("%s-%s-%s.hl7", event.Time, event.Type, event.Event)

	eventNode, err := config.StructToNode(event)
	if err != nil {
		return message, err
	}

	for _, seg := range mDef.Segments {
		if seg.Cardinality == "0..0" {
			//This segment was set to zero cardinality for some reason
			continue
		}

		min, max, err := utils.ParseCardinality(seg.Cardinality)
		if err != nil {
			return message, err
		}

		sDef, ok := builder.Segments[seg.Type]
		if !ok {
			if min == 0 {
				continue
			}
			return message, fmt.Errorf("no definition for required segment %s for message %s^%s", sDef.Id, mDef.Type, mDef.Event)
		}

		var inputNodes []*yaml.Node
		// inputNodes := event.Data[sDef.DataKey]
		if sDef.RepetitionKey != "" {
			path, err := yamlpath.NewPath(sDef.RepetitionKey)
			if err != nil {
				return message, err
			}
			nodes, err := path.Find(&eventNode)
			if err != nil {
				return message, err
			}
			count := len(nodes)
			if count < max {
				max = count
			}

			inputNodes = make([]*yaml.Node, 0)
			inputNodes = append(inputNodes, nodes...)

		} else {
			inputNodes = make([]*yaml.Node, 1)
			inputNodes[0] = &eventNode
		}

		for segN := 1; segN <= max; segN++ {
			segment := modeling.Segment{}
			segment.Name = seg.Type

			inputIndex := 0
			if segN <= len(inputNodes) {
				inputIndex = segN - 1
			}
			inputNode := inputNodes[inputIndex]
			for _, fDef := range sDef.Fields {
				segment.NextField()
				if fDef.Cardinality == "0..0" {
					// This is a zero cardinality field
					continue
				}

				field := modeling.Field{}
				for _, vDef := range fDef.Values {
					values := make([]string, 0)
					open := strings.IndexRune(vDef, '{')
					close := strings.IndexRune(vDef, '}')
					if open == -1 || close == -1 || close < open {
						// use the template as-is since there is nothing to substitute
						values = append(values, vDef)
					} else if vDef[open+1] == '{' || (close+1 < len(vDef) && vDef[close+1] == '{') {
						// use the template as-is since the braces are escaped
						values = append(values, vDef)
					} else {
						pathStr, function, err := parseDataPath(vDef[open+1:close], segN)
						if err != nil {
							return message, err
						}

						if strings.Index(pathStr, "meta.sequence") == 0 {
							values = append(values, strconv.Itoa(segN))
						} else if strings.Index(pathStr, "meta.") == 0 {
							value, err := getMetaValue(pathStr, event)
							if err != nil {
								return message, err
							}
							values = append(values, value)
						} else {

							path, err := yamlpath.NewPath(pathStr)
							if err != nil {
								return message, err
							}

							nodes, err := path.Find(inputNode)
							if err != nil {
								return message, err
							}
							for _, node := range nodes {
								value, err := function.processData(node.Value, event)
								if err != nil {
									return message, err
								}
								values = append(values, value)
							}
						}
					}

					for _, value := range values {
						field.PushComponent(value)
					}
				}
				segment.PushField(field)
			}
			message.PushSegment(segment)
		}
	}
	return message, nil
}

func parseDataPath(path string, n int) (string, DataFunction, error) {
	pathStr := strings.ToLower(path)
	pathStr = strings.ReplaceAll(pathStr, "[n]", fmt.Sprintf("[%d]", n))

	function, err := toDataFunction(pathStr)
	if err != nil {
		return "", DataNone, err
	}
	return pathStr, function, nil
}

func (function DataFunction) processData(data string, event Event) (string, error) {
	if function == Distribution {
		return "", errors.New("'distribution' data function is not implemented")
	}
	return data, nil
}

func getMetaValue(pathStr string, event Event) (string, error) {
	switch pathStr {
	case "meta.date_time":
		hl7Time := strings.ReplaceAll(event.Time, "-", "")
		hl7Time = strings.ReplaceAll(hl7Time, ":", "")
		hl7Time = strings.ReplaceAll(hl7Time, "T", "")
		index := strings.IndexRune(hl7Time, '.')
		if index != -1 {
			hl7Time = hl7Time[:index]
		}
		return hl7Time, nil
	case "meta.msg_type":
		return event.Type, nil
	case "meta.msg_trigger":
		return event.Event, nil
	case "meta.control_id":
		return strconv.Itoa(event.Id), nil
	case "meta.reason":
		return "REASON TODO", nil
	default:
		return "", fmt.Errorf("unsupported meta substitution of {%s}", pathStr)
	}
}
