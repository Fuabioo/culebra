package internal

import (
	"reflect"
	"testing"

	lua "github.com/yuin/gopher-lua"
)

func TestIsLuaArray(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tests := []struct {
		name     string
		setup    func() *lua.LTable
		expected bool
	}{
		{
			name: "ValidArray",
			setup: func() *lua.LTable {
				table := L.NewTable()
				table.RawSetInt(1, lua.LString("first"))
				table.RawSetInt(2, lua.LString("second"))
				table.RawSetInt(3, lua.LString("third"))
				return table
			},
			expected: true,
		},
		{
			name: "EmptyTable",
			setup: func() *lua.LTable {
				return L.NewTable()
			},
			expected: false,
		},
		{
			name: "SparseArray",
			setup: func() *lua.LTable {
				table := L.NewTable()
				table.RawSetInt(1, lua.LString("first"))
				table.RawSetInt(3, lua.LString("third")) // Missing index 2
				return table
			},
			expected: false,
		},
		{
			name: "MixedKeys",
			setup: func() *lua.LTable {
				table := L.NewTable()
				table.RawSetInt(1, lua.LString("first"))
				table.RawSetString("key", lua.LString("value")) // Has string key
				return table
			},
			expected: false,
		},
		{
			name: "NonSequentialStart",
			setup: func() *lua.LTable {
				table := L.NewTable()
				table.RawSetInt(0, lua.LString("zero")) // Starts at 0, not 1
				table.RawSetInt(1, lua.LString("one"))
				return table
			},
			expected: false,
		},
		{
			name: "FloatIndices",
			setup: func() *lua.LTable {
				table := L.NewTable()
				L.SetTable(table, lua.LNumber(1.5), lua.LString("float"))
				return table
			},
			expected: false,
		},
		{
			name: "SingleElement",
			setup: func() *lua.LTable {
				table := L.NewTable()
				table.RawSetInt(1, lua.LString("only"))
				return table
			},
			expected: true,
		},
		{
			name: "LargeValidArray",
			setup: func() *lua.LTable {
				table := L.NewTable()
				for i := 1; i <= 100; i++ {
					table.RawSetInt(i, lua.LNumber(i))
				}
				return table
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := tt.setup()
			result := isLuaArray(table)
			if result != tt.expected {
				t.Errorf("isLuaArray() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHasSequentialKeys(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tests := []struct {
		name     string
		setup    func() (*lua.LTable, int)
		expected bool
	}{
		{
			name: "Sequential",
			setup: func() (*lua.LTable, int) {
				table := L.NewTable()
				table.RawSetInt(1, lua.LString("a"))
				table.RawSetInt(2, lua.LString("b"))
				table.RawSetInt(3, lua.LString("c"))
				return table, 3
			},
			expected: true,
		},
		{
			name: "MissingKey",
			setup: func() (*lua.LTable, int) {
				table := L.NewTable()
				table.RawSetInt(1, lua.LString("a"))
				table.RawSetInt(3, lua.LString("c"))
				return table, 3
			},
			expected: false,
		},
		{
			name: "EmptyWithZeroLength",
			setup: func() (*lua.LTable, int) {
				table := L.NewTable()
				return table, 0
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table, length := tt.setup()
			result := hasSequentialKeys(table, length)
			if result != tt.expected {
				t.Errorf("hasSequentialKeys() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHasNonArrayKeys(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tests := []struct {
		name     string
		setup    func() (*lua.LTable, int)
		expected bool
	}{
		{
			name: "OnlyArrayKeys",
			setup: func() (*lua.LTable, int) {
				table := L.NewTable()
				table.RawSetInt(1, lua.LString("a"))
				table.RawSetInt(2, lua.LString("b"))
				return table, 2
			},
			expected: false,
		},
		{
			name: "HasStringKey",
			setup: func() (*lua.LTable, int) {
				table := L.NewTable()
				table.RawSetInt(1, lua.LString("a"))
				table.RawSetString("key", lua.LString("value"))
				return table, 1
			},
			expected: true,
		},
		{
			name: "HasOutOfRangeIndex",
			setup: func() (*lua.LTable, int) {
				table := L.NewTable()
				table.RawSetInt(1, lua.LString("a"))
				table.RawSetInt(5, lua.LString("e")) // Index 5 when length is 1
				return table, 1
			},
			expected: true,
		},
		{
			name: "HasZeroIndex",
			setup: func() (*lua.LTable, int) {
				table := L.NewTable()
				table.RawSetInt(0, lua.LString("zero"))
				table.RawSetInt(1, lua.LString("one"))
				return table, 1
			},
			expected: true,
		},
		{
			name: "HasFloatIndex",
			setup: func() (*lua.LTable, int) {
				table := L.NewTable()
				table.RawSetInt(1, lua.LString("one"))
				L.SetTable(table, lua.LNumber(1.5), lua.LString("float"))
				return table, 1
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table, length := tt.setup()
			result := hasNonArrayKeys(table, length)
			if result != tt.expected {
				t.Errorf("hasNonArrayKeys() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestLuaTableToGoSlice(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tests := []struct {
		name     string
		setup    func() *lua.LTable
		expected []any
	}{
		{
			name: "StringArray",
			setup: func() *lua.LTable {
				table := L.NewTable()
				table.RawSetInt(1, lua.LString("first"))
				table.RawSetInt(2, lua.LString("second"))
				table.RawSetInt(3, lua.LString("third"))
				return table
			},
			expected: []any{"first", "second", "third"},
		},
		{
			name: "NumberArray",
			setup: func() *lua.LTable {
				table := L.NewTable()
				table.RawSetInt(1, lua.LNumber(10))
				table.RawSetInt(2, lua.LNumber(20))
				table.RawSetInt(3, lua.LNumber(30))
				return table
			},
			expected: []any{float64(10), float64(20), float64(30)},
		},
		{
			name: "MixedTypeArray",
			setup: func() *lua.LTable {
				table := L.NewTable()
				table.RawSetInt(1, lua.LString("string"))
				table.RawSetInt(2, lua.LNumber(42))
				table.RawSetInt(3, lua.LBool(true))
				return table
			},
			expected: []any{"string", float64(42), true},
		},
		{
			name: "NestedArray",
			setup: func() *lua.LTable {
				innerTable := L.NewTable()
				innerTable.RawSetInt(1, lua.LString("nested"))

				table := L.NewTable()
				table.RawSetInt(1, lua.LString("first"))
				table.RawSetInt(2, innerTable)
				return table
			},
			expected: []any{"first", []any{"nested"}},
		},
		{
			name: "SingleElement",
			setup: func() *lua.LTable {
				table := L.NewTable()
				table.RawSetInt(1, lua.LString("only"))
				return table
			},
			expected: []any{"only"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := tt.setup()
			result := luaTableToGoSlice(table, true)

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("luaTableToGoSlice() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGoSliceToLuaTable(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	tests := []struct {
		name  string
		input []any
		check func(*lua.LTable) bool
	}{
		{
			name:  "StringSlice",
			input: []any{"a", "b", "c"},
			check: func(table *lua.LTable) bool {
				return table.RawGetInt(1).String() == "a" &&
					table.RawGetInt(2).String() == "b" &&
					table.RawGetInt(3).String() == "c" &&
					table.Len() == 3
			},
		},
		{
			name:  "NumberSlice",
			input: []any{1, 2.5, 3},
			check: func(table *lua.LTable) bool {
				v1, _ := table.RawGetInt(1).(lua.LNumber)
				v2, _ := table.RawGetInt(2).(lua.LNumber)
				v3, _ := table.RawGetInt(3).(lua.LNumber)
				return float64(v1) == 1 &&
					float64(v2) == 2.5 &&
					float64(v3) == 3
			},
		},
		{
			name:  "MixedSlice",
			input: []any{"string", 42, true, nil},
			check: func(table *lua.LTable) bool {
				return table.RawGetInt(1).String() == "string" &&
					table.RawGetInt(2).(lua.LNumber) == 42 &&
					table.RawGetInt(3).(lua.LBool) == true &&
					table.RawGetInt(4) == lua.LNil
			},
		},
		{
			name:  "EmptySlice",
			input: []any{},
			check: func(table *lua.LTable) bool {
				return table.Len() == 0
			},
		},
		{
			name:  "NestedSlice",
			input: []any{"outer", []any{"inner1", "inner2"}},
			check: func(table *lua.LTable) bool {
				if table.RawGetInt(1).String() != "outer" {
					return false
				}
				innerTable, ok := table.RawGetInt(2).(*lua.LTable)
				if !ok {
					return false
				}
				return innerTable.RawGetInt(1).String() == "inner1" &&
					innerTable.RawGetInt(2).String() == "inner2"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := goSliceToLuaTable(L, tt.input)
			if !tt.check(result) {
				t.Errorf("goSliceToLuaTable() failed for %s", tt.name)
			}
		})
	}
}

func TestLuaToGoWithConfigArrayHandling(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Test table conversion with array conversion enabled
	t.Run("TableWithArrayConversion", func(t *testing.T) {
		table := L.NewTable()
		table.RawSetInt(1, lua.LString("first"))
		table.RawSetInt(2, lua.LString("second"))

		result := LuaToGoWithConfig(table, true)
		slice, ok := result.([]any)
		if !ok {
			t.Fatalf("Expected []any, got %T", result)
		}

		if len(slice) != 2 || slice[0] != "first" || slice[1] != "second" {
			t.Errorf("Unexpected slice: %v", slice)
		}
	})

	// Test table conversion without array conversion
	t.Run("TableWithoutArrayConversion", func(t *testing.T) {
		table := L.NewTable()
		table.RawSetInt(1, lua.LString("first"))
		table.RawSetInt(2, lua.LString("second"))

		result := LuaToGoWithConfig(table, false)
		m, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("Expected map[string]any, got %T", result)
		}

		if m["1"] != "first" || m["2"] != "second" {
			t.Errorf("Unexpected map: %v", m)
		}
	})

	// Test non-array table with array conversion enabled
	t.Run("NonArrayTableWithConversion", func(t *testing.T) {
		table := L.NewTable()
		table.RawSetString("key", lua.LString("value"))

		result := LuaToGoWithConfig(table, true)
		m, ok := result.(map[string]any)
		if !ok {
			t.Fatalf("Expected map[string]any, got %T", result)
		}

		if m["key"] != "value" {
			t.Errorf("Unexpected map: %v", m)
		}
	})

	// Test other Lua types
	t.Run("OtherTypes", func(t *testing.T) {
		// These should not be affected by convertArrays flag

		// Default case - unknown type
		userdata := L.NewUserData()
		result := LuaToGoWithConfig(userdata, true)
		if result == nil {
			t.Error("Expected non-nil for unknown type")
		}
	})
}

func TestGoToLuaWithSlices(t *testing.T) {
	L := lua.NewState()
	defer L.Close()

	// Test []any conversion
	t.Run("SliceAny", func(t *testing.T) {
		slice := []any{"a", 1, true}
		result := GoToLua(L, slice)

		table, ok := result.(*lua.LTable)
		if !ok {
			t.Fatalf("Expected *lua.LTable, got %T", result)
		}

		if table.RawGetInt(1).String() != "a" {
			t.Error("First element mismatch")
		}
		if n, _ := table.RawGetInt(2).(lua.LNumber); float64(n) != 1 {
			t.Error("Second element mismatch")
		}
		if b, _ := table.RawGetInt(3).(lua.LBool); bool(b) != true {
			t.Error("Third element mismatch")
		}
	})

	// Test map[string]any conversion
	t.Run("MapStringAny", func(t *testing.T) {
		m := map[string]any{"key": "value", "num": 42}
		result := GoToLua(L, m)

		table, ok := result.(*lua.LTable)
		if !ok {
			t.Fatalf("Expected *lua.LTable, got %T", result)
		}

		if table.RawGetString("key").String() != "value" {
			t.Error("String value mismatch")
		}
		if n, _ := table.RawGetString("num").(lua.LNumber); float64(n) != 42 {
			t.Error("Number value mismatch")
		}
	})
}
