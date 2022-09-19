package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"

	"github.com/Hendryboyz/eth-synchronizer/configs"
	"github.com/Hendryboyz/eth-synchronizer/context"
	"github.com/Hendryboyz/eth-synchronizer/db"
)

func main() {
	configs.Init("local")
	db.Init()
	redisPool := initCachePool(configs.GetConfig())
	pool := work.NewWorkerPool(context.Context{}, 20, "namespace", redisPool)
	pool.PeriodicallyEnqueue("0 * * * * *", "sync_blocks")

	pool.Middleware((*context.Context).Log)

	pool.Job("sync_blocks", (*context.Context).SyncBlocks)

	pool.JobWithOptions("export", work.JobOptions{Priority: 10, MaxFails: 1}, (*context.Context).Export)
	pool.Start()

	interruptWokerPool(pool)
}

func initCachePool(config *viper.Viper) *redis.Pool {
	cacheAddr := fmt.Sprintf("%s:%s", config.GetString("redis.host"), config.GetString("redis.port"))
	fmt.Println(cacheAddr)
	redisPool := &redis.Pool{
		MaxActive: 5,
		MaxIdle:   5,
		Wait:      true,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", cacheAddr)
		},
	}
	return redisPool
}

func interruptWokerPool(pool *work.WorkerPool) {
	signalChain := make(chan os.Signal, 1)
	signal.Notify(signalChain, os.Interrupt, os.Kill)
	<-signalChain

	pool.Stop()
}
