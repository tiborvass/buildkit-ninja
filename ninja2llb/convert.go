package ninja2llb

import (
	"errors"
	"fmt"
	"os"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/util/system"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/tiborvass/buildkit-ninja/ninja"
	"github.com/tiborvass/buildkit-ninja/ninja/config"
)

func Ninja2LLB(cfg *config.Config) (llb.State, *v1.Image, error) {
	st := llb.Image("gcc")

	b, err := ninja.NewBuilder(cfg)
	if err != nil {
		return st, nil, err
	}
	img := &v1.Image{
		Architecture: "amd64",
		OS:           "linux",
	}
	img.RootFS.Type = "layers"
	img.Config.WorkingDir = "/"
	img.Config.Env = []string{"PATH=" + system.DefaultPathEnv}

	cmd, err := b.CommandFor("hello.o")
	if err != nil {
		return st, nil, err
	}
	fmt.Fprintln(os.Stderr, "totodebug", cmd)

	return st, img, errors.New("TODO: not implemented!")
}
