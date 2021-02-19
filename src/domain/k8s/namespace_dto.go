package k8s

import "github.com/shakilbd009/go-gke-api/src/domain/entity"

type Knamespace struct {
	Request     *entity.Request `json:"details,omitempty"`
	ClusterName string          `json:"clusterName"`
	Namespace   string          `json:"namespace,omitempty"`
}
