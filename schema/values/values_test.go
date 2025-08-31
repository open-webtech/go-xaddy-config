package values

import (
	"testing"
)

func TestNewStringValue(t *testing.T) {
	var testString string
	stringValue := NewStringValue(&testString)
	
	if stringValue == nil {
		t.Fatal("NewStringValue() returned nil")
	}
	
	// Test setting value
	err := stringValue.Set("test_value")
	if err != nil {
		t.Errorf("Set() failed: %v", err)
	}
	
	if testString != "test_value" {
		t.Errorf("Expected testString to be 'test_value', got '%s'", testString)
	}
	
	// Test getting value
	got := stringValue.Get()
	if got != "test_value" {
		t.Errorf("Get() returned '%v', expected 'test_value'", got)
	}
	
	// Test string representation
	str := stringValue.String()
	if str != "test_value" {
		t.Errorf("String() returned '%s', expected 'test_value'", str)
	}
}

func TestNewBoolValue(t *testing.T) {
	var testBool bool
	boolValue := NewBoolValue(&testBool)
	
	if boolValue == nil {
		t.Fatal("NewBoolValue() returned nil")
	}
	
	// Test setting true
	err := boolValue.Set("true")
	if err != nil {
		t.Errorf("Set('true') failed: %v", err)
	}
	
	if !testBool {
		t.Error("Expected testBool to be true")
	}
	
	// Test setting false
	err = boolValue.Set("false")
	if err != nil {
		t.Errorf("Set('false') failed: %v", err)
	}
	
	if testBool {
		t.Error("Expected testBool to be false")
	}
	
	// Test invalid boolean
	err = boolValue.Set("maybe")
	if err == nil {
		t.Error("Expected error when setting invalid boolean value")
	}
	
	// Test getting value
	got := boolValue.Get()
	if got != false {
		t.Errorf("Get() returned %v, expected false", got)
	}
}

func TestNewIntValue(t *testing.T) {
	var testInt int
	intValue := NewIntValue(&testInt)
	
	if intValue == nil {
		t.Fatal("NewIntValue() returned nil")
	}
	
	// Test setting positive integer
	err := intValue.Set("42")
	if err != nil {
		t.Errorf("Set('42') failed: %v", err)
	}
	
	if testInt != 42 {
		t.Errorf("Expected testInt to be 42, got %d", testInt)
	}
	
	// Test setting negative integer
	err = intValue.Set("-123")
	if err != nil {
		t.Errorf("Set('-123') failed: %v", err)
	}
	
	if testInt != -123 {
		t.Errorf("Expected testInt to be -123, got %d", testInt)
	}
	
	// Test setting zero
	err = intValue.Set("0")
	if err != nil {
		t.Errorf("Set('0') failed: %v", err)
	}
	
	if testInt != 0 {
		t.Errorf("Expected testInt to be 0, got %d", testInt)
	}
	
	// Test invalid integer
	err = intValue.Set("not_a_number")
	if err == nil {
		t.Error("Expected error when setting invalid integer value")
	}
	
	// Test getting value
	got := intValue.Get()
	if got != 0 {
		t.Errorf("Get() returned %v, expected 0", got)
	}
}

func TestNewUintValue(t *testing.T) {
	var testUint uint
	uintValue := NewUintValue(&testUint)
	
	if uintValue == nil {
		t.Fatal("NewUintValue() returned nil")
	}
	
	// Test setting positive integer
	err := uintValue.Set("42")
	if err != nil {
		t.Errorf("Set('42') failed: %v", err)
	}
	
	if testUint != 42 {
		t.Errorf("Expected testUint to be 42, got %d", testUint)
	}
	
	// Test setting zero
	err = uintValue.Set("0")
	if err != nil {
		t.Errorf("Set('0') failed: %v", err)
	}
	
	if testUint != 0 {
		t.Errorf("Expected testUint to be 0, got %d", testUint)
	}
	
	// Test invalid unsigned integer
	err = uintValue.Set("-123")
	if err == nil {
		t.Error("Expected error when setting negative value for unsigned integer")
	}
	
	// Test getting value
	got := uintValue.Get()
	if got != uint(0) {
		t.Errorf("Get() returned %v, expected 0", got)
	}
}

func TestNewFloat32Value(t *testing.T) {
	var testFloat float32
	floatValue := NewFloat32Value(&testFloat)
	
	if floatValue == nil {
		t.Fatal("NewFloat32Value() returned nil")
	}
	
	// Test setting float value
	err := floatValue.Set("3.14")
	if err != nil {
		t.Errorf("Set('3.14') failed: %v", err)
	}
	
	if testFloat != 3.14 {
		t.Errorf("Expected testFloat to be 3.14, got %f", testFloat)
	}
	
	// Test setting negative float
	err = floatValue.Set("-2.5")
	if err != nil {
		t.Errorf("Set('-2.5') failed: %v", err)
	}
	
	if testFloat != -2.5 {
		t.Errorf("Expected testFloat to be -2.5, got %f", testFloat)
	}
	
	// Test invalid float
	err = floatValue.Set("not_a_float")
	if err == nil {
		t.Error("Expected error when setting invalid float value")
	}
}

