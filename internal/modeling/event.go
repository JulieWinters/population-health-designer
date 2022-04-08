package modeling

type DataType string

const (
	PatientSource  DataType = "patient"
	ProviderSource DataType = "provider"
)

type When struct {
	Field  string `yaml:"field"`
	Mask   string `yaml:"mask"`
	Format string `yaml:"format"`
}

type DataSource struct {
	Type   DataType `yaml:"type"`
	Filter string   `yaml:"filter"`
}

type EventDefinition struct {
	Name         string `yaml:"name"`
	MessageType  string `yaml:"message_type"`
	MessageEvent string `yaml:"message_event"`
	// RootSource   DataType `yaml:"root_source"`
	// Phase        Phase     `yaml:"phase"`
	Cardinality string       `yaml:"cardinality"`
	DataSources []DataSource `yaml:"data_sources"`
	Triggers    []string     `yaml:"triggers"`
	When        When         `yaml:"when"`
	// Triggers     []Trigger `yaml:"triggers"`

}

type ControlId struct {
	Seed      int `yaml:"seed"`
	Increment int `yaml:"increment"`
	CurrentId int `yaml:"current_id,omitempty"`
}

type EventGen struct {
	Stats            string            `yaml:"stats"`
	ControlId        ControlId         `yaml:"control_id"`
	Messages         []string          `yaml:"messages"`
	Segments         []string          `yaml:"segments"`
	EventDefinitions []EventDefinition `yaml:"event_definitions"`
}
