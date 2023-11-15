package config

import (
	"flag"
	"time"
)

var (
	FlagNatsContext    = flag.String("nats", "", "NATS context to use")
	FlagServer         = flag.String("server", "", "NATS server to connect to")
	FlagCreds          = flag.String("creds", "", "NATS credentials file")
	FlagShardPrefix    = flag.String("shard-prefix", "", "Prefix for shard subjects")
	FlagListenToShards = flag.String("shards-to-listen", "", "Comma-separated list of shards to listen to")
	FlagCacheTtl       = flag.Duration("cache-ttl", 1*time.Minute, "Cache TTL")
)
