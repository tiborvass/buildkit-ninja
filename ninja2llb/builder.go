package ninja2llb

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tiborvass/buildkit-ninja/ninja"
)

type scope interface {
	get(k string) (string, bool)
}

type vars map[string]string

func (vars vars) get(k string) (string, bool) {
	v, ok := vars[k]
	return v, ok
}

type rule struct {
	command     string
	description string
}

func (r *rule) get(k string) (string, bool) {
	v := reflect.ValueOf(r).Elem().FieldByName(k)
	if v.Kind() != reflect.Invalid {
		return v.String(), true
	}
	return "", false
}

type edge struct {
	parent  scope
	rule    *rule
	inputs  []string
	outputs []string
	vars    vars
}

func (e *edge) in() string {
	return strings.Join(e.inputs, " ")
}
func (e *edge) in_newline() string {
	return strings.Join(e.inputs, "\n")
}
func (e *edge) out() string {
	return strings.Join(e.outputs, " ")
}

func (e *edge) get(k string) (string, bool) {
	switch k {
	case "in":
		return e.in(), true
	case "in_newline":
		return e.in_newline(), true
	case "out":
		return e.out(), true
	}

	if v, ok := e.vars.get(k); ok {
		return v, true
	}

	if v, ok := e.rule.get(k); ok {
		return v, true
	}

	return e.parent.get(k)
}

type builder struct {
	vars     vars
	rules    map[string]*rule
	defaults map[*edge]struct{}
	edges    map[string]*edge

	llb map[*edge]llb.State
}

func (b *builder) CommandFor(output string) (string, error) {
	e := b.edges[output]
	if e == nil {
		return "", fmt.Errorf("build edge '%s' not found", output)
	}
	cmd, ok := e.get("command")
	if !ok {
		panic("command field not found")
	}
	return expand(e, cmd), nil
}

func newBuilder(cfg *ninja.Config) (*builder, error) {
	b := &builder{
		vars:     vars(cfg.Vars),
		rules:    make(map[string]*rule, len(cfg.Rules)),
		edges:    make(map[string]*edge, len(cfg.BuildEdges)),
		defaults: make(map[*edge]struct{}, len(cfg.Defaults)),
		llb: make(map[*edge]llb.State, len(cfg.BuildEdges),
	}

	for bei, be := range cfg.BuildEdges {
		rulename := be.RuleName
		cfgRule := cfg.Rules[rulename]
		if cfgRule == nil {
			return nil, fmt.Errorf("build edge #%d references unreachable rule '%s'", bei, be.RuleName)
		}

		r := b.rules[rulename]
		if r == nil {
			r = &rule{
				command:     cfgRule.Command,
				description: cfgRule.Description,
			}
			b.rules[rulename] = r
		}

		e := &edge{
			parent:  b.vars,
			rule:    r,
			vars:    vars(be.Vars),
			inputs:  make([]string, len(be.Inputs)),
			outputs: make([]string, len(be.Outputs)),
		}

		for i, in := range be.Inputs {
			e.inputs[i] = expand(e.parent, in)
		}
		for i, out := range be.Outputs {
			out = expand(e.parent, out)
			e.outputs[i] = out
			b.edges[out] = e
		}
	}

	if len(cfg.Defaults) > 0 {
		b.defaults = make(map[*edge]struct{}, len(cfg.Defaults))
		for _, def := range cfg.Defaults {
			d := expand(b.vars, def)
			e := b.edges[d]
			if e != nil {
				b.defaults[e] = struct{}{}
			} else {
				var extra string
				if d != def {
					extra = fmt.Sprintf(" (evaluated to '%s')", d)
				}
				return nil, fmt.Errorf("could not set default build edge to unreachable '%s'%s", def, extra)
			}
		}
	} else {
		b.defaults = make(map[*edge]struct{}, len(b.edges))
		for _, e := range b.edges {
			b.defaults[e] = struct{}{}
		}
	}

	return b, nil
}
