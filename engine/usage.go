package datalarkengine

import (
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/ipld/go-ipld-prime/schema"
	"go.starlark.net/starlark"
)

// See docs on datalark.InjectGlobals.
// Typically you should prefer using functions in the datalark package,
// rather than their equivalents in the datalarkengine package.
func InjectGlobals(globals starlark.StringDict, obj *Object) {
	// Technically this would work on any 'starlark.IterableMapping', but I don't think that makes the function more useful, and would make it *less* self-documenting.
	itr := obj.Iterate()
	defer itr.Done()
	var k starlark.Value
	for itr.Next(&k) {
		v, _, err := obj.Get(k)
		if err != nil {
			panic(err)
		}
		globals[string(k.(starlark.String))] = v
	}
}

// See docs on datalark.ObjOfConstructorsForPrimitives.
// Typically you should prefer using functions in the datalark package,
// rather than their equivalents in the datalarkengine package.
func ObjOfConstructorsForPrimitives() *Object {
	obj := NewObject(7)
	obj.SetKey(starlark.String("Map"), &Prototype{basicnode.Prototype.Map})
	obj.SetKey(starlark.String("List"), &Prototype{basicnode.Prototype.List})
	obj.SetKey(starlark.String("Bool"), &Prototype{basicnode.Prototype.Bool})
	obj.SetKey(starlark.String("Int"), &Prototype{basicnode.Prototype.Int})
	obj.SetKey(starlark.String("Float"), &Prototype{basicnode.Prototype.Float})
	obj.SetKey(starlark.String("String"), &Prototype{basicnode.Prototype.String})
	obj.SetKey(starlark.String("Bytes"), &Prototype{basicnode.Prototype.Bytes})
	obj.Freeze()
	return obj
}

// See docs on datalark.ObjOfConstructorsForPrototypes.
// Typically you should prefer using functions in the datalark package,
// rather than their equivalents in the datalarkengine package.
func ObjOfConstructorsForPrototypes(prototypes ...schema.TypedPrototype) *Object {
	obj := NewObject(len(prototypes))
	for _, npt := range prototypes {
		obj.SetKey(starlark.String(npt.Type().Name()), &Prototype{npt})
	}
	obj.Freeze()
	return obj
}
