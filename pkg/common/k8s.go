//go:build k8s
// +build k8s

package common

import "github.com/eniac/mucache/pkg/nodeIdx"

// in k8s, urls are const, build with -tags k8s,nodex

// const CachedUrl = "memcached" + nodeIdx.NodeIdx + ":11211"
const CachedUrl = "cache" + nodeIdx.NodeIdx + "-redis-master" + ":6379"
const RedisUrl = "redis" + nodeIdx.NodeIdx
const CMUrl = "http://cm" + nodeIdx.NodeIdx
