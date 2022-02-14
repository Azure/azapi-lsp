package schema

type Modifier string

const (
	Optional Modifier = "Optional"
	Required Modifier = "Required"
)

type Property struct {
	Name        string
	Description string
	Modifier    Modifier
	Type        string
	Value       string
}

type PropertySet struct {
	Name       string
	Properties map[string]Property
}
