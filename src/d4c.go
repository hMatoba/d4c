package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func main() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		tags := image.RepoTags
		for _, tag := range tags {
			fmt.Println(tag)
			isLatest := strings.HasSuffix(tag, ":latest")
			if isLatest {
				out, err := cli.ImagePull(ctx, tag, types.ImagePullOptions{})
				if err != nil {
					fmt.Println(err)
					continue
				}
				defer out.Close()
				io.Copy(os.Stdout, out)
			} else {
				fmt.Printf("Passed pulling image: %s\n", tag)
			}
		}
	}

	os.Exit(0)
}
