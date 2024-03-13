//go:build shard
// +build shard

package common

import "os"

const ShardEnabled = true

var ShardCount = os.Getenv("SHARD_COUNT")
var ShardIdx = os.Getenv("SHARD_IDX")
