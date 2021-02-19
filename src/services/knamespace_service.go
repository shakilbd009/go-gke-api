package services

import (
	"context"

	"github.com/shakilbd009/go-gke-api/src/domain/entity"
	"github.com/shakilbd009/go-gke-api/src/domain/k8s"
	"github.com/shakilbd009/go-gke-api/src/utils/k8auth"
	"github.com/shakilbd009/go-utils-lib/rest_errors"
)

var (
	KnamespaceServices knamespaceInterface = &knamespaceService{}
)

type knamespaceService struct{}

type knamespaceInterface interface {
	GetAll(ctx context.Context, e *entity.Request, clusterName string) ([]*k8s.Knamespace, rest_errors.RestErr)
}

func (*knamespaceService) GetAll(ctx context.Context, e *entity.Request, clusterName string) ([]*k8s.Knamespace, rest_errors.RestErr) {

	client, err := k8auth.GetGkekubeConfig(ctx, e.Project, e.Region, clusterName)
	if err != nil {
		return nil, err
	}
	kn := &k8s.Knamespace{
		Request: e,
	}
	result, err := kn.GetAll(ctx, client, clusterName)
	if err != nil {
		return nil, err
	}
	return result, nil
}
