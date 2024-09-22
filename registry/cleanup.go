package registry

import (
	"fmt"
	"strings"

	"github.com/marcinhlybin/docker-env/docker"
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
	var images []*docker.Image

	// Fetch images
	if includeImages {
		images, err = reg.fetchImagesForMultipleProjects(projects)
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

func (reg *DockerProjectRegistry) fetchImagesForMultipleProjects(projects []*project.Project) ([]*docker.Image, error) {
	var images []*docker.Image
	isErr := false
	for _, p := range projects {
		projectImages, err := reg.fetchImages(p)
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
