package context

import (
	"fmt"

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
	fmt.Println("Sync Blocks")
	return nil
}

func (c *Context) Export(job *work.Job) error {
	return nil
}
