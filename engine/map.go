package datalarkengine

import (
	"errors"
	"fmt"

	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/ipld/go-ipld-prime/printer"
	"github.com/ipld/go-ipld-prime/schema"
	"go.starlark.net/starlark"
)

func ConstructMap(np datamodel.NodePrototype, _ *starlark.Thread, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	// TODO should have several different options:
	//  - accept a single arg of a starlark dict and bounce it in.
	//    - maybe even multiple, and merge them?  maybe.
	//  - accept as many kwargs as you like.
	//  - positional args?  eh, hrm.  maybe.  or maybe not.  i dislike pairs-styled interfaces.
	//  - list of tuples?  does starlark have tuples?  idk.  if this is syntactically obvious, then sure, let's have it.
	// TODO I think there may be some destructing syntax (`**kwargs`...?) that we should understand before going wild with this.
	//  update:
	//    - these are medium powerful.  they do let you unpack a dict.  nice.
	//    - "keyword argument may not follow **kwargs"... so you can't use them to easily override some values.
	//    - this syntax actually allows you to sneak the same key in twice (probably a bug!).
	//    - you can't use kwargs for all strings, so in general, no we can't use kwargs as the exclusive option for any of this.
	//    - yes, you can use `**{"ay":"bee"}` and it works.
	//    - yes, you can use `**uwot()` on a function that returns a dict and it works.
	//    - no, you can not use more than one doublestar in the same function invocaton.
	//    - unknown if you can use doublestar on something that's iterable (or whatever) but not literally a built-in dict.  (probably?)
	//    - tangentally: no, i'm pretty sure you can't add a '+' binary operator on the built-in dict type.

	// FUTURE: (*far* future...) We could also try to accept a `starlark.Callable` as a positional arg, and hand it assemblers.  May be nice for performance since it can build in-place and do less allocs and less copying (same reasons as in direct golang).

	if len(args) > 0 && len(kwargs) > 0 {
		return starlark.None, fmt.Errorf("datalark.Map: can either use positional or keyword arguments, but not both")
	}

	nb := np.NewBuilder()
	if len(args) > 1 {
		return starlark.None, fmt.Errorf("datalark.Map: if using positional arguments, only one is expected")
	}
	if len(args) == 1 {
		// If there's one arg, and it's a starlark dict, 'assignish' will do the right thing and restructure that into us.
		// FUTURE: I'd like this to also be clever enough to, if the map key type is a struct or something (but with stringy representation), have this also magic that into rightness.  Unclear where exactly that magic should live, though.
		if err := assignish(nb, args[0]); err != nil {
			return starlark.None, fmt.Errorf("datalark.Map: %w", err)
		}
	}
	if len(kwargs) > 0 {
		ma, err := nb.BeginMap(int64(len(kwargs))) // FUTURE: this could... need to take into account more things.
		if err != nil {
			return starlark.None, err
		}
		for _, kwarg := range kwargs {
			if err := assignish(ma.AssembleKey(), kwarg[0]); err != nil {
				return starlark.None, err
			}
			if err := assignish(ma.AssembleValue(), kwarg[1]); err != nil {
				return starlark.None, err
			}
		}
		if err := ma.Finish(); err != nil {
			return starlark.None, err
		}
	}
	return &Map{nb.Build()}, nil
}

var _ starlark.Mapping = (*Map)(nil)

type Map struct {
	val datamodel.Node
}

func (g *Map) Node() datamodel.Node {
	return g.val
}
func (g *Map) Type() string {
	if tn, ok := g.val.(schema.TypedNode); ok {
		return fmt.Sprintf("datalark.Map<%T>", tn.Type().Name())
	}
	return fmt.Sprintf("datalark.Map")
}
func (g *Map) String() string {
	return printer.Sprint(g.val)
}
func (g *Map) Freeze() {}
func (g *Map) Truth() starlark.Bool {
	return true
}
func (g *Map) Hash() (uint32, error) {
	return 0, errors.New("TODO")
}

// Get implements part of `starlark.Mapping`.
func (g *Map) Get(in starlark.Value) (out starlark.Value, found bool, err error) {
	if _, ok := in.(Value); ok {
		// TODO: unbox it and use LookupByNode.
	}
	// TODO: coerce to string?  (don't use the String method, it's a printer, not what want.)
	// TODO: it has now become high time to standardize the "not found" errors from the Node API!
	ks := in.String() // FIXME placeholder; objectively and clearly wrong.
	n, err := g.val.LookupByString(ks)
	if err != nil {
		return nil, false, err
	}
	w, err := Wrap(n)
	return w, true, err
}

// TODO: Items?  Keys?  Len?  Iterate?  Attr?  AttrNames?
