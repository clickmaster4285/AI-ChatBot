package registry

import (
	"encoding/json"
	"os"
)

type ProjectConfig struct {
	ProjectName      string `json:"projectName"`
	DBType           string `json:"dbType"`
	DBName           string `json:"dbName"`
	DBUri            string `json:"dbUri"`
	OnlyGetApiAccess bool   `json:"onlyGetApiAccess"`
	AIModel          string `json:"aiModel"`
}

var Projects = make(map[string]ProjectConfig)

func LoadProjects() error {
	file, err := os.ReadFile("config/projects.json")
	if err != nil {
		return err
	}

	return json.Unmarshal(file, &Projects)
}

func GetProject(name string) (ProjectConfig, bool) {
	project, exists := Projects[name]
	return project, exists
}