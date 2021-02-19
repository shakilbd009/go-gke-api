package services

import (
	"context"

	"github.com/shakilbd009/go-gke-api/src/domain/entity"
	"github.com/shakilbd009/go-gke-api/src/domain/k8s"
	"github.com/shakilbd009/go-gke-api/src/utils/k8auth"
	"github.com/shakilbd009/go-utils-lib/rest_errors"
)

var (
	KCommonServices kCommonInterface = &kCommonService{}
)

type kCommonService struct{}

type kCommonInterface interface {
	GetAllDeployments(ctx context.Context, e *entity.Request) (map[string]map[string][]*k8s.K8sDeployment, rest_errors.RestErr)
}

func (*kCommonService) GetAllDeployments(ctx context.Context, e *entity.Request) (map[string]map[string][]*k8s.K8sDeployment, rest_errors.RestErr) {

	clusters, err := KgkeServices.GetClusters(ctx, e)
	if err != nil {
		return nil, err
	}
	result := make(map[string]map[string][]*k8s.K8sDeployment)
	for _, cluster := range clusters {
		client, err := k8auth.GetGkekubeConfig(ctx, e.Project, e.Region, cluster.ClusterName)
		if err != nil {
			return nil, err
		}
		kn := &k8s.Knamespace{
			Request: e,
		}
		namespaces, err := kn.GetAll(ctx, client, cluster.ClusterName)
		if err != nil {
			return nil, err
		}
		for _, n := range namespaces {
			var k k8s.K8sDeployment
			deployments, err := k.GetAll(ctx, client, cluster.ClusterName, n.Namespace)
			if err != nil {
				return nil, err
			}
			result[cluster.ClusterName] = map[string][]*k8s.K8sDeployment{
				n.Namespace: deployments,
			}
		}
	}
	return result, nil
}
