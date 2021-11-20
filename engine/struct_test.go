package datalarkengine

import (
	"testing"
)

func TestStructs(t *testing.T) {
	evalWithUltramagic(t,
		`
		type FooBar struct {
			foo String
			bar String
		}
	`, `
		print(mytypes.FooBar(foo="one", bar="two"))
	`, `
		struct<FooBar>{
			foo: string<String>{"one"}
			bar: string<String>{"two"}
		}
	`)
}

func TestStructSpecs(t *testing.T) {
	testFixture(t, "../docs/using-structs.md")
}