package utils

import (
	"reflect"
	"testing"
)

type testStruct struct {
	Field1 int
	Field2 string
}

func TestIfElse(t *testing.T) {
	tt := []struct {
		name      string
		condition bool
		trueVal   interface{}
		falseVal  interface{}
		expected  interface{}
	}{
		{"TestTrue", true, "True case", "False case", "True case"},
		{"TestFalse", false, "True case", "False case", "False case"},
		{"TestNumbers", true, 123, 456, 123},
		{"TestFunctionTrue", testBooleanCondition(true), "True case", "False case", "True case"},
		{"TestFunctionFalse", testBooleanCondition(false), "True case", "False case", "False case"},
		{"TestStringEqualityTrue", "Hans" == "Hans", "equal", "not equal", "equal"},
		{"TestStringEqualityFalse", "Hans" == "Karl", "equal", "not equal", "not equal"},
		{"TestMap", true, map[string]int{"one": 1}, map[string]int{"two": 2}, map[string]int{"one": 1}},
		{"TestStruct", false, testStruct{1, "one"}, testStruct{2, "two"}, testStruct{2, "two"}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := IfElse(tc.condition, tc.trueVal, tc.falseVal)
			checkExpected(t, tc.expected, result)
		})
	}
}

func testBooleanCondition(val bool) bool {
	return val
}

func checkExpected(t *testing.T, expected, actual interface{}) {
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("Expected %v, but got %v", expected, actual)
	}
}
