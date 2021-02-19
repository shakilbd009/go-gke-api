package services

import (
	"context"
	"strings"

	"github.com/shakilbd009/go-gke-api/src/domain/entity"
	"github.com/shakilbd009/go-gke-api/src/domain/k8s"
	"github.com/shakilbd009/go-gke-api/src/utils/k8auth"
	"github.com/shakilbd009/go-k8s-api/src/client/kubernetes/corev1"

	"github.com/shakilbd009/go-k8s-api/src/utils/utility"
	"github.com/shakilbd009/go-utils-lib/rest_errors"
)

var (
	KDeploymentServices kDeploymentInterface = &kDeploymentService{}
	msg                                      = "no route to host"
)

type kDeploymentService struct{}

type kDeploymentInterface interface {
	CreateDeployment(context.Context, *k8s.K8sDeployment) (*k8s.K8sDeployment, rest_errors.RestErr)
	DeleteDeployment(context.Context, *k8s.K8sDeployment) (*k8s.K8sDeployment, rest_errors.RestErr)
	CreateMultiContainerDeployment(ctx context.Context, k *k8s.K8sDeployments) (*k8s.K8sDeployments, rest_errors.RestErr)
	GetAll(ctx context.Context, r *entity.Request, clusterName, namespace string) ([]*k8s.K8sDeployment, rest_errors.RestErr)
}

func (*kDeploymentService) GetAll(ctx context.Context, r *entity.Request, clusterName, namespace string) ([]*k8s.K8sDeployment, rest_errors.RestErr) {

	var k k8s.K8sDeployment
	client, err := k8auth.GetGkekubeConfig(ctx, r.Project, r.Region, clusterName)
	if err != nil {
		return nil, err
	}
	resp, err := k.GetAll(ctx, client, clusterName, namespace)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (*kDeploymentService) CreateDeployment(ctx context.Context, k *k8s.K8sDeployment) (*k8s.K8sDeployment, rest_errors.RestErr) {

	if err := k.ValidateCreateDeployment(); err != nil {
		return nil, err
	}
	//_, err := AuthService.GetToken(k.Token)
	// if err != nil {
	// 	return nil, err
	// }
	client, rest_err := k8auth.GetGkekubeConfig(ctx, k.Request.Project, k.Request.Region, k.ClusterName)
	if rest_err != nil {
		return nil, rest_err
	}
	if err := corev1.Namespace.Get(ctx, client, k.Namespace); err != nil {
		if strings.Contains(err.Error(), msg) {
			return nil, rest_errors.NewInternalServerError("Server unavailable", utility.ErrDatabase)
		}
		return nil, rest_errors.NewBadRequestError(err.Error())
	}

	result, err := k.Create(ctx, client)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (*kDeploymentService) CreateMultiContainerDeployment(ctx context.Context, k *k8s.K8sDeployments) (*k8s.K8sDeployments, rest_errors.RestErr) {

	if err := k.ValidateCreateDeployment(); err != nil {
		return nil, err
	}
	client, rest_err := k8auth.GetGkekubeConfig(ctx, k.Request.Project, k.Request.Region, k.ClusterName)
	if rest_err != nil {
		return nil, rest_err
	}
	if err := corev1.Namespace.Get(ctx, client, k.Namespace); err != nil {
		if strings.Contains(err.Error(), msg) {
			return nil, rest_errors.NewInternalServerError("Server unavailable", utility.ErrDatabase)
		}
		return nil, rest_errors.NewBadRequestError(err.Error())
	}

	result, err := k.CreateMultiContainer(ctx, client)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (*kDeploymentService) DeleteDeployment(ctx context.Context, k *k8s.K8sDeployment) (*k8s.K8sDeployment, rest_errors.RestErr) {

	if err := k.ValidateDeleteDeployment(); err != nil {
		return nil, err
	}
	client, rest_err := k8auth.GetGkekubeConfig(ctx, k.Request.Project, k.Request.Region, k.ClusterName)
	if rest_err != nil {
		return nil, rest_err
	}
	if err := corev1.Namespace.Get(ctx, client, k.Namespace); err != nil {
		return nil, rest_errors.NewBadRequestError(err.Error())
	}
	var resp k8s.K8sDeployment
	result, err := k.Delete(ctx, client)
	if err != nil {
		return nil, err
	}
	resp.DeploymentName = result.DeploymentName
	resp.Namespace = result.Namespace
	resp.Status = result.Status
	return &resp, nil
}
