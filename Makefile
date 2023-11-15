build:
	go build main.go

example-run:
	./main -shards-to-listen 1,2,3 -shard-prefix shards -nats lol -cache-ttl 1m

example-curl:
	curl -v -XGET localhost:1324/get\?previous_hash_id=kek\&shard_id=3