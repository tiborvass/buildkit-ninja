package ninja2llb

import "testing"

func TestExpand(t *testing.T) {
	variables := vars{"foo": "bar", "foofoo": "baz", "has spaces": "qux", "has": "incorrect substitution"}
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
			in:       "hi $foo!",
			expected: "hi bar!",
		},
		{
			name: "curly substitute",
			vars: variables,
			in: "hi ${foo}!",
			expected: "hi bar!",
		},
		{
			name: "has spaces",
			vars: variables,
			in: "hi ${has spaces}",
			expected: "hi qux",
		},
		{
			name: "precedence",
			vars: variables,
			in: "hi foo$foofoo$foo",
			expected: "hi foobazbar",
		},
		{
			name: "escape",
			vars: nil,
			in: "$$$ $:$\n",
			expected: "$ :\n",
		},
	} {
		t.Run(x.name, func(t *testing.T) {
			out := expand(x.vars, x.in)
			if out != x.expected {
				t.Fatalf("\nexpand(\n\t%#v,\n\t%q\n) = %q (expected %q)", x.vars, x.in, out, x.expected)
			}
		})
	}
}
