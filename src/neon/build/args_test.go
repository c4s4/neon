package build

import "testing"

type TestArgs struct {
	Test string
	Num  int `optional`
}

func TestValidateTaskArgsNominal(t *testing.T) {
	argsi := make(map[string]interface{})
	argsi["test"] = "Hello World!"
	argsi["num"] = 3
	err := ValidateTaskArgs(argsi, TestArgs{})
	if err != nil {
		t.Errorf("failed args validation: %#v", err)
	}
}

func TestValidateTaskArgsMissingArg(t *testing.T) {
	argsi := make(map[string]interface{})
	argsi["num"] = 3
	err := ValidateTaskArgs(argsi, TestArgs{})
	if err == nil || err.Error() != "missing mandatory field 'test'" {
		t.Errorf("failed args validation")
	}
}

func TestValidateTaskArgsMissingArgOptional(t *testing.T) {
	argsi := make(map[string]interface{})
	argsi["test"] = "Hello World!"
	err := ValidateTaskArgs(argsi, TestArgs{})
	if err != nil {
		t.Errorf("failed args validation: %#v", err)
	}
}

func TestValidateTaskArgsBadType(t *testing.T) {
	argsi := make(map[string]interface{})
	argsi["test"] = 1
	err := ValidateTaskArgs(argsi, TestArgs{})
	if err == nil || err.Error() != "field 'test' must be of type 'string' ('int' provided)" {
		t.Errorf("failed args validation")
	}
}

func TestEvaluateTaskArgsNominal(t *testing.T) {
	argsi := make(map[string]interface{})
	argsi["test"] = "Hello World!"
	argsi["num"] = 3
	args := TestArgs{}
	err := EvaluateTaskArgs(argsi, &args, nil)
	if err != nil {
		t.Errorf("failed args evaluation: %#v", err)
	}
	if args.Test != "Hello World!" || args.Num != 3 {
		t.Errorf("failed args evaluation: %#v", args)
	}
}
