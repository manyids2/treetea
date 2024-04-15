package task

import (
	"fmt"
	"sort"
	"testing"
)

func TestTaskProjects(t *testing.T) {
	_, err := TaskProjects()
	if err != nil {
		t.Fatalf("Error: %s\n", err)
	}
}

func TestSaveProjects(t *testing.T) {
	projects, err := TaskProjects() // TODO: Use synth
	if err != nil {
		t.Fatalf("Error: %s\n", err)
	}
	taskdata, _ := GetEnv()
	err = SaveProjects(projects, taskdata+"/projects-test.data")
	if err != nil {
		t.Fatalf("Error: %s\n", err)
	}
	projects, err = ReadProjects(taskdata + "/projects-test.data")
	if err != nil {
		t.Fatalf("Error: %s\n", err)
	}
}

func equalProjects(want, projects []Project, do_sort bool) bool {
	if do_sort {
		// Sort alphabetically
		sort.Slice(projects, func(i, j int) bool {
			return projects[i].Name < projects[j].Name
		})
		sort.Slice(want, func(i, j int) bool {
			return want[i].Name < want[j].Name
		})
	}
	for i := range projects {
		if (want[i].Name != projects[i].Name) ||
			(len(want[i].Subprojects) != len(projects[i].Subprojects)) {
			return false
		}
		for j := range want[i].Subprojects {
			if want[i].Subprojects[j] != projects[i].Subprojects[j] {
				return false
			}
		}
	}
	return true
}

func TestMergeProjects(t *testing.T) {
	projects_a := []Project{
		{Name: "a", Subprojects: []string{}},
		{Name: "b", Subprojects: []string{"b.1"}},
		{Name: "b.1", Subprojects: []string{}},
	}
	projects_b := []Project{
		{Name: "a", Subprojects: []string{}},
		{Name: "d", Subprojects: []string{}},
	}
	want := []Project{
		{Name: "a", Subprojects: []string{}},
		{Name: "b", Subprojects: []string{"b.1"}},
		{Name: "b.1", Subprojects: []string{}},
		{Name: "d", Subprojects: []string{}},
	}
	projects := MergeProjects(projects_a, projects_b)
	if !equalProjects(want, projects, true) {
		t.Fatalf("Error: \nwant: %#v\ngot : %#v\n", want, projects)
	}
	fmt.Println(projects)
}

func TestImportProjects(t *testing.T) {
	projects_from := []Project{
		{Name: "a", Subprojects: []string{}},
		{Name: "d", Subprojects: []string{}},
	}
	projects_to := []Project{
		{Name: "a", Subprojects: []string{}},
		{Name: "b", Subprojects: []string{"b.1"}},
		{Name: "b.1", Subprojects: []string{}},
	}
	want := []Project{
		{Name: "a", Subprojects: []string{"a.a", "a.d"}},
		{Name: "a.a", Subprojects: []string{}},
		{Name: "a.d", Subprojects: []string{}},
		{Name: "b", Subprojects: []string{"b.1"}},
		{Name: "b.1", Subprojects: []string{}},
	}
	projects := ImportProjects(projects_from, projects_to, "a")
	if !equalProjects(want, projects, true) {
		t.Fatalf("Error: \nwant: %#v\ngot : %#v\n", want, projects)
	}
	fmt.Println(projects)
}
