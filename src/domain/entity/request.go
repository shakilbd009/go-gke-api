package entity

type Request struct {
	Project     string `json:"project"`
	Region      string `json:"region"`
	ClusterName string `json:"clusterName"`
}
