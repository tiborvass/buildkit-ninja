/*
Package config defines types and functions to deal with Ninja configurations.

	c := &Config{
		Vars: map[string]string{
			"cc": "gcc",
			"cflags": "-Wall",
			"obj": "hello.o",
		},
		Rules: map[string]*Rule{
			"compile": {Command: "$cc $cflags -c $in -o $out"},
			"link": {Command: "$cc $in -o $out"},
		},
		BuildEdges: []*BuildEdge{
			{
				RuleName: "compile",
				Inputs: []string{"hello.c"},
				Outputs: []string{"$obj"},
			},
			{
				RuleName: "link",
				Inputs: []string{"hello.o"},
				Outputs: []string{"hello"},
			},
		},
		Default: "hello",
	}

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
				"outputs": ["$obj"],
			},
			{
				"rule": "link",
				"inputs": ["hello.o"],
				"outputs": ["hello"],
			}
		]
	}

*/
package config

type Vars map[string]string
type Rules map[string]*Rule
type BuildEdges []*BuildEdge

type Config struct {
	Vars       Vars
	Rules      Rules
	BuildEdges BuildEdges
	Defaults   []string
}

type Rule struct {
	Command     string
	Description string
}

type BuildEdge struct {
	RuleName string
	Inputs   []string
	Outputs  []string
	Vars     Vars
}
