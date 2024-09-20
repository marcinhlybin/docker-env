package registry

import (
	"encoding/json"
	"strings"

	"github.com/marcinhlybin/docker-env/logger"
	"github.com/marcinhlybin/docker-env/project"
)

type DockerComposeImages struct {
	Id            string `json:"ID"`
	ContainerName string `json:"ContainerName"`
	Repository    string `json:"Repository"`
	Tag           string `json:"Tag"`
	Size          uint64 `json:"Size"`
}

func (reg *DockerProjectRegistry) fetchProjectImages(p *project.Project) ([]*DockerComposeImages, error) {
	logger.Debug("Fetching images")
	dc := reg.dockerCmd.FetchImagesCommand(p)
	jsonOutput, err := dc.ExecuteWithOutput()
	if err != nil {
		return nil, err
	}
	jsonString := strings.Join(jsonOutput, "")

	return reg.createImagesFromJson(jsonString)
}

func (reg *DockerProjectRegistry) createImagesFromJson(jsonString string) ([]*DockerComposeImages, error) {
	var images []*DockerComposeImages

	err := json.Unmarshal([]byte(jsonString), &images)
	if err != nil {
		return nil, err
	}

	return images, nil
}
