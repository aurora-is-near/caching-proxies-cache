# Project Name: Caching Proxies Cache Service

## Overview
The Caching Proxies Cache Service is a Go-based solution aimed to give access to cached chunks submitted by Caching Proxies Terminal. This service is a key component of the RPC Speedup project and is designed to handle HTTP requests from RPC nodes efficiently. It does not require authentication, simplifying the interaction with RPC nodes.

## Installation

### Prerequisites
- Go 1.20 or higher
- Access to a NATS server

### Steps
1. Clone the repository:
   ```bash
   git clone [repository-url]
   ```
2. Navigate to the project directory:
   ```bash
   cd rpc-speedup-cache
   ```
3. Build the project using the Makefile:
   ```bash
   make build
   ```

## Configuration Options
The service can be configured using command-line flags:

- **`-nats`** (string): Specifies the NATS context.
  ```bash
  ./main -nats=[context]
  ```

- **`-server`** (string): Sets the NATS server address.
  ```bash
  ./main -server=[server_address]
  ```

- **`-creds`** (string): Path to the NATS credentials file.
  ```bash
  ./main -creds=[path_to_credentials]
  ```

- **`-shard-prefix`** (string): Prefix for shard subjects.
  ```bash
  ./main -shard-prefix=[prefix]
  ```

- **`-shards-to-listen`** (string): List of shards to listen to, separated by commas.
  ```bash
  ./main -shards-to-listen=1,2,3
  ```

- **`-cache-ttl`** (duration): Time-to-live for the cache.
  ```bash
  ./main -cache-ttl=1m
  ```

## Usage

### Running the Service
- Start the service using the configured flags as needed:
  ```bash
  ./main -shards-to-listen=1,2,3 -shard-prefix=shards -nats=context -cache-ttl=1m
  ```

### Interacting with the Service
- Retrieve a cached chunk:
  ```bash
  curl -v -XGET localhost:1324/get?previous_hash_id=[id]&shard_id=[id]
  ```

## Contributing
Contributions to the RPC Speedup Cache Service are welcome. Please adhere to the standard Git workflow: fork the repository, make your changes, and submit a pull request for review.