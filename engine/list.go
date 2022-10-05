package datalarkengine

import (
	"errors"
	"fmt"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/printer"
	"go.starlark.net/starlark"
)

type listValue struct {
	node datamodel.Node
}

var _ Value = (*listValue)(nil)
var _ starlark.Sequence = (*listValue)(nil)

func newListValue(node datamodel.Node) Value {
	return &listValue{node}
}

func (v *listValue) Node() datamodel.Node {
	return v.node
}
func (v *listValue) Type() string {
	// TODO(dustmop): Can a list be a TypedNode? I believe so, it
	// is used for a homogeneous typed list.
	return fmt.Sprintf("datalark.List")
}
func (v *listValue) String() string {
	return printer.Sprint(v.node)
}
func (v *listValue) Freeze() {}
func (v *listValue) Truth() starlark.Bool {
	return true
}
func (v *listValue) Hash() (uint32, error) {
	return 0, errors.New("TODO")
}

// NewList converts a starlark.List into a datalark.Value
func NewList(starList *starlark.List) (Value, error) {
	nb := basicnode.Prototype.List.NewBuilder()
	size := starList.Len()
	la, err := nb.BeginList(int64(size))
	if err != nil {
		return nil, err
	}
	for i := 0; i < size; i++ {
		item := starList.Index(i)
		if err := assembleFrom(la.AssembleValue(), item); err != nil {
			return nil, fmt.Errorf("cannot add %v of type %T", item, item)
		}
	}
	if err := la.Finish(); err != nil {
		return nil, err
	}
	return newListValue(nb.Build()), nil
}

// starlark.Sequence

func (v *listValue) Iterate() starlark.Iterator {
	panic(fmt.Errorf("TODO(dustmop): listValue.Iterate not implemented for %T", v))
}

func (v *listValue) Len() int {
	return int(v.node.Length())
}

// starlark.HasAttrs : starlark.List

var listMethods = []string{"clear", "copy", "fromkeys", "get", "items", "keys", "pop", "popitem", "setdefault", "update", "values"}

func (v *listValue) Attr(gstrName string) (starlark.Value, error) {
	return starlark.None, nil
}

func (v *listValue) AttrNames() []string {
	return listMethods
}
