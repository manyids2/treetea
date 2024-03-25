package task

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

type Project struct {
	Children map[string]Project
}

type Filters struct {
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
	Contexts map[string]Filters
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

// ReadRc Parse taskrc and return map of keys and values
func ReadRc(path string) (cfg map[string]string, err error) {
	r := regexp.MustCompile(`(?P<Key>[0-9a-zA-Z_\-\.]+)=(?P<Val>.*)`)
	readFile, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	cfg = map[string]string{}
	for fileScanner.Scan() {
		if len(fileScanner.Text()) == 0 {
			continue
		}
		if fileScanner.Text()[0] == '#' {
			continue
		}
		match := r.FindStringSubmatch(fileScanner.Text())
		if len(match) != 3 { // Only of form [key]=[value]
			continue
		}
		cfg[match[1]] = match[2]
	}
	readFile.Close()
	return cfg, nil
}

// ExtractKey Extract 2 level nested keys, and rest are kept at 2 levels
// if one level, the second key is ""
// if no dots, just query the map, do not use this function
func ExtractKey(key string, cfg map[string]string) (values map[string]map[string]string) {
	values = map[string]map[string]string{}
	for k, v := range cfg {
		parts := strings.Split(k, ".")
		if (parts[0] == key) && (len(parts) > 1) {
			_, ok := values[parts[1]]
			if !ok {
				values[parts[1]] = map[string]string{}
			}
			subkey := ""
			if len(parts) > 2 {
				subkey = parts[2]
			}
			if len(parts) > 3 {
				subkey = strings.Join(parts[2:], ".")
			}
			values[parts[1]][subkey] = v
		}
	}
	return values
}

func LoadTaskRC(path string) (rc TaskRC, err error) {
	cfg, err := ReadRc(path)
	if err != nil {
		return rc, err
	}

	// Initialize
	rc = TaskRC{
		DataLocation:    cfg["data.location"],
		DefaultCommand:  cfg["default.command"],
		Weekstart:       cfg["weekstart"],
		CaseSensitive:   cfg["search.case.sensitive"],
		Regex:           cfg["regex"],
		ListAllProjects: cfg["list.all.projects"],
		ListAllTags:     cfg["list.all.tags"],
		Verbose:         cfg["verbose"],
		TaskshAutoclear: cfg["tasksh.autoclear"],
		NewsVersion:     cfg["news.version"],
	}

	// Contexts
	rc.Contexts = map[string]Filters{}
	for k, v := range cfg {
		if len(k) < 9 {
			continue
		}
		if k[:8] == "context." {
			parts := strings.Split(k, ".")
			if _, ok := rc.Contexts[parts[1]]; !ok {
				rc.Contexts[parts[1]] = Filters{}
			}
			if c, ok := rc.Contexts[parts[1]]; ok {
				if parts[2] == "read" {
					c.Read = v
				} else if parts[2] == "write" {
					c.Write = v
				}
				rc.Contexts[parts[1]] = c
			}

		}
	}

	// Reports
	rc.Reports = map[string]Report{}
	for k, v := range cfg {
		if len(k) < 8 {
			continue
		}
		if k[:7] == "report." {
			parts := strings.Split(k, ".")
			if _, ok := rc.Reports[parts[1]]; !ok {
				rc.Reports[parts[1]] = Report{}
			}
			if c, ok := rc.Reports[parts[1]]; ok {
				switch parts[2] {
				case "description":
					c.Description = v
				case "columns":
					c.Columns = strings.Split(v, ",")
				case "labels":
					c.Labels = strings.Split(v, ",")
				case "sort":
					c.Sort = v
				case "filter":
					c.Filter = v
				}
				rc.Reports[parts[1]] = c
			}
		}
	}
	return rc, err
}
