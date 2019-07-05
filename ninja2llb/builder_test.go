package ninja2llb

import (
	"testing"

	"github.com/tiborvass/buildkit-ninja/ninja"
)

func TestBuilder(t *testing.T) {
	cfg := &ninja.Config{
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
				RuleName: "link",
				Inputs:   []string{"hello.o"},
				Outputs:  []string{"hello"},
			},
		},
		Defaults: []string{"hello"},
	}

	b, err := newBuilder(cfg)
	if err != nil {
		t.Fatal(err)
	}

	for _, x := range []struct {
		output   string
		expected string
	}{
		{"hello.o", "gcc -Wall -c hello.c -o hello.o"},
		{"hello", "gcc hello.o -o hello"},
	} {
		t.Run(x.output, func(t *testing.T) {
			cmd, err := b.CommandFor(x.output)
			if err != nil {
				t.Fatal(err)
			}
			if cmd != x.expected {
				t.Fatalf("expected %q got %q", x.expected, cmd)
			}
		})
	}
}
