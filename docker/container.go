package docker

import (
	"strings"
)

type Container struct {
	CreatedAt string `json:"CreatedAt"`
	Name      string `json:"Names"`
	State     string `json:"State"`
	Labels    string `json:"Labels"`
}

// Project name without the stack prefix
func (c *Container) ProjectName() string {
	name := c.ComposeProjectName()
	return strings.SplitN(name, "-", 2)[1]
}

// Full project name with stack prefix
func (c *Container) ComposeProjectName() string {
	return c.LabelValue("com.docker.compose.project")
}

func (c *Container) ServiceName() string {
	projectName := c.LabelValue("com.docker.compose.project")
	return strings.TrimPrefix(c.Name, projectName+"-")

}

func (c *Container) LabelValue(labelName string) string {
	label := ""
	for _, l := range strings.Split(c.Labels, ",") {
		if strings.Contains(l, labelName+"=") {
			label = strings.Split(l, "=")[1]
			break
		}
	}

	return label
}
