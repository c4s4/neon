package build

import "testing"

type TestArgs struct {
	Foo  int
	Bar  string `mandatory:"false"`
	Test bool  `name:"toto" mandatory:"false"`
}

func TestParseArgs(t *testing.T) {
	testArgs := TestArgs{Foo: 123, Bar: "test"}
	args, err := ParseArgs(nil, testArgs)
	if err != nil {
		t.Errorf("failed with error: %#v", err)
	}
	println(args)
}
