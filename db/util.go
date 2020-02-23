package db

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

const (
	dataFile       = "data.dat"
	rootPosition   = 0
	maxstoragesize = 1     //GB
	capsize        = 16144 //128 (key) + 16000 (Value) + 8 (position) + 8 (next position)
	maxKeyLength   = 32    //characters
	maxValLength   = 16000 //corresponds to 16KB
)

func handleError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func encode(n *Node) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(n)
	handleError(err)

	return buf.Bytes()
}

func encodeString(val string) []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(val)
	handleError(err)

	return buf.Bytes()
}

func decode(data []byte) *Node {
	var ret Node

	buf := bytes.NewReader(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&ret)
	handleError(err)

	return &ret
}

func decodeString(data []byte) string {
	var ret string

	buf := bytes.NewReader(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&ret)
	handleError(err)

	return ret
}

func truncateKey(key string) string {
	//max allowed key lenght is 32
	if len(key) > maxKeyLength {
		key = key[0:maxKeyLength]
	}

	return key

}
