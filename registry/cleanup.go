package registry

import (
	"fmt"
	"strings"

	"github.com/marcinhlybin/docker-env/helpers"
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/marcinhlybin/docker-env/project"
)

func (reg *DockerProjectRegistry) Cleanup() error {
	logger.Info("Cleaning up")
	includeStopped := true
	projects, err := reg.fetchProjects(includeStopped)
	if err != nil {
		return err
	}

	// Fetch images from existing projects
	images, fetchImagesErr := reg.fetchImagesForProjects(projects)

	// Remove projects and images
	removeProjectsErr := reg.removeProjects(projects)
	removeImagesErr := reg.removeImages(images)

	// Check for errors
	var errors []string
	if fetchImagesErr != nil {
		errors = append(errors, fetchImagesErr.Error())
	}
	if removeProjectsErr != nil {
		errors = append(errors, removeProjectsErr.Error())
	}
	if removeImagesErr != nil {
		errors = append(errors, removeImagesErr.Error())
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors during cleanup: %s", strings.Join(errors, ", "))
	}

	return nil
}

func (reg *DockerProjectRegistry) fetchImagesForProjects(projects []*project.Project) ([]*DockerComposeImages, error) {
	var images []*DockerComposeImages
	isErr := false
	for _, p := range projects {
		projectImages, err := reg.fetchProjectImages(p)
		if err != nil {
			isErr = true
			logger.Warning("Could not fetch images for %s: %v", p.String(), err)
			continue
		}
		images = append(images, projectImages...)
	}

	if isErr {
		return nil, fmt.Errorf("one or more images could not be fetched")
	}

	return images, nil
}

func (reg *DockerProjectRegistry) removeProjects(projects []*project.Project) error {
	isErr := false
	for _, p := range projects {
		dc := reg.dockerCmd.RemoveProjectCommand(p)
		err := dc.Execute()
		if err != nil {
			isErr = true
			logger.Warning("Could not remove %s", p.String())
			continue
		}
	}

	if isErr {
		return fmt.Errorf("one or more projects could not be removed")
	}

	return nil
}

func (reg *DockerProjectRegistry) removeImages(images []*DockerComposeImages) error {
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
