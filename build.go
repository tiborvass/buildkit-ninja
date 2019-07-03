package bkninja

import (
	"context"
	"fmt"
	"os"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
)

func Build(ctx context.Context, c client.Client) (*client.Result, error) {
	buildFile, err := getBuildFile(c)
	if err != nil {
		return nil, err
	}

	return nil, errors.New("TODO: not implemented!")
}

func getBuildFile(c client.Client) ([]byte, error) {
	opts := c.BuildOpts().Opts
	filename := opts["filename"]
	if filename == "" {
		filename = "build.ninja"
	}

	name := "load build.ninja"
	if filename != "build.ninja" {
		name += " from " + filename
	}

	src := llb.Local("dockerfile",
		llb.IncludePatterns([]string{filename}),
		llb.SessionID(c.BuildOpts().SessionID),
		llb.SharedKeyHint("build.ninja"),
		llb.WithCustomName("[internal] "+name),
	)

	def, err := src.Marshal()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal local source")
	}

	var buildFile []byte
	res, err := c.Solve(ctx, client.SolveRequest{
		Definition: def.ToPB(),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to resolve dockerfile")
	}

	ref, err := res.SingleRef()
	if err != nil {
		return nil, err
	}

	buildFile, err = ref.ReadFile(ctx, client.ReadRequest{
		Filename: filename,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read dockerfile")
	}
}
