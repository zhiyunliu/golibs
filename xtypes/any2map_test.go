package xtypes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type TestStruct2 struct {
	ID   int    `json:"id"`
	Name string `json:"-"`
}

type TestStruct3 struct {
	ID   int
	Name string
}

func TestAnyToMap(t *testing.T) {
	obj := TestStruct{
		ID:   1,
		Name: "John",
	}

	expected := map[string]interface{}{
		"id":   1,
		"name": "John",
	}

	result, err := AnyToMap(obj)
	assert.Equal(t, nil, err, "TestAnyToMap not error")
	assert.Equal(t, expected, result, "TestAnyToMap should convert struct to map correctly")
}

func TestAnyToMapWithPointer(t *testing.T) {
	obj := &TestStruct{
		ID:   1,
		Name: "John",
	}

	expected := map[string]interface{}{
		"id":   1,
		"name": "John",
	}

	result, err := AnyToMap(obj)
	assert.Equal(t, nil, err, "TestAnyToMap not error")
	assert.Equal(t, expected, result, "TestAnyToMap should convert struct pointer to map correctly")
}

func TestAnyToMapWithIgnoredField(t *testing.T) {
	obj := TestStruct2{
		ID:   1,
		Name: "John",
	}

	expected := map[string]interface{}{
		"id": 1,
	}

	result, err := AnyToMap(obj)
	assert.Equal(t, nil, err, "TestAnyToMap not error")
	assert.Equal(t, expected, result, "TestAnyToMap should ignore fields with '-' tag")
}

func TestAnyToMapWithEmptyTag(t *testing.T) {
	obj := TestStruct3{
		ID:   1,
		Name: "John",
	}

	expected := map[string]interface{}{
		"ID":   1,
		"Name": "John",
	}

	result, err := AnyToMap(obj)
	assert.Equal(t, nil, err, "TestAnyToMap not error")
	assert.Equal(t, expected, result, "TestAnyToMap should use field name when json tag is empty")
}

func TestAnyToMapWithNestedAnonymousStruct(t *testing.T) {
	type Nested struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"field2"`
	}
	obj := struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Nested
	}{
		ID:   1,
		Name: "John",
		Nested: Nested{
			Field1: "NestedField",
			Field2: 2,
		},
	}

	expected := map[string]interface{}{
		"id":     1,
		"name":   "John",
		"field1": "NestedField",
		"field2": 2,
	}

	result, err := AnyToMap(obj)
	assert.Equal(t, nil, err, "TestAnyToMap not error")
	assert.Equal(t, expected, result, "Test failed, expected: '%v', got:  '%v'", expected, result)
}

func TestAnyToMapWithNestedAnonymousStruct2(t *testing.T) {
	obj := struct {
		ID     int    `json:"id"`
		Name   string `json:"name"`
		Nested struct {
			Field1 string `json:"field1"`
			Field2 int    `json:"field2"`
		} `json:"nested"`
	}{
		ID:   1,
		Name: "John",
		Nested: struct {
			Field1 string `json:"field1"`
			Field2 int    `json:"field2"`
		}{
			Field1: "NestedField",
			Field2: 2,
		},
	}

	expected := map[string]interface{}{
		"id":   1,
		"name": "John",
		"nested": map[string]interface{}{
			"field1": "NestedField",
			"field2": 2,
		},
	}

	result, err := AnyToMap(obj)
	assert.Equal(t, nil, err, "TestAnyToMap not error")
	assert.Equal(t, expected, result, "Test failed, expected: '%v', got:  '%v'", expected, result)
}

type Match struct {
	MatchId   string   `json:"match_id"`
	MapName   string   `json:"map_name"`
	BeginTime string   `json:"begin_time"`
	EndTime   string   `json:"won_time"`
	WonDate   int      `json:"won_date"`
	TeamCount int      `json:"team_count"`
	GameType  string   `json:"game_type"`
	TeamSize  int      `json:"team_size"`
	Winners   []Player `json:"winners"`
}

type Player struct {
	Id    string `json:"playerId"`
	Name  string `json:"name"`
	Kills int    `json:"kills"`
}
type UserItem struct {
	Match
	UserId     int    `json:"user_id"`
	LoginId    string `json:"login_id"`
	PubgName   string `json:"pubg_nickname"`
	PubgUserId string `json:"pubg_id"`
}

func TestAnyToMapWithNestedAndAnonymousStruct(t *testing.T) {
	obj := UserItem{
		Match: Match{
			MatchId:   "1",
			MapName:   "Map1",
			BeginTime: "10:00",
			EndTime:   "11:00",
			WonDate:   20220202,
			TeamCount: 2,
			GameType:  "Type1",
			TeamSize:  5,
			Winners: []Player{
				{
					Id:    "Player1",
					Name:  "John",
					Kills: 10,
				},
				{
					Id:    "Player2",
					Name:  "Doe",
					Kills: 15,
				},
			},
		},
		UserId:     1,
		LoginId:    "User1",
		PubgName:   "User1",
		PubgUserId: "PubgUser1",
	}

	expected := map[string]interface{}{
		"match_id":   "1",
		"map_name":   "Map1",
		"begin_time": "10:00",
		"won_time":   "11:00",
		"won_date":   20220202,
		"team_count": 2,
		"game_type":  "Type1",
		"team_size":  5,
		"winners": []interface{}{
			map[string]interface{}{
				"playerId": "Player1",
				"name":     "John",
				"kills":    10,
			},
			map[string]interface{}{
				"playerId": "Player2",
				"name":     "Doe",
				"kills":    15,
			},
		},
		"user_id":       1,
		"login_id":      "User1",
		"pubg_nickname": "User1",
		"pubg_id":       "PubgUser1",
	}

	result, err := AnyToMap(obj)
	assert.Equal(t, nil, err, "TestAnyToMap not error")
	assert.Equal(t, expected, result, "Test failed, expected: '%v', got:  '%v'", expected, result)
}
