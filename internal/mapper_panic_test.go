package internal

import (
	"testing"

	lua "github.com/yuin/gopher-lua"
)

// TestGoToLuaPanicRecovery tests that GoToLua doesn't panic on unexpected types
func TestGoToLuaPanicRecovery(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	testCases := []struct {
		name  string
		value any
	}{
		{"struct", struct{ Name string }{Name: "test"}},
		{"channel", make(chan int)},
		{"function", func() {}},
		{"complex", complex(1, 2)},
		{"pointer", new(int)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// This should not panic, but return a string representation
			result := GoToLua(L, tc.value)
			
			// Should return a LString with the string representation
			if result.Type() != lua.LTString {
				t.Errorf("Expected LString for %s, got %v", tc.name, result.Type())
			}
		})
	}
}

// TestGoToLuaAllNumericTypes tests that all numeric types are handled correctly
func TestGoToLuaAllNumericTypes(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	testCases := []struct {
		name  string
		value any
	}{
		{"int", int(42)},
		{"int8", int8(42)},
		{"int16", int16(42)},
		{"int32", int32(42)},
		{"int64", int64(42)},
		{"uint", uint(42)},
		{"uint8", uint8(42)},
		{"uint16", uint16(42)},
		{"uint32", uint32(42)},
		{"uint64", uint64(42)},
		{"float32", float32(42.5)},
		{"float64", float64(42.5)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GoToLua(L, tc.value)
			
			// Should return a LNumber
			if result.Type() != lua.LTNumber {
				t.Errorf("Expected LNumber for %s, got %v", tc.name, result.Type())
			}
		})
	}
}

// TestGoToLuaByteSlice tests that []byte is converted to string
func TestGoToLuaByteSlice(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	data := []byte("hello world")
	result := GoToLua(L, data)

	if result.Type() != lua.LTString {
		t.Errorf("Expected LString for []byte, got %v", result.Type())
	}

	if result.String() != "hello world" {
		t.Errorf("Expected 'hello world', got %s", result.String())
	}
}