package bkninja

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/containerd/containerd/platforms"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
	"github.com/tiborvass/buildkit-ninja/ninja2llb"
)

const (
	defaultBuildNinjaFilename = "build.ninja"
)

func Build(ctx context.Context, c client.Client) (*client.Result, error) {
	opts := c.BuildOpts().Opts

	_, ignoreCache := opts["no-cache"]

	cfg, err := getNinjaConfig(ctx, c)
	if err != nil {
		return nil, err
	}

	src := llb.Local("context",
		llb.SessionID(c.BuildOpts().SessionID),
		llb.SharedKeyHint("context"),
		llb.WithCustomName("[internal] load build context"),
	)

	st, img, err := ninja2llb.Ninja2LLB(cfg, src, llb.Image("gcc"), ignoreCache)
	if err != nil {
		return nil, err
	}

	def, err := st.Marshal()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal local source")
	}
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

	config, err := json.Marshal(img)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal image config")
	}
	k := platforms.Format(platforms.DefaultSpec())

	res.AddMeta(fmt.Sprintf("%s/%s", exptypes.ExporterImageConfigKey, k), config)
	res.SetRef(ref)

	return res, nil
}

func getNinjaConfig(ctx context.Context, c client.Client) (*ninja2llb.Config, error) {
	opts := c.BuildOpts()
	filename := opts.Opts["filename"]
	if filename == "" {
		filename = defaultBuildNinjaFilename
	}

	name := "load " + defaultBuildNinjaFilename
	if filename != defaultBuildNinjaFilename {
		name += " from " + filename
	}

	src := llb.Local("dockerfile",
		llb.IncludePatterns([]string{filename}),
		llb.SessionID(opts.SessionID),
		llb.SharedKeyHint(defaultBuildNinjaFilename),
		llb.WithCustomName("[internal] "+name),
	)

	def, err := src.Marshal()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal local source")
	}

	// TODO: rulesFile
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

	//return ninja2llb.NewParser(buildFile)
	_ = buildFile
	// TODO: use an actual ninja parser, in the meantime hardcode a config

	// TODO: we should be able to replace hello.o by $obj in the first build rule.
	r := strings.NewReader(`
{
	"vars": {"cc": "gcc", "cflags": "-Wall", "obj": "hello.o"},
	"rules": {
		"compile": {"command": "$cc $cflags -c $in -o $out"},
		"link": {"command": "$cc $in -o $out"}
	},
	"builds": [
		{
			"rule": "compile",
			"inputs": ["hello.c"],
			"outputs": ["hello.o"]
		},
		{
			"rule": "compile",
			"inputs": ["main.c"],
			"outputs": ["main.o"]
		},
		{
			"rule": "link",
			"inputs": ["hello.o", "main.o"],
			"outputs": ["hello"]
		}
	],
	"defaults": ["hello"]
}`)

	cfg := &ninja2llb.Config{}
	if err := ninja2llb.Parse(cfg, r); err != nil {
		return nil, err
	}
	return cfg, nil

	/*
		&ninja.Config{
			Vars: ninja.Vars{
				"cc":     "gcc",
				"cflags": "-Wall",
				"obj":    "hello.o",
			},
			Rules: ninja.Rules{
				"compile": {Command: "$cc $cflags -c $in -o $out"},
				"link":    {Command: "$cc $in -o $out"},
			},
			BuildEdges: ninja.BuildEdges{
				{
					RuleName: "compile",
					Inputs:   []string{"hello.c"},
					Outputs:  []string{"$obj"},
				},
				{
					RuleName: "compile",
					Inputs:   []string{"main.c"},
					Outputs:  []string{"main.o"},
				},
				{
					RuleName: "link",
					Inputs:   []string{"hello.o", "main.o"},
					Outputs:  []string{"hello"},
				},
			},
			Defaults: []string{"hello"},
		}, nil
	*/
}
