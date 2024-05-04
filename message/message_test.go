package message

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageID(t *testing.T) {
	id := NewMessageIDFromString("test case")
	msg := Message([]MessageSegment{Reply(id)})
	data, err := json.Marshal(&msg)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
	if !assert.Equal(t, `[{"type":"reply","data":{"id":"test case"}}]`, string(data)) {
		t.Fail()
	}
	id = NewMessageIDFromInteger(-1008611)
	paras := map[string]interface{}{}
	paras["message_id"] = id
	data, err = json.Marshal(&paras)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
	if !assert.Equal(t, `{"message_id":-1008611}`, string(data)) {
		t.Fail()
	}
	id = NewMessageIDFromString("test case")
	paras = map[string]interface{}{}
	paras["message_id"] = id
	data, err = json.Marshal(&paras)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
	if !assert.Equal(t, `{"message_id":"test case"}`, string(data)) {
		t.Fail()
	}
}
