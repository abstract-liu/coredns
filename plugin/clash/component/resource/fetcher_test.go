package resource

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type Demo struct {
	UserId    int    `json:"userId"`
	Id        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var idx = 0

func ParseDemo(data []byte) (Demo, error) {
	var demo Demo
	err := json.Unmarshal(data, &demo)
	return demo, err
}

func OnUpdate(demo Demo) {
	fmt.Sprintf("demo-%d: %+v", idx, demo)
	idx++
}

func TestFetcher(t *testing.T) {
	fetcher := NewFetcher("test", "https://jsonplaceholder.typicode.com/todos/1", 3*time.Second, ParseDemo, OnUpdate)
	initialDemo, err := fetcher.Initial()
	assert.Nil(t, err)
	assert.NotNil(t, initialDemo)

}
