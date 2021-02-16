package k8s

type Kcontainer struct {
	ContainerNamespace string `json:"pod_namespace"`
	PodName            string `json:"pod_name"`
	ContainerName      string `json:"container"`
	Port               int32  `json:"port"`
	LivelinessProbe    string `json:"liveliness"`
	ReadynessProbe     string `json:"readyness"`
	Image              string `json:"image"`
}
