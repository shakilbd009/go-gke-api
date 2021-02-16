package services

import (
	"context"

	"github.com/shakilbd009/go-gke-api/src/domain/k8s"
	"github.com/shakilbd009/go-gke-api/src/utils/k8auth"
	"github.com/shakilbd009/go-utils-lib/rest_errors"
)

var KContainerServices kcontainerInterface = &kcontainerService{}

type kcontainerInterface interface {
	GetContainerLogs(ctx context.Context, projectID, region, clusterName, namespace, podName, container string) ([]string, rest_errors.RestErr)
	GetContainers(ctx context.Context, projectID, region, clusterName, namespace, podName string) ([]k8s.Kcontainer, rest_errors.RestErr)
}
type kcontainerService struct{}

func (*kcontainerService) GetContainers(ctx context.Context, projectID, region, clusterName, namespace, podName string) ([]k8s.Kcontainer, rest_errors.RestErr) {
	var cont k8s.Kcontainer
	client, rest_err := k8auth.GetGkekubeConfig(ctx, projectID, region, clusterName)
	if rest_err != nil {
		return nil, rest_err
	}
	containers, err := cont.GetContainers(ctx, client, namespace, podName)
	if err != nil {
		return nil, err
	}
	return containers, nil
}

func (*kcontainerService) GetContainerLogs(ctx context.Context, projectID, region, clusterName, namespace, podName, container string) ([]string, rest_errors.RestErr) {

	client, rest_err := k8auth.GetGkekubeConfig(ctx, projectID, region, clusterName)
	if rest_err != nil {
		return nil, rest_err
	}
	var cont k8s.Kcontainer
	logs, err := cont.GetContainerLogs(ctx, client, namespace, podName, container)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
