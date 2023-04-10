// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

/*
Package address contains logic for parsing a Terraform address.

The Terraform address grammar is documented at
https://www.terraform.io/docs/internals/resource-addressing.html

Parsing is implemented using Pigeon, a PEG parser generator.
*/
package address

import (
	"fmt"
	"strings"
)

// Address holds the parsed components of a Terraform address.
type Address struct {
	ModulePath   ModulePath
	ResourceSpec ResourceSpec
}

// NewAddress parses the given address `a` into an Address struct. Returns an
// error if we find a malformed address.
// [module path][resource spec]
func NewAddress(a string) (*Address, error) {
	addr, err := Parse(a, []byte(a))
	if err != nil {
		return nil, err
	}
	return addr.(*Address), nil
}

// Clone copies the memory containing the address structure.
func (a *Address) Clone() *Address {
	mp := make(ModulePath, len(a.ModulePath))
	copy(mp, a.ModulePath)
	return &Address{
		mp,
		a.ResourceSpec,
	}
}

// String representation of the address.
func (a *Address) String() string {
	var prefix string
	if len(a.ModulePath) > 0 {
		prefix = a.ModulePath.String() + "."
	}
	return prefix + a.ResourceSpec.String()
}

// ModulePath holds a list of modules contained in the address. The furthest
// module on the left-hand side (outer-most) of the address is at index 0.
type ModulePath []Module

// String representation of the path component of an address.
func (m ModulePath) String() string {
	modules := make([]string, len(m))
	for i, c := range m {
		modules[i] = c.String()
	}
	return strings.Join(modules, ".")
}

// Index of either a module or a resource. Can either be an int or a string.
type Index struct {
	Value interface{}
}

// String representation of an index. If the index is a string, it will be
// quoted and escaped using go's string escaping semantics.
func (i *Index) String() string {
	if i == nil || i.Value == nil {
		return ""
	}
	switch v := i.Value.(type) {
	case int:
		return fmt.Sprintf("%d", v)
	case string:
		return fmt.Sprintf("%q", v)
	default:
		panic(fmt.Errorf("got unknown type %T", v))
	}
}

// Module represents a module component of an address.
// module.module_name[module index]
type Module struct {
	// Name of the module
	Name string
	// Index of the module. May be nil.
	Index Index
}

// String representation of the module. The literal `module.` will be
// prepended.
func (m *Module) String() string {
	if idx := m.Index.String(); idx != "" {
		return fmt.Sprintf("module.%s[%s]", m.Name, idx)
	}
	return fmt.Sprintf("module.%s", m.Name)
}

// ResourceSpec describes the resource of an address.
// resource_type.resource_name[resource index]
type ResourceSpec struct {
	Type  string
	Name  string
	Index Index
}

// String representation of the resource component of an address.
func (r *ResourceSpec) String() string {
	if idx := r.Index.String(); idx != "" {
		return fmt.Sprintf("%s.%s[%s]", r.Type, r.Name, idx)
	}
	return fmt.Sprintf("%s.%s", r.Type, r.Name)
}
