package context

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/Hendryboyz/eth-synchronizer/configs"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gocraft/work"
)

type Context struct {
	blockId string
}

func (c *Context) Log(job *work.Job, next work.NextMiddlewareFunc) error {
	fmt.Println("Starting job:", job.Name)
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
	n := config.GetInt64("scans.n")
	for i := int64(blockNum) - n; i <= int64(blockNum); i++ {
		doSync(client, i)
	}
	return nil
}

func doSync(client *ethclient.Client, blockNum int64) {
	fmt.Println("Current block number is :", blockNum)
	block, err := client.BlockByNumber(
		context.Background(),
		big.NewInt(blockNum),
	)
	if err != nil {
		panic(err)
	}

	blockInfo := make(map[string]string)
	blockInfo["num"] = block.Number().String()
	blockInfo["hash"] = block.Hash().String()
	blockInfo["time"] = strconv.FormatUint(block.Time(), 10)
	blockInfo["parentHash"] = block.ParentHash().String()
	fmt.Println(blockInfo)
}

func (c *Context) Export(job *work.Job) error {
	return nil
}
