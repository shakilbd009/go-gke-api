package gked

import (
	"fmt"

	"github.com/shakilbd009/go-gke-api/src/domain/entity"
	"github.com/shakilbd009/go-utils-lib/rest_errors"
)

var (
	msgTemplate    = "%v field is missing values"
	required_field = "required fields are missing"
	project        = "project"
	region         = "region"
)

type GkeCluster struct {
	*entity.Request
	NodePools     []string `json:"nodepools"`
	NodeCount     int32    `json:"nodeCount"`
	MasterVersion string   `json:"version"`
	ClusterName   string   `json:"clusterName"`
}

func (g *GkeCluster) Validate() rest_errors.RestErr {
	if g == nil {
		return rest_errors.NewBadRequestError(required_field)
	}
	if g.Project == "" {
		return rest_errors.NewBadRequestError(fmt.Sprintf(msgTemplate, project))
	}
	if g.Region == "" {
		return rest_errors.NewBadRequestError(fmt.Sprintf(msgTemplate, region))
	}
	return nil
}
