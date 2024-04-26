package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func Status() error {
	apiClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		log.Fatal("Error creating docker client: ", err)
		return err
	}
	defer apiClient.Close()

	containers, err := apiClient.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		log.Fatal("Error creating docker client: ", err)
		return err
	}

	for _, ctr := range containers {
		fmt.Printf("Container tag: %s\n", ctr.Labels["tag"])
		fmt.Printf("%s %s (status: %s)\n", ctr.ID, ctr.Image, ctr.Status)
	}
	return nil
}