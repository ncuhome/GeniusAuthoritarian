package tokenStore

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"sync/atomic"
	"time"
)

type Node struct {
	// this is a unique id for each node process, it
	// will be reallocated every day.
	// ID must be smaller than 100.
	ID uint64

	keyNodeIDPrefix string

	// use for refresh fields.
	Lock       *sync.Mutex
	IDTimeMark uint64
	TokenID    *atomic.Uint64
}

func (node *Node) keyNodeID(timeMark uint64) string {
	return fmt.Sprintf("%s-%d", node.keyNodeIDPrefix, timeMark)
}

func (node *Node) currentTimeMark() uint64 {
	return uint64(time.Now().YearDay())
}

func (node *Node) WithClient(client *redis.Client) NodeWithClient {
	return NodeWithClient{
		Node:   node,
		client: client,
	}
}

type NodeWithClient struct {
	*Node
	client *redis.Client
}

func (node NodeWithClient) GenID(ctx context.Context) (uint64, error) {
	currentTimeMark := node.currentTimeMark()
	if node.IDTimeMark != currentTimeMark {
		node.Lock.Lock()
		if node.IDTimeMark == currentTimeMark {
			node.Lock.Unlock()
			return node.GenID(ctx)
		}
		defer node.Lock.Unlock()
		newNodeID, err := node.client.Incr(ctx, node.keyNodeID(currentTimeMark)).Uint64()
		if err != nil {
			return 0, err
		}
		node.ID = newNodeID % 100
		node.TokenID.Store(0)
		node.IDTimeMark = currentTimeMark
	}
	tokenID := node.TokenID.Add(1)
	tokenID = (tokenID << 5) + (node.ID << 3) + node.IDTimeMark
	return tokenID, nil
}
