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
	key         = "key"
	val         = "val"
	Addr        string
	Mock        *redismock.ClientMock
	MockBackend *Backend
	case1       = `{"url":"http://czan.io","time":"2020-09-29T16:23:56.71891-06:00"}`
)

func TestMain(m *testing.M) {
	mr, err := miniredis.Run()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	Mock = redismock.NewNiceMock(client)

	client = redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	// Addr is the address of the mock redis instance
	Addr = mr.Addr()
	_ = fmt.Sprintf("%s", Addr)

	code := m.Run()
	os.Exit(code)
}

func TestNew(t *testing.T) {
	var err error

	_ = redismock.NewNiceMock(client)

	MockBackend, err = New(context.Background(), Addr, "", 0)
	if err != nil {
		t.Fatalf("Failed with err: %s", err)
	}
	assert.NoError(t, err)
}

func TestSet(t *testing.T) {
	var err error
	exp := time.Duration(0)

	Mock.On("Set", key, val, exp).Return(redis.NewStatusResult("", nil))

	route := &internal.Route{}
	err = json.Unmarshal([]byte(case1), &route)
	if err != nil {
		t.Fatalf("Failed with err: %s", err)
	}
	err = MockBackend.Put(context.Background(), key, route)
	if err != nil {
		t.Fatalf("Failed with err: %s", err)
	}
	assert.NoError(t, err)
}

func TestGet(t *testing.T) {
	var err error
	Mock.On("Get", key).Return(redis.NewStringResult(val, nil))

	route, err := MockBackend.Get(context.Background(), key)
	js, err := json.Marshal(route)
	if err != nil {
		t.Fatalf("Failed to get value: %s", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, case1, string(js))
}

func TestDel(t *testing.T) {
	var err error
	Mock.On("Del", key).Return(redis.NewStringResult(val, nil))

	err = MockBackend.Del(context.Background(), key)
	if err != nil {
		t.Fatalf("Failed to delete key: %s", err)
	}

	assert.NoError(t, err)
}

// TestNextID ensures that nextID will be created and set to 1
// if it doesn't already exist
func TestNextID(t *testing.T) {
	var err error

	next, err := MockBackend.NextID(context.Background())
	if err != nil {
		t.Fatalf("Failed to increment nextID: %s", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(1), next)
}

// TestNextID2 makes sure nextID can increment further
// in the case it already exists (in this case, from above test)
func TestNextID2(t *testing.T) {
	var err error

	next, err := MockBackend.NextID(context.Background())
	if err != nil {
		t.Fatalf("Failed to increment nextID: %s", err)
	}

	assert.NoError(t, err)
	assert.Equal(t, uint64(2), next)
}
