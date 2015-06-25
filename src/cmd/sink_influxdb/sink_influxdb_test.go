package main

import "testing"
import "reflect"
import "time"
import "fmt"

func TestTagSplit(t *testing.T) {
	var tagString string = "KeyName:cpu:region:us-west-2b:instance_id:i-24e2rgw"
	var tagsResult = map[string]string{
		"KeyName":     "cpu",
		"region":      "us-west-2b",
		"instance_id": "i-24e2rgw",
	}

	parsedTags := TagSplit(tagString)

	if len(parsedTags) != len(tagsResult) {
		t.Error("Length of slices did not match.")
	}

	for i := range parsedTags {
		if parsedTags[i] != tagsResult[i] {
			t.Error("Error tags do not match.")
		}
	}
}

func TestCreateMessage(t *testing.T) {
	var msgString string = "KeyName:cpu:region:us-west-2b:" +
		"instance_id:i-24e2rgw|10.31|2015-06-24T05:48:56.865215158Z"
	var tags = map[string]string{
		"region":      "us-west-2b",
		"instance_id": "i-24e2rgw",
	}

	const RFC3339 = "2006-01-02T15:04:05Z07:00"
	timestamp, _ := time.Parse(RFC3339, "2015-06-24T05:48:56.865215158Z")

	var msgResult = Message{
		Key:   "cpu",
		Value: 10.3100004196167,
		Tags:  tags,
		Time:  timestamp,
	}

	m := CreateMessage(msgString)

	eq := reflect.DeepEqual(m, msgResult)
	if eq == false {
		t.Error("Message maps do not match.\n")
		fmt.Printf("Test(m): %+v\n\nResult(msgResult)): %+v\n", m, msgResult)
	}
}
