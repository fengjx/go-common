package event

import (
	"testing"
	"time"
)

func TestOnDuplicateHandler(t *testing.T) {
	testEvent := Event[string]("foo")

	On(testEvent, func(data string) {
		t.Log("bar", data)
	})
	On(testEvent, func(data string) {
		t.Log("abc", data)
	})

	Emit(testEvent, "testData")

	time.Sleep(time.Second)
}
