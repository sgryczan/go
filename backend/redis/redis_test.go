package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	redismock "github.com/elliotchance/redismock/v8"
	redis "github.com/go-redis/redis/v8"
	"github.com/kellegous/go/internal"
	"github.com/stretchr/testify/assert"
)

var (
	client *redis.Client
)

var (
	key   = "key"
	val   = "val"
	Addr  string
	case1 = `{"url":"http://czan.io","time":"2020-09-29T16:23:56.71891-06:00"}`
)

func TestMain(m *testing.M) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	client = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	// Addr is the address of the mock redis instance
	Addr := mr.Addr()
	_ = fmt.Sprintf("%s", Addr)

	code := m.Run()
	os.Exit(code)
}

func TestSet(t *testing.T) {
	exp := time.Duration(0)

	mock := redismock.NewNiceMock(client)
	mock.On("Set", key, val, exp).Return(redis.NewStatusResult("", nil))

	backend, err := New(context.Background(), Addr, "", 0)
	route := &internal.Route{}
	err = json.Unmarshal([]byte(case1), &route)
	if err != nil {
		t.Fatalf("Failed with err: %s", err)
	}
	err = backend.Put(context.Background(), key, route)
	if err != nil {
		t.Fatalf("Failed with err: %s", err)
	}
	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	mock := redismock.NewNiceMock(client)
	mock.On("Get", key).Return(redis.NewStringResult(val, nil))

	backend, err := New(context.Background(), Addr, "", 0)
	route, err := backend.Get(context.Background(), key)
	js, err := json.Marshal(route)
	if err != nil {
		t.Fatalf("Failed to get value: %s", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, case1, string(js))
}
