package docker

import (
	"encoding/base64"
	"encoding/json"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

type HandlerOutFunc func(io.ReadCloser) error

func PushImage(username, password, tag string, handlerOut HandlerOutFunc) error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	authConfig := types.AuthConfig{
		Username: username,
		Password: password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return err
	}
	authStr := base64.URLEncoding.EncodeToString(encodedJSON)
	out, err := cli.ImagePush(ctx, tag, types.ImagePushOptions{RegistryAuth: authStr})
	if err != nil {
		return err
	}

	defer out.Close()
	return handlerOut(out)
}
