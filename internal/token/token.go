package token

import (
	"github.com/google/uuid"
)

// GenerateUUID generates a UUID using the `github.com/google/uuid` package
func GenerateUUID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return id.String()
}

// CreateTokenMaps ...
func CreateTokenMaps(names []string) (map[string]string, map[string]string) {
	name2token := make(map[string]string)
	token2name := make(map[string]string)
	for i := 0; i < len(names); i++ {
		token := GenerateUUID()
		name2token[names[i]] = token
		token2name[token] = names[i]
	}
	return name2token, token2name
}

// TODO caching
