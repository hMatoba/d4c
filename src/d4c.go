package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

func main() {
	var (
		excludeFlag string
	)

	flag.StringVar(&excludeFlag, "exclude", "", "exclude flag")
	flag.Parse()

	_excludeFlag := []string{}
	if excludeFlag != "" {
		_excludeFlag = strings.Split(excludeFlag, ",")
	}

	fmt.Println(_excludeFlag)

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
		for _, dockerImage := range tags {
			fmt.Println(dockerImage)
			tag := strings.Split(dockerImage, ":")
			_tag := tag[1]
			isExcluded := contains(_excludeFlag, _tag)
			if !isExcluded {
				out, err := cli.ImagePull(ctx, dockerImage, types.ImagePullOptions{})
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

func contains(arr []string, str string) bool {
	for _, s := range arr {
		if str == s {
			return true
		}
	}
	return false
}
