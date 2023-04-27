package v0

var (
	InputShowType       = map[string]bool{"list": true, "table": true, "tableRow": true, "tableMerge": true, "tableMix": true, "info": true}
	InputFieldInputType = map[string]bool{
		"input":        true,
		"inputs":       true,
		"numberInput":  true,
		"numberInputs": true,
		"switch":       true,
		"select":       true,
		"radio":        true,
		"radioButton":  true,
		"file":         true,
		"files":        true,
		"date":         true,
		"time":         true,
		"dateTime":     true,
		"numberRange":  true,
		"ip":           true,
		"button":       true,
		"numberSlide":  true,
		"textArea":     true,
		"text":         true,
		"hide":         true,
		"password":     true,
	}
	InputReactReactType = map[string]bool{
		"showInput":               true,
		"showInputsByNumber":      true,
		"showInputByUrl":          true,
		"showInputsByNumberByUrl": true,
	}
	InputDependType = map[string]bool{
		"hosts":       true,
		"credentials": true,
	}
)
