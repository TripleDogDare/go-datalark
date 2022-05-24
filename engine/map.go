package datalarkengine

import (
	"errors"
	"fmt"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/printer"
	"github.com/ipld/go-ipld-prime/schema"
	"go.starlark.net/starlark"
)

type mapValue struct {
	node datamodel.Node
}

var _ starlark.Mapping = (*mapValue)(nil)
var _ Value = (*mapValue)(nil)

func newMapValue(node datamodel.Node) Value {
	return &mapValue{node}
}

func (v *mapValue) Node() datamodel.Node {
	return v.node
}
func (v *mapValue) Type() string {
	if tn, ok := v.node.(schema.TypedNode); ok {
		return fmt.Sprintf("datalark.Map<%T>", tn.Type().Name())
	}
	return fmt.Sprintf("datalark.Map")
}
func (v *mapValue) String() string {
	return printer.Sprint(v.node)
}
func (v *mapValue) Freeze() {}
func (v *mapValue) Truth() starlark.Bool {
	return true
}
func (v *mapValue) Hash() (uint32, error) {
	return 0, errors.New("TODO")
}

// Get returns a value from a map, implementing starlark.Mapping
// example:
//
//   d = {'a': 'apple', 'b': 'banana'}
//   d['a'] # calls d.Get('a')
//
func (v *mapValue) Get(in starlark.Value) (out starlark.Value, found bool, err error) {
	if _, ok := in.(Value); ok {
		// TODO: unbox it and use LookupByNode.
	}
	// TODO: coerce to string?  (don't use the String method, it's a printer, not what want.)
	// TODO: it has now become high time to standardize the "not found" errors from the Node API!
	ks := in.String() // FIXME placeholder; objectively and clearly wrong.
	n, err := v.node.LookupByString(ks)
	if err != nil {
		return nil, false, err
	}
	val, err := ToValue(n)
	return val, true, err
}

// TODO: Items?  Keys?  Len?  Iterate?  Attr?  AttrNames?
