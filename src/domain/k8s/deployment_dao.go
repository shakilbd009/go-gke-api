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
