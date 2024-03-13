package cm

import (
	"bufio"
	"fmt"
	"github.com/eniac/mucache/pkg/cm"
	"github.com/eniac/mucache/pkg/common"
	"github.com/eniac/mucache/pkg/nodeIdx"
	"github.com/eniac/mucache/pkg/utility"
	"github.com/golang/glog"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	serviceName cm.ServiceName
	port        int
	//memcachedClient *memcache.Client
	cacheClient   *redis.Client
	cmAddresses   map[cm.ServiceName]string
	printTimeFreq int
}

// ReadCacheManagerAddressFile returns a tuple of
// 1. a map from service name to cache manager address
// 2. the service name this cache manager is responsible for
func ReadCacheManagerAddressFile(cmAddsFile string) (map[cm.ServiceName]string, cm.ServiceName) {
	cmAddresses := make(map[cm.ServiceName]string)
	var serviceName cm.ServiceName

	readFile, err := os.Open(cmAddsFile)
	if err != nil {
		panic(err)
	}
	defer readFile.Close()

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	idx := 0
	for fileScanner.Scan() {
		idx++
		line := fileScanner.Text()
		lineNoSuffix := strings.TrimSuffix(line, "\n")
		tokens := strings.Split(lineNoSuffix, " ")
		if !common.ShardEnabled {
			cmAddresses[cm.ServiceName(tokens[0])] = tokens[1]
			if fmt.Sprint(idx) == nodeIdx.NodeIdx {
				serviceName = cm.ServiceName(tokens[0])
			}
		} else {
			shard, err := strconv.Atoi(common.ShardCount)
			if err != nil {
				panic(err)
			}
			//node, err := strconv.Atoi(nodeIdx.NodeIdx)
			//if err != nil {
			//	panic(err)
			//}
			//shardIdx, err := strconv.Atoi(common.ShardIdx)
			//if err != nil {
			//	panic(err)
			//}
			//if idx == (node-1)*shard+shardIdx {
			//	serviceName = cm.ServiceName(tokens[0])
			//}
			for i := 1; i < shard+1; i++ {
				cmAddresses[cm.ServiceName(tokens[0]+fmt.Sprint(i))] = tokens[1] + fmt.Sprint(i)
			}
			if fmt.Sprint(idx) == nodeIdx.NodeIdx {
				serviceName = cm.ServiceName(tokens[0] + common.ShardIdx)
			}
		}
	}
	glog.Infof("Service name: %v", serviceName)
	glog.Infof("Service Map: %v", cmAddresses)
	return cmAddresses, serviceName
}

func InitConfig(httpP int, cmAddsFile string, printTimeFreq int) *Config {
	cmAddresses, serviceName := ReadCacheManagerAddressFile(cmAddsFile)
	// Initialize a memcached client
	//mc := cm.CreateCacheClient(memcachedP)
	c := cm.GetOrCreateCacheClient()
	cfg := Config{
		serviceName:   serviceName,
		port:          httpP,
		cacheClient:   c,
		cmAddresses:   cmAddresses,
		printTimeFreq: printTimeFreq,
	}
	return &cfg
}

func (cfg *Config) Close() {
	cfg.cacheClient.Close()
}

func (cfg *Config) GetCacheManagerAddress(name cm.ServiceName) string {
	if val, ok := cfg.cmAddresses[name]; ok {
		return val
	} else {
		panic(fmt.Sprintf("%s not exisit in %s", name, cfg.cmAddresses))
	}
}

func (cfg *Config) GetNeighbors() []string {
	utility.Assert(common.ShardEnabled)
	var neighbors []string
	for k, v := range cfg.cmAddresses {
		if k == cfg.serviceName {
			continue
		}
		if strings.TrimRight(string(cfg.serviceName), "0123456789") == strings.TrimRight(string(k), "0123456789") {
			neighbors = append(neighbors, v)
		}
	}
	return neighbors
}
