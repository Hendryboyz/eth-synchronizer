package context

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/Hendryboyz/eth-synchronizer/configs"
	"github.com/Hendryboyz/eth-synchronizer/db"
	"github.com/Hendryboyz/eth-synchronizer/models"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gocraft/work"
)

type Context struct {
	blockId string
}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Printf("Starting job[%s]\n", job.Name)
	return next()
}

func (c *Context) SyncBlocks(job *work.Job) error {
	config := configs.GetConfig()
	client, err := ethclient.Dial(config.GetString("scans.nodeEndpoint"))
	if err != nil {
		panic(err)
	}
	blockNum, err := client.BlockNumber(context.Background())
	if err != nil {
		return err
	}
	syncTime := time.Now()
	fmt.Printf("Start sync block at %s\n", syncTime)
	n := config.GetInt64("scans.n")
	wg := new(sync.WaitGroup)
	for i := int64(blockNum) - n; i <= int64(blockNum); i++ {
		wg.Add(1)
		go doSync(client, i, wg)
	}
	wg.Wait()
	fmt.Printf("Finish sync at %s\n", syncTime)
	return nil
}

func doSync(client *ethclient.Client, blockNum int64, wg *sync.WaitGroup) {
	defer wg.Done()
	block, err := client.BlockByNumber(
		context.Background(),
		big.NewInt(blockNum),
	)
	if err != nil {
		panic(err)
	}

	blockEntity := createBlockEntity(block)
	blockRepo := db.BlocksRepository{DB: db.GetDBInstance()}
	blockRepo.CreateBlock(blockEntity)
}

func createBlockEntity(block *types.Block) models.Blocks {
	return models.Blocks{
		Number:     block.Number().Int64(),
		Hash:       block.Hash().String(),
		ParentHash: block.ParentHash().String(),
		Time:       time.Unix(int64(block.Time()), 0),
	}
}

func (c *Context) Export(job *work.Job) error {
	return nil
}
