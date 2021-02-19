package services

import (
	"context"

	"github.com/shakilbd009/go-gke-api/src/domain/entity"
	"github.com/shakilbd009/go-gke-api/src/domain/gked"
	"github.com/shakilbd009/go-utils-lib/rest_errors"
)

var (
	KgkeServices kgkeInterface = &kgkeService{}
)

type kgkeService struct{}

type kgkeInterface interface {
	GetClusters(ctx context.Context, n *entity.Request) ([]*gked.GkeCluster, rest_errors.RestErr)
}

func (*kgkeService) GetClusters(ctx context.Context, n *entity.Request) ([]*gked.GkeCluster, rest_errors.RestErr) {

	g := &gked.GkeCluster{
		Request: n,
	}
	if err := g.Validate(); err != nil {
		return nil, err
	}
	resp, err := g.GetClusters(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
