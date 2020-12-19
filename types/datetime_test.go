package types

import (
	"encoding/json"
	"encoding/xml"

	"testing"
	"time"
)

type TestTime struct {
	NowTime *DateTime `json:"nowTime" xml:"nowTime"`
}

func TestDatetime(t *testing.T) {
	nowtime := time.Now()
	timefmt := "2006-01-02 15:04:05"
	expect := nowtime.Format(timefmt)
	var input = TestTime{
		NowTime: NewDateTime(nowtime),
	}

	resultJSON := &TestTime{}

	jsonBytes, err := json.Marshal(input)
	if err != nil {
		t.Errorf("json.Marshal出错:%w", err)
	}
	err = json.Unmarshal(jsonBytes, resultJSON)
	if err != nil {
		t.Errorf("json.Unmarshal出错:%w", err)
	}
	t.Log("jsonBytes:", string(jsonBytes))

	resultXML := &TestTime{}

	xmlBytes, err := xml.Marshal(input)
	if err != nil {
		t.Errorf("xml.Marshal出错:%w", err)
	}
	err = xml.Unmarshal(xmlBytes, resultXML)
	if err != nil {
		t.Errorf("json.Unmarshal出错:%w", err)
	}

	t.Log("xmlBytes:", string(xmlBytes))
	t.Log("resultJson:", resultJSON.NowTime.String())
	t.Log("resultXml:", resultXML.NowTime.String())

	if expect != resultJSON.NowTime.String() {
		t.Errorf("ResultJSON not equal. got:%s,expect:%s", resultJSON.NowTime.String(), expect)
	}
	if expect != resultXML.NowTime.String() {
		t.Errorf("resultXML not equal. got:%s,expect:%s", resultJSON.NowTime.String(), expect)
	}
}
