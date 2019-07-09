package ninja2llb

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
)

type Getter interface {
	Get(key string) (value string, ok bool)
}

type Vars map[string]string

func (vars Vars) Get(key string) (string, bool) {
	v, ok := vars[key]
	return v, ok
}

type Rule struct {
	Command     string
	Description string
}

func (r *Rule) Get(key string) (string, bool) {
	v := reflect.ValueOf(r).Elem().FieldByNameFunc(func(k string) bool {
		return strings.ToLower(k) == key
	})
	if v.IsValid() {
		return v.String(), true
	}
	return "", false
}

type BuildEdge struct {
	Vars    Vars
	Rule    string
	Inputs  []string
	Outputs []string
}

func (be *BuildEdge) String() string {
	return fmt.Sprintf("%v : %s <- %v", be.Outputs, be.Rule, be.Inputs)
}

type Config struct {
	Vars     Vars
	Rules    map[string]Rule
	Builds   []BuildEdge
	Defaults []string
}

func (be *BuildEdge) Get(key string) (string, bool) {
	switch key {
	case "in":
		return strings.Join(be.Inputs, " "), true
	case "in_newline":
		return strings.Join(be.Inputs, "\n"), true
	case "out":
		return strings.Join(be.Outputs, " "), true
	}
	return be.Vars.Get(key)
}

type Scope []Getter

func (s Scope) Get(key string) (string, bool) {
	for _, getter := range s {
		if v, ok := getter.Get(key); ok {
			return v, true
		}
	}
	return "", false
}

func Parse(cfg *Config, r io.Reader) error {
	return json.NewDecoder(r).Decode(cfg)
}