func TestNewFloat64Value(t *testing.T) {
	var testFloat float64
	floatValue := NewFloat64Value(&testFloat)
	
	if floatValue == nil {
		t.Fatal("NewFloat64Value() returned nil")
	}
	
	// Test setting float value
	err := floatValue.Set("3.141592653589793")
	if err != nil {
		t.Errorf("Set() failed: %v", err)
	}
	
	expected := 3.141592653589793
	if testFloat != expected {
		t.Errorf("Expected testFloat to be %f, got %f", expected, testFloat)
	}
}

func TestNewStringsValue(t *testing.T) {
	var testStrings []string
	stringsValue := NewStringsValue(&testStrings)
	
	if stringsValue == nil {
		t.Fatal("NewStringsValue() returned nil")
	}
	
	// Test setting multiple values
	values := []string{"value1", "value2", "value3"}
	SetList(stringsValue, values)
	
	if len(testStrings) != 3 {
		t.Errorf("Expected 3 strings, got %d", len(testStrings))
	}
	
	for i, expected := range values {
		if testStrings[i] != expected {
			t.Errorf("Expected testStrings[%d] to be '%s', got '%s'", i, expected, testStrings[i])
		}
	}
}

func TestNewBoolsValue(t *testing.T) {
	var testBools []bool
	boolsValue := NewBoolsValue(&testBools)
	
	if boolsValue == nil {
		t.Fatal("NewBoolsValue() returned nil")
	}
	
	// Test setting multiple boolean values
	values := []string{"true", "false", "true"}
	SetList(boolsValue, values)
	
	if len(testBools) != 3 {
		t.Errorf("Expected 3 booleans, got %d", len(testBools))
	}
	
	expectedBools := []bool{true, false, true}
	for i, expected := range expectedBools {
		if testBools[i] != expected {
			t.Errorf("Expected testBools[%d] to be %t, got %t", i, expected, testBools[i])
		}
	}
}

func TestNewIntsValue(t *testing.T) {
	var testInts []int
	intsValue := NewIntsValue(&testInts)
	
	if intsValue == nil {
		t.Fatal("NewIntsValue() returned nil")
	}
	
	// Test setting multiple integer values
	values := []string{"1", "2", "3"}
	SetList(intsValue, values)
	
	if len(testInts) != 3 {
		t.Errorf("Expected 3 integers, got %d", len(testInts))
	}
	
	expectedInts := []int{1, 2, 3}
	for i, expected := range expectedInts {
		if testInts[i] != expected {
			t.Errorf("Expected testInts[%d] to be %d, got %d", i, expected, testInts[i])
		}
	}
}

func TestNewUintsValue(t *testing.T) {
	var testUints []uint
	uintsValue := NewUintsValue(&testUints)
	
	if uintsValue == nil {
		t.Fatal("NewUintsValue() returned nil")
	}
	
	// Test setting multiple unsigned integer values
	values := []string{"1", "2", "3"}
	SetList(uintsValue, values)
	
	if len(testUints) != 3 {
		t.Errorf("Expected 3 unsigned integers, got %d", len(testUints))
	}
	
	expectedUints := []uint{1, 2, 3}
	for i, expected := range expectedUints {
		if testUints[i] != expected {
			t.Errorf("Expected testUints[%d] to be %d, got %d", i, expected, testUints[i])
		}
	}
}

func TestSetListWithEmptySlice(t *testing.T) {
	var testStrings []string
	stringsValue := NewStringsValue(&testStrings)
	
	// Test setting empty list
	values := []string{}
	SetList(stringsValue, values)
	
	if len(testStrings) != 0 {
		t.Errorf("Expected 0 strings, got %d", len(testStrings))
	}
}

func TestAccumulator(t *testing.T) {
	var testStrings []string
	
	accumulator := NewAccumulator(&testStrings, func(v interface{}) Value {
		return NewStringValue(v.(*string))
	})
	
	if accumulator == nil {
		t.Fatal("NewAccumulator() returned nil")
	}
	
	// Test accumulating values
	values := []string{"val1", "val2", "val3"}
	SetList(accumulator, values)
	
	if len(testStrings) != 3 {
		t.Errorf("Expected 3 accumulated strings, got %d", len(testStrings))
	}
	
	for i, expected := range values {
		if testStrings[i] != expected {
			t.Errorf("Expected testStrings[%d] to be '%s', got '%s'", i, expected, testStrings[i])
		}
	}
}