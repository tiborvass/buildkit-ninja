package ninja

import "testing"

func TestExpand(t *testing.T) {
	variables := vars{"foo": "bar", "foofoo": "baz"}
	for _, x := range []struct {
		name     string
		vars     vars
		in       string
		expected string
	}{
		{
			name:     "empty",
			vars:     nil,
			in:       "",
			expected: "",
		},
		{
			name:     "no substitution",
			vars:     variables,
			in:       "hi foo bar",
			expected: "hi foo bar",
		},
		{
			name:     "simple substitution",
			vars:     variables,
			in:       "hi $foo !",
			expected: "hi bar !",
		},
	} {
		t.Run(x.name, func(t *testing.T) {
			out := expand(x.vars, x.in)
			if out != x.expected {
				t.Fatalf("expand(%v, %q) = %q    expected %q", x.vars, x.in, out, x.expected)
			}
		})
	}
}
