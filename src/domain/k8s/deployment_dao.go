package k8s

import (
	"context"
	"fmt"

	"github.com/shakilbd009/go-gke-api/src/client/kubernetes/appsv1"
	"github.com/shakilbd009/go-utils-lib/rest_errors"
	"k8s.io/client-go/kubernetes"
)

var (
	statusMsg = "deployment_%s_request_is_successfull"
	creation  = "creation"
	deletion  = "deletion"
)

func (k *K8sDeployment) Get(ctx context.Context, client *kubernetes.Clientset, clusterName, namespace, deploymentName string) (*K8sDeployment, rest_errors.RestErr) {

	resp, err := appsv1.Deployment.Get(ctx, client, namespace, deploymentName)
	if err != nil {
		return nil, rest_errors.NewBadRequestError(err.Error())
	}
	return &K8sDeployment{
		ClusterName:    clusterName,
		Namespace:      namespace,
		DeploymentName: resp.Name,
		Replicas:       *resp.Spec.Replicas,
		CreationTime:   resp.CreationTimestamp.String(),
		Labels:         resp.Spec.Template.Labels,
	}, nil
}

func (k *K8sDeployment) GetAll(ctx context.Context, client *kubernetes.Clientset, clusterName, namespace string) ([]*K8sDeployment, rest_errors.RestErr) {

	resp, err := appsv1.Deployment.GetAll(ctx, client, namespace)
	if err != nil {
		return nil, rest_errors.NewBadRequestError(err.Error())
	}
	result := make([]*K8sDeployment, 0, len(resp))
	for _, deployment := range resp {
		result = append(result, &K8sDeployment{
			ClusterName:    clusterName,
			Namespace:      deployment.Namespace,
			DeploymentName: deployment.Name,
			Replicas:       *deployment.Spec.Replicas,
			CreationTime:   deployment.CreationTimestamp.String(),
			Labels:         deployment.Spec.Template.Labels,
		})
	}
	return result, nil
}

func (k *K8sDeployment) Create(ctx context.Context, client *kubernetes.Clientset) (*K8sDeployment, rest_errors.RestErr) {

	resp, err := appsv1.Deployment.Create(ctx,
		client,
		k.Namespace,
		k.DeploymentName,
		k.ContainerName,
		k.Image,
		&k.Replicas)
	if err != nil {
		return nil, rest_errors.NewBadRequestError(err.Error())
	}
	var result K8sDeployment
	result.DeploymentName = k.DeploymentName
	result.Status = fmt.Sprintf(statusMsg, creation)
	result.CreationTime = resp
	return &result, nil
}

func (k *K8sDeployment) Delete(ctx context.Context, client *kubernetes.Clientset) (*K8sDeployment, rest_errors.RestErr) {

	err := appsv1.Deployment.Delete(ctx,
		client,
		k.Namespace,
		k.DeploymentName)
	if err != nil {
		return nil, rest_errors.NewBadRequestError(err.Error())
	}
	k.Status = fmt.Sprintf(statusMsg, deletion)
	return k, nil
}

func (k *K8sDeployments) CreateMultiContainer(ctx context.Context, client *kubernetes.Clientset) (*K8sDeployments, rest_errors.RestErr) {

	resp, err := appsv1.Deployment.CreateMultiContainer(ctx,
		client, k.Namespace, k.DeploymentName, k.Containers, &k.Replicas)
	if err != nil {
		return nil, rest_errors.NewBadRequestError(err.Error())
	}
	var result K8sDeployments
	result.DeploymentName = k.DeploymentName
	result.Status = fmt.Sprintf(statusMsg, creation)
	result.CreationTime = resp
	return &result, nil
}
