package db

import (
	"aichatbot/internal/registry"
)

type DBConnection struct {
	URI    string
	DBType string
	DBName string
}

// SELECT DATABASE BASED ON PROJECT
func ResolveDB(project registry.ProjectConfig) DBConnection {

	switch project.DBType {

	case "mongodb":
		return DBConnection{
			URI:    project.DBUri,
			DBType: "mongodb",
			DBName: project.DBName,
		}

	case "postgres":
		return DBConnection{
			URI:    project.DBUri,
			DBType: "postgres",
			DBName: project.DBName,
		}

	case "mysql":
		return DBConnection{
			URI:    project.DBUri,
			DBType: "mysql",
			DBName: project.DBName,
		}

	default:
		return DBConnection{
			URI:    "",
			DBType: "unknown",
			DBName: "",
		}
	}
}