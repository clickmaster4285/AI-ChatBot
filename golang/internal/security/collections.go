package security

var AllowedCollections = map[string][]string{
	"clickmastererp": {
		"projects",
		"tasks",
		"users",
	},
}

func IsCollectionAllowed(project string, collection string) bool {

	collections, exists := AllowedCollections[project]
	if !exists {
		return false
	}

	for _, c := range collections {
		if c == collection {
			return true
		}
	}

	return false
}