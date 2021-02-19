package services

import (
	"context"

	"github.com/shakilbd009/go-gke-api/src/domain/k8s"
	"github.com/shakilbd009/go-gke-api/src/utils/k8auth"
	"github.com/shakilbd009/go-utils-lib/rest_errors"
)

var KPodServices kpodInterface = &kpodService{}

type kpodInterface interface {
	GetPods(ctx context.Context, projectID, region, clusterName, namespace string) ([]k8s.Kpod, rest_errors.RestErr)
	GetPodLogs(ctx context.Context, projectID, region, clusterName, namespace string, podName string) ([]string, rest_errors.RestErr)
}
type kpodService struct{}

func (*kpodService) GetPods(ctx context.Context, projectID, region, clusterName, namespace string) ([]k8s.Kpod, rest_errors.RestErr) {

	// if _, err := AuthService.GetToken(token); err != nil {
	// 	return nil, err
	// }
	client, rest_err := k8auth.GetGkekubeConfig(ctx, projectID, region, clusterName)
	if rest_err != nil {
		return nil, rest_err
	}
	var pod k8s.Kpod
	pods, err := pod.GetPods(ctx, client, namespace)
	if err != nil {
		return nil, err
	}
	return pods, nil
}

func (*kpodService) GetPodLogs(ctx context.Context, projectID, region, clusterName, namespace string, podName string) ([]string, rest_errors.RestErr) {
	client, rest_err := k8auth.GetGkekubeConfig(ctx, projectID, region, clusterName)
	if rest_err != nil {
		return nil, rest_err
	}
	var pod k8s.Kpod
	logs, err := pod.GetPodLogs(ctx, client, namespace, podName)
	if err != nil {
		return nil, rest_errors.NewBadRequestError(err.Error())
	}
	return logs, nil
}
