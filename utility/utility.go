// Package utility contains functions to be used across the application.
package utility

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

var logPanic = log.Panic

func HandleError(err error) {
	if err != nil {
		logPanic(err)
	}
}

func ToBytes(i interface{}) []byte {
	var aBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBuffer)
	HandleError(encoder.Encode(i))
	return aBuffer.Bytes()
}

func FromBytes(i interface{}, data []byte) {
	decoder := gob.NewDecoder(bytes.NewBuffer(data))
	HandleError(decoder.Decode(i))
}

func ToJson(i interface{}) []byte {
	b, err := json.Marshal(i)
	HandleError(err)
	return b
}

// Hash takes an interface, hashes it and returns the hex encoding the data to the interface.
func Hash(i interface{}) string {
	s := fmt.Sprintf("%v", i)
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

func Splitter(s, sep string, i int) string {
	r := strings.Split(s, sep)
	if len(r)-1 < i {
		return ""
	}
	return r[i]
}
