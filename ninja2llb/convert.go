package ninja2llb

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/util/system"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

const srcPrefix = "/src"

type defaultOutput struct {
	edgeIdx   int
	outputIdx int
}

type converter struct {
	*Config
	outs     map[string]llb.State
	defaults []llb.State

	source  llb.State
	builder llb.State
}

func (c *converter) addEdge(be *BuildEdge) error {
	rule, ok := c.Rules[be.Rule]
	if !ok {
		return fmt.Errorf("rule %q referenced by %s not found", be.Rule, be)
	}
	scope := Scope{be, &rule, c.Vars}
	cmd, ok := scope.Get("command")
	if !ok {
		return errors.New("'command' field not found")
	}
	cmd = Expand(cmd, scope)

	r := c.builder.Run(llb.Args([]string{"sh", "-c", cmd}), llb.Dir(srcPrefix))
	for _, in := range be.Inputs {
		st, ok := c.outs[in]
		prefixedIn := filepath.Join(srcPrefix, in)
		if ok {
			_ = r.AddMount(prefixedIn, st, llb.SourcePath(prefixedIn), llb.Readonly)
		} else {
			_ = r.AddMount(prefixedIn, c.source, llb.SourcePath(in), llb.Readonly)
		}
	}
	for _, out := range be.Outputs {
		c.outs[out] = r.AddMount(srcPrefix, llb.Scratch())
	}
	return nil
}

func (c *converter) Convert() (llb.State, error) {
defaults:
	for _, def := range c.Defaults {
		// more likely to be closer to the last build edges
		for i := len(c.Builds) - 1; i >= 0; i-- {
			e := c.Builds[i]
			for _, out := range e.Outputs {
				if out == def {
					if err := c.addEdge(&e); err != nil {
						return llb.State{}, err
					}
					st, ok := c.outs[out]
					if !ok {
						panic(fmt.Errorf("expected to find output %q", out))
					}
					c.defaults = append(c.defaults, st)
					continue defaults
				}
			}
		}
	}

	if len(c.defaults) == 1 {
		return c.defaults[0], nil
	}

	return llb.State{}, errors.New("TODO: multiple defaults not implemented")
	/*
		for _, def := range c.defaults {
			_ = def
		}
	*/
}

func Ninja2LLB(cfg *Config, src, builder llb.State) (llb.State, *v1.Image, error) {

	c := &converter{
		Config:  cfg,
		outs:    make(map[string]llb.State),
		source:  src,
		builder: builder,
	}

	st, err := c.Convert()
	if err != nil {
		return llb.State{}, nil, err
	}

	img := &v1.Image{
		Architecture: "amd64",
		OS:           "linux",
	}
	img.RootFS.Type = "layers"
	img.Config.WorkingDir = "/"
	img.Config.Env = []string{"PATH=" + system.DefaultPathEnv}

	return st, img, nil
}
