package args

import (
	"testing"

	"github.com/open-webtech/go-xaddy-config/schema/values"
)

func TestNewArgDef(t *testing.T) {
	var testValue string
	target := values.NewStringValue(&testValue)
	
	argDef := NewArgDef(target, String)
	
	if argDef == nil {
		t.Fatal("NewArgDef() returned nil")
	}
	
	if argDef.Type() != String {
		t.Errorf("Expected type String, got %v", argDef.Type())
	}
	
	if !argDef.Required() {
		t.Error("Expected argument to be required by default")
	}
	
	if argDef.Variadic() {
		t.Error("Expected argument to not be variadic by default")
	}
	
	if argDef.Target() != target {
		t.Error("Expected target to match the provided value")
	}
}

func TestNewArgDefWithOptional(t *testing.T) {
	var testValue string
	target := values.NewStringValue(&testValue)
	
	argDef := NewArgDef(target, String, Optional)
	
	if argDef.Required() {
		t.Error("Expected argument to be optional when Optional attribute is provided")
	}
}

func TestNewVariadicArgDef(t *testing.T) {
	var testValues []string
	target := values.NewStringsValue(&testValues)
	
	argDef := NewVariadicArgDef(target, String)
	
	if argDef == nil {
		t.Fatal("NewVariadicArgDef() returned nil")
	}
	
	if !argDef.Variadic() {
		t.Error("Expected variadic argument to be variadic")
	}
	
	if !argDef.Required() {
		t.Error("Expected variadic argument to be required by default")
	}
}

func TestNewVariadicArgDefWithOptional(t *testing.T) {
	var testValues []string
	target := values.NewStringsValue(&testValues)
	
	argDef := NewVariadicArgDef(target, String, Optional)
	
	if argDef.Required() {
		t.Error("Expected variadic argument to be optional when Optional attribute is provided")
	}
}

func TestArgDefMethods(t *testing.T) {
	var testValue string
	target := values.NewStringValue(&testValue)
	
	argDef := NewArgDef(target, String)
	
	// Test all getter methods
	if argDef.Type() != String {
		t.Errorf("Type() returned %v, expected String", argDef.Type())
	}
	
	if !argDef.Required() {
		t.Error("Required() returned false, expected true")
	}
	
	if argDef.Variadic() {
		t.Error("Variadic() returned true, expected false")
	}
	
	if argDef.Target() != target {
		t.Error("Target() returned different value than expected")
	}
}

func TestStringArg(t *testing.T) {
	var testValue string
	
	argDef := StringArg(&testValue)
	
	if argDef == nil {
		t.Fatal("StringArg() returned nil")
	}
	
	if argDef.Type() != String {
		t.Errorf("Expected type String, got %v", argDef.Type())
	}
	
	if !argDef.Required() {
		t.Error("Expected string argument to be required by default")
	}
	
	// Test setting value through target
	err := argDef.Target().Set("test_string")
	if err != nil {
		t.Errorf("Setting value failed: %v", err)
	}
	
	if testValue != "test_string" {
		t.Errorf("Expected testValue to be 'test_string', got '%s'", testValue)
	}
}

func TestStringArgWithOptional(t *testing.T) {
	var testValue string
	
	argDef := StringArg(&testValue, Optional)
	
	if argDef.Required() {
		t.Error("Expected string argument to be optional when Optional attribute is provided")
	}
}

func TestIntArg(t *testing.T) {
	var testValue int
	
	argDef := IntArg(&testValue)
	
	if argDef == nil {
		t.Fatal("IntArg() returned nil")
	}
	
	if argDef.Type() != Int {
		t.Errorf("Expected type Int, got %v", argDef.Type())
	}
	
	// Test setting value through target
	err := argDef.Target().Set("42")
	if err != nil {
		t.Errorf("Setting integer value failed: %v", err)
	}
	
	if testValue != 42 {
		t.Errorf("Expected testValue to be 42, got %d", testValue)
	}
}

func TestBoolArg(t *testing.T) {
	var testValue bool
	
	argDef := BoolArg(&testValue)
	
	if argDef == nil {
		t.Fatal("BoolArg() returned nil")
	}
	
	if argDef.Type() != Bool {
		t.Errorf("Expected type Bool, got %v", argDef.Type())
	}
	
	// Test setting value through target
	err := argDef.Target().Set("true")
	if err != nil {
		t.Errorf("Setting boolean value failed: %v", err)
	}
	
	if !testValue {
		t.Errorf("Expected testValue to be true, got %t", testValue)
	}
}

func TestVariadicStringArg(t *testing.T) {
	var testValues []string
	
	argDef := VariadicStringArg(&testValues)
	
	if argDef == nil {
		t.Fatal("VariadicStringArg() returned nil")
	}
	
	if !argDef.Variadic() {
		t.Error("Expected variadic string argument to be variadic")
	}
	
	if argDef.Type() != String {
		t.Errorf("Expected type String, got %v", argDef.Type())
	}
}

func TestArgListOperations(t *testing.T) {
	var val1, val2 string
	argsList := ArgsList{
		StringArg(&val1),
		StringArg(&val2, Optional),
	}
	
	if len(argsList) != 2 {
		t.Errorf("Expected ArgsList length 2, got %d", len(argsList))
	}
	
	if !argsList[0].Required() {
		t.Error("Expected first argument to be required")
	}
	
	if argsList[1].Required() {
		t.Error("Expected second argument to be optional")
	}
}

// Test error cases
func TestArgDefWithInvalidValues(t *testing.T) {
	var testValue int
	argDef := IntArg(&testValue)
	
	// Test setting invalid integer value
	err := argDef.Target().Set("not_a_number")
	if err == nil {
		t.Error("Expected error when setting invalid integer value")
	}
}

func TestBoolArgWithInvalidValues(t *testing.T) {
	var testValue bool
	argDef := BoolArg(&testValue)
	
	// Test setting invalid boolean value
	err := argDef.Target().Set("maybe")
	if err == nil {
		t.Error("Expected error when setting invalid boolean value")
	}
}