package address

import (
	"fmt"
	"strings"
)

type Address struct {
	ModulePath   ModulePath
	ResourceSpec ResourceSpec
}

// NewAddress parses the given address `a` into an Address struct
func NewAddress(a string) (*Address, error) {
	addr, err := Parse(a, []byte(a))
	if err != nil {
		return nil, err
	}
	return addr.(*Address), nil
}

func (a *Address) String() string {
	var prefix string
	if len(a.ModulePath) > 0 {
		prefix = a.ModulePath.String() + "."
	}
	return prefix + a.ResourceSpec.String()
}

type ModulePath []Module

func (m ModulePath) String() string {
	modules := make([]string, len(m))
	for i, c := range m {
		modules[i] = c.String()
	}
	return strings.Join(modules, ".")
}

type Index struct {
	Value interface{}
}

func (i *Index) String() string {
	if i == nil || i.Value == nil {
		return ""
	}
	switch v := i.Value.(type) {
	case int:
		return fmt.Sprintf("[%d]", v)
	case string:
		return fmt.Sprintf("[%q]", v)
	default:
		panic(fmt.Errorf("got unknown type %T", v))
	}
}

type Module struct {
	Name  string
	Index Index
}

func (m *Module) String() string {
	return fmt.Sprintf("module.%s%s", m.Name, m.Index.String())
}

type ResourceSpec struct {
	Type  string
	Name  string
	Index Index
}

func (r *ResourceSpec) String() string {
	return fmt.Sprintf("%s.%s%s", r.Type, r.Name, r.Index.String())
}
