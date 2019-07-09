package ninja2llb

import (
	"errors"
	"fmt"
	"os"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/util/system"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/tiborvass/buildkit-ninja/ninja"
)

func addEdges(b *builder, dst llb.State, outs []*edge) (llb.State, error) {
	if len(edges) == 0 {
		return llb.State{}, errors.New("empty list of edges")
	}
	dest := dst
	var err error
	for _, e := range out {
		// merge all destinations into one
		dest, err = addEdge(b, dest, e)
		if err != nil {
			return llb.State{}, err
		}
	}
	return dest, nil
}

var gcc = llb.Image("gcc")

func addEdge(b *builder, dst llb.State, out *edge) (llb.State, error) {
	if st, ok := b.llb[out]; ok {
		return st, nil
	}

	cmd, ok := out.get("command")
	if !ok {
		return llb.State{}, fmt.Errorf("command field not found")
	}
	cmd = expand(out, cmd)

	run := gcc.Run(llb.Args([]string{"sh", "-c", cmd}))

	

	st := run.AddMount("/in", dst)

	








	// TODO: if len(out.inputs) == 0 then leaf
	ins := make([]*edge, len(out.inputs))
	for i, input := range out.inputs {
		ins[i], ok = b.edges[input]
		if !ok {
			// if it's not the output of a build, then it's just an input file
			return 
			return llb.State{}, fmt.Errorf("input %q not found", in)
		}
	}
	st, err := addEdges(b, gcc, ins)
	if err != nil {
		return llb.State{}, err
	}
	b.llb[out] = st
	return st, nil
}

func Ninja2LLB(cfg *ninja.Config, dst llb.State) (llb.State, *v1.Image, error) {
	b, err := newBuilder(cfg)
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

	st, err := addEdges(b, dst, b.defaults)

	/*
	cmd, err := b.CommandFor("hello.o")
	if err != nil {
		return st, nil, err
	}
	fmt.Fprintln(os.Stderr, "totodebug", cmd)
	*/

	return st, img, err // errors.New("TODO: not implemented!")
}
