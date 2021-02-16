package k8s

import "github.com/shakilbd009/go-gke-api/src/domain/entity"

type Kpod struct {
	Request      *entity.Request `json:"details,omitempty"`
	PodNamespace string          `json:"pod_namespace"`
	PodName      string          `json:"pod_name"`
	Status       string          `json:"status"`
	Restarts     int32           `json:"restarts"`
	PodIP        string          `json:"pod_ip"`
}
