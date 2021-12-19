package utility

import (
	"encoding/hex"
	"fmt"
	"testing"
)

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

func ExampleHash() {
	s := struct{ Test string }{Test: "test_word"}
	x := Hash(s)
	fmt.Println(x)
	// Output:8757833d76c277fc3112bec45bfa35ea7f88e7627dc1e9c5ddd74bfd26aa6f6a
}
