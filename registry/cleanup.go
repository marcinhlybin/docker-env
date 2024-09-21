package registry

import (
	"fmt"
	"strings"

	"github.com/marcinhlybin/docker-env/helpers"
	"github.com/marcinhlybin/docker-env/logger"
	"github.com/marcinhlybin/docker-env/project"
)

// Cleanup removes all projects and images
func (reg *DockerProjectRegistry) Cleanup(includeImages bool) error {
	message := "Cleaning up"
	if !includeImages {
		message += " (without images)"
	}
	logger.Info(message)

	// Fetch projects
	includeStopped := true
	projects, err := reg.fetchProjects(includeStopped)
	if err != nil {
		return err
	}

	var errors []string
	var images []*DockerComposeImages

	// Fetch images
	if includeImages {
		images, err = reg.fetchImagesForProjects(projects)
		if err != nil {
			errors = append(errors, err.Error())
		}
	}

	// Order of removal is important, first projects, then images
	// and to get images we need to have projects
	// Remove projects
	err = reg.removeProjects(projects)
	if err != nil {
		errors = append(errors, err.Error())
	}

	// Remove images
	if includeImages {
		err = reg.removeImages(images)
		if err != nil {
			errors = append(errors, err.Error())
		}
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
