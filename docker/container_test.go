package docker_test

import (
	"testing"

	"github.com/marcinhlybin/docker-env/docker"
)

func TestGetProjectName(t *testing.T) {
	container := docker.Container{
		Labels: "com.docker.compose.project=stackname-projectname",
		Name:   "stackname-projectname",
	}

	expected := "projectname"
	result := container.ProjectName()

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestGetComposeProjectName(t *testing.T) {
	container := docker.Container{
		Labels: "com.docker.compose.project=stackname-projectname",
	}

	expected := "stackname-projectname"
	result := container.ComposeProjectName()

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestGetServiceName(t *testing.T) {
	container := docker.Container{
		Labels: "com.docker.compose.project=stackname-projectname",
		Name:   "stackname-projectname-servicename",
	}

	expected := "servicename"
	result := container.ServiceName()

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}

func TestGetLabelValue(t *testing.T) {
	container := docker.Container{
		Labels: "com.docker.another.value=value,com.docker.compose.project=stackname-projectname,another.label=value",
	}

	expected := "stackname-projectname"
	result := container.LabelValue("com.docker.compose.project")

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}

	expected = "value"
	result = container.LabelValue("another.label")

	if result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
