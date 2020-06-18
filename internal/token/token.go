package token

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/google/uuid"
)

type TokenMap struct {
	NameToToken map[string]string
	TokenToName map[string]string
}

// Add adds an entry for every name in `names`, it does not change anything if an entry is allready pressent
func (tm *TokenMap) Add(names ...string) {
	for _, name := range names {
		if _, ok := tm.NameToToken[name]; !ok {
			token := generateUUID()
			tm.NameToToken[name] = token
			tm.TokenToName[token] = name
		}
	}
}

// Remove removes the entry for every name in `names`
func (tm *TokenMap) Remove(names ...string) {
	for _, name := range names {
		if token, ok := tm.NameToToken[name]; ok {
			delete(tm.TokenToName, token)
			delete(tm.NameToToken, name)
		}
	}
}

// Contains returns a `[]bool` an bool that is true for every name in the TokenMap, and false for all that is not.
func (tm *TokenMap) Contains(names ...string) []bool {
	out := make([]bool, len(names))
	for index, name := range names {
		_, ok := tm.NameToToken[name]
		out[index] = ok
	}
	return out
}

// GetNames ...
func (tm *TokenMap) GetNames() []string {
	keys := make([]string, len(tm.NameToToken))
	for key := range tm.NameToToken {
		keys = append(keys, key)
	}
	return keys
}

func (tm *TokenMap) SaveToFile(fileDir string) error {
	var err = os.Remove(fileDir)
	/* if err != nil {
		return err
	} */
	file, err := os.OpenFile(fileDir, os.O_CREATE|os.O_RDWR, 0660)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(tm.generateFileString())
	if err != nil {
		return err
	}
	file.Sync()
	return nil
}

func (tm *TokenMap) LoadFromFile(fileDir string) error {
	data, err := ioutil.ReadFile(fileDir)
	if err != nil {
		return err
	}
	tm.parseFileString(string(data))
	return nil
}

// generateUUID generates a UUID using the `github.com/google/uuid` package
func generateUUID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return id.String()
}

func (tm *TokenMap) generateFileString() string {
	var out strings.Builder
	for key, val := range tm.NameToToken {
		out.WriteString(fmt.Sprintf("%s=%s\n", key, val))
	}
	return out.String()
}

func (tm *TokenMap) parseFileString(str string) {
	strs := strings.Split(str, "\n")
	for _, val := range strs {
		i := strings.LastIndex(val, "=")
		if i != -1 {
			name := val[:i]
			token := val[i+1:]
			tm.NameToToken[name] = token
			tm.TokenToName[token] = name
		}
	}
}
