package main

import (
	"context"
	"flag"
	"strings"
	"time"

	"caching-proxies-cache/config"
	"caching-proxies-cache/support/cache"
	"caching-proxies-cache/support/connection"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func main() {
	flag.Parse()
	ctx := context.Background()

	logrus.SetFormatter(&logrus.JSONFormatter{PrettyPrint: true})

	cch := cache.NewCache(ctx, *config.FlagCacheTtl)

	// Echo instance
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Middleware
	e.Use(middleware.Recover())

	// Routes
	e.GET("/get", requestProcessor(cch))

	// Start server
	go func() {
		e.Logger.Fatal(e.Start("0.0.0.0:1324"))
	}()

	ns, _, _ := connection.Establish(*config.FlagNatsContext, *config.FlagServer, *config.FlagCreds)

	mergedInput := make(chan *nats.Msg, 1024)

	for _, shard := range strings.Split(*config.FlagListenToShards, ",") {
		if shard != "" {
			go listenForShard(ctx, ns, shard, mergedInput)
		}
	}

	for {
		select {
		case msg := <-mergedInput:
			processMessage(msg, cch)
		case <-ctx.Done():
			// We don't need to drain the mergedInput
			return
		}
	}
}

func requestProcessor(cch *cache.Cache) func(c echo.Context) error {
	return func(c echo.Context) error {
		previousHashID := c.QueryParam("previous_hash_id")
		shardID := c.QueryParam("shard_id")
		ans := cch.Get(previousHashID, shardID)
		if ans == nil {
			logrus.WithFields(map[string]interface{}{
				"previous_hash_id": previousHashID,
				"shard_id":         shardID,
				"result":           "not_found",
			}).Infof("Message with previous hash id %s and shard id %s not found", previousHashID, shardID)
			return c.String(404, "Not found")
		}

		wasInCacheFor := time.Since(ans.WasCachedAt)
		logrus.WithFields(map[string]interface{}{
			"previous_hash_id":    previousHashID,
			"shard_id":            shardID,
			"result":              "found",
			"was_in_cache_for_ns": wasInCacheFor.Nanoseconds(),
		}).
			Infof(
				"Message with prevous hash id %s and shard id %s found and was in cache for %s",
				previousHashID, shardID, wasInCacheFor.String(),
			)
		return c.Blob(200, "application/octet-stream", ans.Blob)
	}
}

func processMessage(msg *nats.Msg, cch *cache.Cache) {
	blockHash := msg.Header.Get("X-Block-Hash")
	previousHashID := msg.Header.Get("X-Previous-Hash-Id")
	shardID := msg.Header.Get("X-Shard-Id")
	blob := msg.Data

	cch.SetIfBlobIsBigger(previousHashID, shardID, blob)
	logrus.WithFields(map[string]interface{}{
		"previous_hash_id": previousHashID,
		"shard_id":         shardID,
		"block_hash":       blockHash,
		"msg_id":           msg.Header.Get(nats.MsgIdHdr),
	}).Info("Processed message on nats subject: ", msg.Subject)
}

func listenForShard(ctx context.Context, ns *nats.Conn, shard string, output chan *nats.Msg) {
	shardSubject := *config.FlagShardPrefix + ":" + shard
	subscription, err := ns.ChanSubscribe(shardSubject, output)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Info("Listening for shard subject: ", shardSubject)

	<-ctx.Done()
	_ = subscription.Drain()
}
