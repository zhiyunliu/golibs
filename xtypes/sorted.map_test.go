package xtypes

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewSortedMap(t *testing.T) {

	sortMap := NewSortedMap[int, string](func(a, b int) bool {
		return a < b
	})

	result := map[int]string{
		1: "a",
		2: "b",
		5: "ee",
		4: "d",
		3: "c",
		0: "0",
	}

	sortMap.Put(1, "a")
	sortMap.Put(2, "b")
	sortMap.Put(5, "e")
	sortMap.Put(4, "d")
	sortMap.Put(5, "ee")
	sortMap.Put(3, "c")
	sortMap.Put(0, "0")

	sortMap.Each(func(i int, s string) {
		if !strings.EqualFold(result[i], s) {
			t.Errorf("排序不正确:%d,%s", i, s)
		}
	})

	get, ok := sortMap.Get(1)
	if !ok {
		t.Errorf("Get error not exists")
	}
	if get != "a" {
		t.Errorf("Get error not equal")
	}

	sortMap.Remove(1)
	get1, ok := sortMap.Get(1)
	if ok {
		t.Errorf("Get2 error exists")
	}
	if get1 == "a" {
		t.Errorf("Get2 error equal")
	}

	ok = sortMap.Contains(1)
	if ok {
		t.Errorf("Contains error exists")
	}

	keys := sortMap.Keys()
	if len(keys) != sortMap.Len() {
		t.Errorf("Keys count not equal")
	}

	exactualKeys := []int{0, 2, 3, 4, 5}
	for i := range exactualKeys {
		if keys[i] != exactualKeys[i] {
			t.Errorf("keys not equal ")
		}
	}

	values := sortMap.Values()
	if len(values) != sortMap.Len() {
		t.Errorf("values count not equal")
	}

	exactualVal := []string{"0", "b", "c", "d", "ee"}
	for i := range exactualVal {
		if values[i] != exactualVal[i] {
			t.Errorf("values not equal ")
		}
	}

}
func TestMarshalJSON(t *testing.T) {

	sortMap := NewSortedMap[int, string](func(a, b int) bool {
		return a < b
	})

	result := map[int]string{
		1: "a",
	}

	sortMap.Put(1, "a")

	smbytes, _ := json.Marshal(sortMap)
	nmbytes, _ := json.Marshal(result)

	if !bytes.EqualFold(smbytes, nmbytes) {
		t.Errorf("Marshal bytes not equal")
	}

	bbytes, _ := sortMap.MarshalBinary()

	if !bytes.EqualFold(bbytes, nmbytes) {
		t.Errorf("MarshalBinary bytes not equal")
	}
}
