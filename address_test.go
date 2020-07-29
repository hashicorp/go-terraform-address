package address

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidAddresses(t *testing.T) {
	var tests = []struct {
		name string
	}{
		{"module.module"},
		{"foo.bar"},
		{`foo.bar["xyz"]`},
		{`module.a.foo.bar`},
		{`module.a.foo.bar["xyz"]`},
		{`module.a.module.b.foo.bar`},
		{`module.a.module.b.foo.bar["xyz"]`},
		{`module.a["xyz"].foo.bar`},
		{`module.a["xyz"].foo.bar["xyz"]`},
		{`module.a["xyz"].module.b.foo.bar`},
		{`module.a["xyz"].module.b.foo.bar["xyz"]`},
		{`module.a.foo.bar[0]`},
		{`module.a.module.b.foo.bar`},
		{`module.a.module.b.foo.bar[0]`},
		{`module.a[0].foo.bar`},
		{`module.a[0].foo.bar[0]`},
		{`module.a[0].module.b.foo.bar`},
		{`module.a[0].module.b.foo.bar[0]`},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			a, err := Parse(tt.name, []byte(tt.name)) //, Debug(true))
			require.NoError(t, err)
			require.IsType(t, &Address{}, a)
			require.Equal(t, tt.name, a.(*Address).String())
		})
	}
}

func TestIndexEdgecases(t *testing.T) {
	var tests = []string{
		`foo"bar"`,
		`123`,
		"a`b",
		"a'b",
		"!@*(ÔASd//\\",
	}
	tpl := "module.foo[%q].a.b[%q]"
	for _, tt := range tests {
		addr := fmt.Sprintf(tpl, tt, tt)
		t.Run(addr, func(t *testing.T) {
			a, err := Parse(addr, []byte(addr)) //, Debug(true))
			require.NoError(t, err)
			require.IsType(t, &Address{}, a)
			require.Equal(t, addr, a.(*Address).String())
		})
	}
}

func TestInvalidAddresses(t *testing.T) {
	var tests = []struct {
		name string
	}{
		{"foo"},
		{`foo["xyz"]`},
		{`foo["xyz"`},
		{`foo["xyz]`},
		{`module.foo.bar`},
		{`module.a.foo.bar["x"yz"]`},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			a, err := Parse(tt.name, []byte(tt.name)) //, Debug(true))
			require.Error(t, err)
			require.Zero(t, a)
		})
	}
}

func TestNewAddress(t *testing.T) {
	a, err := NewAddress("foo.bar")
	require.NoError(t, err)
	require.Equal(t, "foo.bar", a.String())
}
