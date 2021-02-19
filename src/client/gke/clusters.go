package gke

import (
	"context"
	"fmt"

	ctn "cloud.google.com/go/container/apiv1"
	ctnpb "google.golang.org/genproto/googleapis/container/v1"
)

var (
	Kubernetes gkeInterface = &gke{}
)

type gkeInterface interface {
	GetClusters(ctx context.Context, project, region string) ([]*ctnpb.Cluster, error)
	GetClient(ctx context.Context) (*ctn.ClusterManagerClient, error)
}

type gke struct{}

func (g *gke) GetClusters(ctx context.Context, project, region string) ([]*ctnpb.Cluster, error) {

	client, err := g.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := client.ListClusters(ctx, &ctnpb.ListClustersRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", project, region)})
	if err != nil {
		return nil, err
	}
	return resp.Clusters, nil
}

func (*gke) GetClient(ctx context.Context) (*ctn.ClusterManagerClient, error) {

	client, err := ctn.NewClusterManagerClient(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}
