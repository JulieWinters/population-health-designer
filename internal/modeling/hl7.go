package modeling

import (
	"strings"
)

type FieldDefinition struct {
	Name        string   `yaml:"name"`
	Cardinality string   `yaml:"cardinality"`
	Values      []string `yaml:"values,omitempty"`
}

type SegmentDefinition struct {
	Id            string            `yaml:"id"`
	Fields        []FieldDefinition `yaml:"fields,omitempty"`
	RepetitionKey string            `yaml:"repetition_key"`
}

type SegmentReference struct {
	Type        string `yaml:"type"`
	Cardinality string `yaml:"cardinality"`
}

type MessageDefinition struct {
	Type     string             `yaml:"type"`
	Event    string             `yaml:"event"`
	Segments []SegmentReference `yaml:"segments"`
}

type Field struct {
	Components []string
}

type Segment struct {
	Name   string
	Fields [][]Field
}

type Hl7v2Message struct {
	Name     string
	Segments []Segment
}

func (msg *Hl7v2Message) Render() string {
	var sb strings.Builder
	for s, seg := range msg.Segments {
		if s > 0 {
			sb.WriteRune('\n')
		}
		sb.WriteString(seg.Name)
		for _, field := range seg.Fields {
			sb.WriteRune('|')
			for f2, rep := range field {
				if f2 > 0 {
					sb.WriteRune('~')
				}
				for c, comp := range rep.Components {
					if c > 0 {
						sb.WriteRune('^')
					}
					sb.WriteString(comp)
				}
			}
		}
	}
	return sb.String()
}

func (msg *Hl7v2Message) PushSegment(segment Segment) {
	if msg.Segments == nil {
		msg.Segments = make([]Segment, 0)
	}
	msg.Segments = append(msg.Segments, segment)
}

func (seg *Segment) NextField() {
	if seg.Fields == nil {
		seg.Fields = make([][]Field, 0)
	}
	seg.Fields = append(seg.Fields, make([]Field, 0))
}

func (seg *Segment) PushField(field Field) {
	if seg.Fields == nil {
		seg.NextField()
	}
	len := len(seg.Fields) - 1
	seg.Fields[len] = append(seg.Fields[len], field)
}

func (field *Field) PushComponent(component string) {
	if field.Components == nil {
		field.Components = make([]string, 0)
	}
	field.Components = append(field.Components, component)
}
