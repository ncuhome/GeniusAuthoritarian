package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/ncuhome/GeniusAuthoritarian/internal/rpc/sshDev/rpcModel"
)

func SubScribeSshDev() *redis.PubSub {
	return Client.Subscribe(context.Background(), keySshDevSub.String())
}

func PublishSshDev(messages []rpcModel.SshAccountMsg) error {
	rpcMsgBytes, err := json.Marshal(messages)
	if err != nil {
		return err
	}
	return Client.Publish(context.Background(), keySshDevSub.String(), rpcMsgBytes).Err()
}
