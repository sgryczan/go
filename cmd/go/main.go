package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/kellegous/go/backend"
	"github.com/kellegous/go/backend/firestore"
	"github.com/kellegous/go/backend/leveldb"
	"github.com/kellegous/go/backend/redis"
	"github.com/kellegous/go/web"
)

func main() {
	pflag.String("addr", ":8067", "default bind address")
	pflag.Bool("admin", false, "allow admin-level requests")
	pflag.String("version", "", "version string")
	pflag.String("backend", "leveldb", "backing store to use. 'leveldb' and 'firestore' currently supported.")
	pflag.String("data", "data", "The location of the leveldb data directory")
	pflag.String("project", "", "The GCP project to use for the firestore backend. Will attempt to use application default creds if not defined.")
	pflag.String("redis-addr", "", "Address of the redis DB to use")
	pflag.String("redis-pw", "", "Password to the redis DB")
	pflag.String("redis-db", "", "Redis DB to use.")
	pflag.Bool("redis-debug", false, "Enable redis debug logging")
	pflag.String("host", "", "The host field to use when gnerating the source URL of a link. Defaults to the Host header of the generate request")
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Panic(err)
	}

	// allow env vars to set pflags
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	var backend backend.Backend

	switch viper.GetString("backend") {
	case "leveldb":
		var err error
		backend, err = leveldb.New(viper.GetString("data"))
		if err != nil {
			log.Panic(err)
		}
	case "firestore":
		var err error

		backend, err = firestore.New(context.Background(), viper.GetString("project"))
		if err != nil {
			log.Panic(err)
		}
	case "redis":
		var err error
		redis.Debug = viper.GetBool("redis-debug")
		backend, err = redis.New(context.Background(), viper.GetString("redis-addr"), viper.GetString("redis-pw"), viper.GetInt("redis-db"))
		if err != nil {
			log.Panic(err)
		}
	default:
		log.Panic(fmt.Sprintf("unknown backend %s", viper.GetString("backend")))
	}

	defer backend.Close()

	log.Panic(web.ListenAndServe(backend))
}
