package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	redis "github.com/go-redis/redis/v8"
	"github.com/kellegous/go/internal"
)

var Debug bool

const (
	nextIDKey = "nextID"
)

// Backend provides access to Redis
type Backend struct {
	client *redis.Client
}

func dbgLogf(format string, v ...interface{}) {
	if Debug {
		log.Printf(format, v...)
	}
}

// New instantiates a new Backend
func New(ctx context.Context, addr, pw string, db int) (*Backend, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
		DB:       db,
	})

	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Print(err)
		return nil, err
	}
	dbgLogf("Redis: %s", pong)
	backend := &Backend{
		client: client,
	}

	return backend, nil
}

// Close the Backend and release associated resources
func (backend *Backend) Close() error {
	return backend.client.Close()
}

// Get retreives a shortcut from the data store.
func (backend *Backend) Get(ctx context.Context, name string) (*internal.Route, error) {
	dbgLogf("[Redis] GET %s\n", name)
	val, err := backend.client.Get(ctx, name).Result()
	if err != nil {
		if err == redis.Nil {
			log.Printf("Route %s does not exist\n", name)
			return nil, internal.ErrRouteNotFound
		}
		log.Print(err)
		return nil, err
	}
	route := &internal.Route{}
	err = json.Unmarshal([]byte(val), &route)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return route, nil
}

// Put stores a new route in the data store
func (backend *Backend) Put(ctx context.Context, key string, rt *internal.Route) error {
	dbgLogf("[Redis] SET %s\n", key)
	val, err := json.Marshal(rt)
	if err != nil {
		log.Print(err)
		return err
	}
	err = backend.client.Set(ctx, key, string(val), 0).Err()
	if err != nil {
		log.Print(err)
	}
	return nil
}

// Del deletes a route from the data store
func (backend *Backend) Del(ctx context.Context, key string) error {
	dbgLogf("[Redis] DEL %s\n", key)
	res, err := backend.client.Del(ctx, key).Result()
	if err != nil {
		log.Print(err)
		return err
	}
	log.Printf("Route %s has been deleted. Result: %d", key, res)
	return nil
}

// List all routes in an iterator, starting with the key prefix of start
func (backend *Backend) List(ctx context.Context, start string) (internal.RouteIterator, error) {
	dbgLogf("[Redis] LIST %s\n", start)
	cmd := backend.client.Scan(ctx, 0, fmt.Sprintf("%s*", start), 0)
	iterator := cmd.Iterator()
	keys, cursor, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		dbgLogf("%s", key)
	}
	dbgLogf("cursor: %d\n", cursor)

	return &RouteIterator{
		it:     iterator,
		ctx:    ctx,
		pos:    int(cursor),
		client: backend.client,
	}, nil
}

// NextID generates the next numeric ID to be used for an auto-named route
func (backend *Backend) NextID(ctx context.Context) (uint64, error) {
	dbgLogf("[Redis] NextID\n")
	result, err := backend.client.Incr(ctx, nextIDKey).Uint64()
	if err != nil {
		log.Print(err)
		return 0, err
	}
	return result, nil
}

// GetAll dumps everything in the db for backup purposes
func (backend *Backend) GetAll(ctx context.Context) (map[string]internal.Route, error) {
	dbgLogf("[Redis] GetAll\n")
	golinks := map[string]internal.Route{}
	cmd := backend.client.Scan(ctx, 0, "*", 0)
	keys, cursor, err := cmd.Result()
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		if key == nextIDKey {
			continue
		}
		dbgLogf("%s", key)
		val, err := backend.client.Get(ctx, key).Result()
		if err != nil {
			if err == redis.Nil {
				log.Printf("Route %s does not exist\n", key)
				return golinks, nil
			}
			log.Print(err)
			return nil, err
		}
		route := &internal.Route{}
		err = json.Unmarshal([]byte(val), &route)
		if err != nil {
			return nil, err
		}
		golinks[key] = *route
	}
	dbgLogf("cursor: %d\n", cursor)
	dbgLogf("[Redis] Getall - RouteMap: %+v\n", golinks)

	return golinks, nil
}
