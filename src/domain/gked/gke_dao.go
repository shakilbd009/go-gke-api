package gked

import (
	"context"

	"github.com/shakilbd009/go-gke-api/src/client/gke"
	"github.com/shakilbd009/go-gke-api/src/domain/entity"
	"github.com/shakilbd009/go-utils-lib/rest_errors"
)

func (g *GkeCluster) GetClusters(ctx context.Context) ([]*GkeCluster, rest_errors.RestErr) {

	clusters, err := gke.Kubernetes.GetClusters(ctx, g.Project, g.Region)
	if err != nil {
		return nil, rest_errors.NewBadRequestError(err.Error())
	}
	result := make([]*GkeCluster, 0, len(clusters))

	for _, c := range clusters {
		nodePools := make([]string, 0)

		for _, nodePool := range c.NodePools {
			nodePools = append(nodePools, nodePool.Name)
		}
		result = append(result, &GkeCluster{
			ClusterName: c.Name,
			Request: &entity.Request{
				Project: g.Project,
				Region:  g.Region,
			},
			MasterVersion: c.CurrentMasterVersion,
			NodeCount:     c.CurrentNodeCount,
			NodePools:     nodePools,
		})
	}
	return result, nil
}
