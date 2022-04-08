package events

import (
	"fmt"
	"strings"
)

//============================================
//===========  Trigger Functions  ============
//============================================
type TriggerFunction string

const (
	TriggerNone TriggerFunction = ""
	Date        TriggerFunction = "date"
)

func toTriggerFunction(input string) (TriggerFunction, error) {
	function := TriggerNone
	if strings.Index(input, string(Date+"(")) == 0 {
		if input[len(input)-1] != ')' {
			return TriggerNone, fmt.Errorf("malformed threshold path '%s'", input)
		}
		function = Date

	}
	return function, nil
}

//============================================
//=============  Data Functions  =============
//============================================
type DataFunction string

const (
	DataNone     DataFunction = ""
	Distribution DataFunction = "distribution"
	AgeToDate    DataFunction = "age_to_date"
)

func toDataFunction(input string) (DataFunction, error) {
	function := DataNone
	if strings.Index(input, string(Distribution+"(")) == 0 {
		if input[len(input)-1] != ')' {
			return DataNone, fmt.Errorf("malformed data path '%s'", input)
		}
		function = Distribution
	} else if strings.Index(input, string(AgeToDate+"(")) == 0 {
		if input[len(input)-1] != ')' {
			return DataNone, fmt.Errorf("malformed data path '%s'", input)
		}
		function = AgeToDate
	}
	return function, nil
}
