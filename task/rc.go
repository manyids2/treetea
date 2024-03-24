package task

type Project struct {
	Name     string
	Children []string
}

type Context struct {
	Read  string
	Write string
}

type Report struct {
	Description string
	Columns     []string
	Labels      []string
	Sort        string
	Filter      string
}

type Urgency struct {
	Coeff float32
}

type UDA struct {
	Type    string
	Label   string
	Values  []string
	Default []string
}

type TaskRC struct {
	Projects map[string]Project
	Contexts map[string]Context
	Reports  map[string]Report
	Urgencys map[string]Urgency
	UDAs     map[string]UDA

	DataLocation   string
	DefaultCommand string
	Weekstart      string
	CaseSensitive  string // yes/no
	Regex          string // yes/no

	ListAllProjects string
	ListAllTags     string
	Verbose         string

	TaskshAutoclear string // true/false
	NewsVersion     string
}
