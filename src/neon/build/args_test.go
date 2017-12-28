package build

import (
	"testing"
	"reflect"
)

type TestArgs struct {
	Bool   bool    `optional`
	Int    int     `optional`
	Float  float64 `optional`
	String string
}

func TestValidateTaskArgsNominal(t *testing.T) {
	argsi := make(map[string]interface{})
	argsi["bool"] = true
	argsi["int"] = 3
	argsi["float"] = 3.14
	argsi["string"] = "Hello World!"
	err := ValidateTaskArgs(argsi, TestArgs{})
	if err != nil {
		t.Errorf("failed args validation: %#v", err)
	}
}

func TestValidateTaskArgsMissingArg(t *testing.T) {
	argsi := make(map[string]interface{})
	argsi["int"] = 3
	err := ValidateTaskArgs(argsi, TestArgs{})
	if err == nil || err.Error() != "missing mandatory field 'string'" {
		t.Errorf("failed args validation: %v", err)
	}
}

func TestValidateTaskArgsMissingArgOptional(t *testing.T) {
	argsi := make(map[string]interface{})
	argsi["string"] = "Hello World!"
	err := ValidateTaskArgs(argsi, TestArgs{})
	if err != nil {
		t.Errorf("failed args validation: %#v", err)
	}
}

func TestValidateTaskArgsBadType(t *testing.T) {
	argsi := make(map[string]interface{})
	argsi["string"] = 1
	err := ValidateTaskArgs(argsi, TestArgs{})
	if err == nil || err.Error() != "field 'string' must be of type 'string' ('int' provided)" {
		t.Errorf("failed args validation")
	}
}

func TestEvaluateTaskArgsNominal(t *testing.T) {
	argsi := make(map[string]interface{})
	argsi["bool"] = true
	argsi["int"] = 3
	argsi["string"] = "Hello World!"
	args := TestArgs{}
	err := EvaluateTaskArgs(argsi, &args, nil)
	if err != nil {
		t.Errorf("failed args evaluation: %#v", err)
	}
	if args.String != "Hello World!" || args.Int != 3 || !args.Bool {
		t.Errorf("failed args evaluation: %#v", args)
	}
}

func TestFieldIs(t *testing.T) {
	field := reflect.StructField{Tag: "test"}
	if !FieldIs(field, "test") {
		t.Errorf("failed FieldIs test")
	}
	if FieldIs(field, "foo") {
		t.Errorf("failed FieldIs test")
	}
	field = reflect.StructField{Tag: "foo bar"}
	if !FieldIs(field, "foo") || !FieldIs(field, "bar") {
		t.Errorf("failed FieldIs test")
	}
	if FieldIs(field, "test") {
		t.Errorf("failed FieldIs test")
	}
}