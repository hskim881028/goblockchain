package utility

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestHandleErr(t *testing.T) {
	preLogPanic := logPanic
	defer func() {
		logPanic = preLogPanic
	}()

	called := false
	logPanic = func(v ...interface{}) {
		called = true
	}

	err := errors.New("Test")
	HandleError(err)
	if !called {
		t.Error("HandleError should call logPanic")
	}
}

func TestToBytes(t *testing.T) {
	s := "test_word"
	b := ToBytes(s)
	k := reflect.TypeOf(b).Kind()
	if k != reflect.Slice {
		t.Errorf("ToBytes should return a slice of bytes got %s", k)
	}
}

func TestFromBytes(t *testing.T) {
	type testStruct struct {
		Test string
	}

	var restored testStruct
	ts := testStruct{"test"}
	b := ToBytes(ts)
	FromBytes(&restored, b)
	if !reflect.DeepEqual(ts, restored) {
		t.Errorf("FromBytes should restore struct")
	}
}

func TestHash(t *testing.T) {
	hash := "8757833d76c277fc3112bec45bfa35ea7f88e7627dc1e9c5ddd74bfd26aa6f6a"
	s := struct{ Test string }{Test: "test_word"}
	x := Hash(s)
	t.Run("Hash is always same", func(t *testing.T) {
		if x != hash {
			t.Errorf("Expected : %s, got : %s", hash, x)
		}
	})

	t.Run("Hash is hex encoded", func(t *testing.T) {
		_, err := hex.DecodeString(x)
		if err != nil {
			t.Error("Hash should be hex encoded")
		}
	})
}

func TestToJson(t *testing.T) {
	type testStruct struct {
		Test string
	}

	ts := testStruct{"test"}
	b := ToJson(ts)
	k := reflect.TypeOf(b).Kind()
	if k != reflect.Slice {
		t.Errorf("Expected %v and get %v", reflect.Slice, k)
	}

	var restored testStruct
	json.Unmarshal(b, &restored)
	if !reflect.DeepEqual(ts, restored) {
		t.Error("ToJson should encode to JSON")
	}
}

func ExampleHash() {
	s := struct{ Test string }{Test: "test_word"}
	x := Hash(s)
	fmt.Println(x)
	// Output:8757833d76c277fc3112bec45bfa35ea7f88e7627dc1e9c5ddd74bfd26aa6f6a
}

func TestSplitter(t *testing.T) {
	type test struct {
		input  string
		sep    string
		index  int
		output string
	}

	tests := []test{
		{"0:1:2", ":", 0, "0"},
		{"0:1:2", ":", 1, "1"},
		{"0:1:2", ":", 2, "2"},
		{"0:1:2", ":", 3, ""},
		{"0:1:2", "?", 0, "0:1:2"},
	}

	for _, tc := range tests {
		result := Splitter(tc.input, tc.sep, tc.index)
		if result != tc.output {
			t.Errorf("Expected %s and get %s", tc.output, result)
		}
	}
}
