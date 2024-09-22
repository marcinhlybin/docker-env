package docker

type Image struct {
	Id            string `json:"ID"`
	ContainerName string `json:"ContainerName"`
	Repository    string `json:"Repository"`
	Tag           string `json:"Tag"`
	Size          uint64 `json:"Size"`
}
