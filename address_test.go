// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package address

import (
	"bufio"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const testdata = "./testdata/config/addresses"

func getTestAddresses(t *testing.T) []string {
	f, err := os.Open(testdata)
	require.NoError(t, err)

	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	require.NoError(t, scanner.Err())

	return lines
}

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
		{`module.data.module.b.data.data`},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			a, err := Parse(tt.name, []byte(tt.name)) //, Debug(true))
			require.NoError(t, err)
			require.IsType(t, &Address{}, a)
			require.Equal(t, ManagedResourceMode, a.(*Address).Mode)
			require.Equal(t, tt.name, a.(*Address).String())
		})
	}
}

func TestValidDataAddresses(t *testing.T) {
	var tests = []struct {
		name string
	}{
		{"data.module.module"},
		{"data.foo.bar"},
		{`data.foo.bar["xyz"]`},
		{`module.a.data.foo.bar`},
		{`module.a.data.foo.bar["xyz"]`},
		{`module.a.module.b.data.foo.bar`},
		{`module.a.module.b.data.foo.bar["xyz"]`},
		{`module.a["xyz"].data.foo.bar`},
		{`module.a["xyz"].data.foo.bar["xyz"]`},
		{`module.a["xyz"].module.b.data.foo.bar`},
		{`module.a["xyz"].module.b.data.foo.bar["xyz"]`},
		{`module.a.data.foo.bar[0]`},
		{`module.a.module.b.data.foo.bar`},
		{`module.a.module.b.data.foo.bar[0]`},
		{`module.a[0].data.foo.bar`},
		{`module.a[0].data.foo.bar[0]`},
		{`module.a[0].module.b.data.foo.bar`},
		{`module.a[0].module.b.data.foo.bar[0]`},
		{`module.data.module.b.data.data.data`},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			a, err := Parse(tt.name, []byte(tt.name)) //, Debug(true))
			require.NoError(t, err)
			require.IsType(t, &Address{}, a)
			require.Equal(t, DataResourceMode, a.(*Address).Mode)
			require.Equal(t, tt.name, a.(*Address).String())
		})
	}
}

func TestIndex(t *testing.T) {
	var tests = []struct {
		name     string
		expected string
		given    string
	}{
		{"string", `"foo"`, `module.foo["foo"].a.b["foo"]`},
		{"int", "123", "module.foo[123].a.b[123]"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := NewAddress(tt.given)
			require.NoError(t, err)
			require.Equal(t, a.ResourceSpec.Index.String(), tt.expected)
			require.Equal(t, a.ModulePath[0].Index.String(), tt.expected)
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

func TestEmptyModule(t *testing.T) {
	a, err := NewAddress("foo.bar")
	require.NoError(t, err)
	require.Empty(t, a.ModulePath.String())
}

func TestCopy(t *testing.T) {
	orig := `module.a["xyz"].module.b.foo.bar["xyz"]`
	expected := `module.c.module.b["abc"].foo.baz["xyz"]`
	a, err := NewAddress(orig)
	require.NoError(t, err)
	b := a.Clone()
	b.ModulePath[0] = Module{Name: "c"}
	b.ModulePath[1].Index = Index{"abc"}
	b.ResourceSpec.Name = "baz"

	require.Equal(t, expected, b.String())
	require.Equal(t, orig, a.String())
}

func TestGeneratedAddresses(t *testing.T) {
	addresses := getTestAddresses(t)
	for _, tt := range addresses {
		t.Run(tt, func(t *testing.T) {
			a, err := Parse(tt, []byte(tt)) //, Debug(true))
			require.NoError(t, err)
			require.IsType(t, &Address{}, a)
			require.Equal(t, tt, a.(*Address).String())
		})
	}
}
