package main

import (
	"encoding/json"
)

//KV  IntArray
type KV struct {
	K string
	V string
}

//IntArray IntArray
type IntArray []KV

//Marshal Marshal
func (i IntArray) Marshal() (string, error) {
	bytes, err := json.Marshal(i)
	return string(bytes), err
}

//Unmarshal Unmarshal
func (i *IntArray) Unmarshal(val string) error {
	err := json.Unmarshal([]byte(val), i)
	return err
}
