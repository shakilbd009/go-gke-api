package k8s

import (
	"context"

	"github.com/shakilbd009/go-gke-api/src/client/kubernetes/corev1"
	"github.com/shakilbd009/go-gke-api/src/domain/entity"
	"github.com/shakilbd009/go-utils-lib/rest_errors"
	"k8s.io/client-go/kubernetes"
)

func (kn *Knamespace) GetAll(ctx context.Context, client *kubernetes.Clientset, cluster string) ([]*Knamespace, rest_errors.RestErr) {

	namespaces, err := corev1.Namespace.GetAll(ctx, client)
	if err != nil {
		return nil, rest_errors.NewBadRequestError(err.Error())
	}
	result := make([]*Knamespace, 0, len(namespaces))
	for _, n := range namespaces {
		result = append(result, &Knamespace{
			Namespace:   n.Name,
			ClusterName: cluster,
			Request: &entity.Request{
				Project: kn.Request.Project,
				Region:  kn.Request.Region,
			},
		})
	}
	return result, nil
}
