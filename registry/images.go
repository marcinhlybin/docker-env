package registry

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/marcinhlybin/docker-env/docker"
	"github.com/marcinhlybin/docker-env/helpers"
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/marcinhlybin/docker-env/project"
)

func (reg *DockerProjectRegistry) fetchImages(p *project.Project) ([]*docker.Image, error) {
	logger.Debug("Fetching images")
	dc := reg.dockerCmd.FetchImagesCommand(p)
	jsonOutput, err := dc.ExecuteWithOutput()
	if err != nil {
		return nil, err
	}
	jsonString := strings.Join(jsonOutput, "")

	return reg.createImagesFromJson(jsonString)
}

func (reg *DockerProjectRegistry) createImagesFromJson(jsonString string) ([]*docker.Image, error) {
	var images []*docker.Image

	err := json.Unmarshal([]byte(jsonString), &images)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (reg *DockerProjectRegistry) removeImages(images []*docker.Image) error {
	var removedIds []string
	isErr := false
	for _, img := range images {
		if helpers.Contains(removedIds, img.Id) {
			continue
		}
		removedIds = append(removedIds, img.Id)
		dc := reg.dockerCmd.RemoveImageCommand(img.Id)
		err := dc.Execute()
		if err != nil {
			isErr = true
			logger.Warning("Could not remove image %s", img.Id)
		}
	}

	if isErr {
		return fmt.Errorf("one or more images could not be removed")
	}

	return nil
}
