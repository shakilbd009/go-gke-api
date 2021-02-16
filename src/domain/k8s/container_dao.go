package k8s

import (
	"context"

	"github.com/shakilbd009/go-gke-api/src/client/kubernetes/corev1"
	"github.com/shakilbd009/go-utils-lib/rest_errors"
	"k8s.io/client-go/kubernetes"
)

func (k *Kcontainer) GetContainers(ctx context.Context, client *kubernetes.Clientset, namespace, podName string) ([]Kcontainer, rest_errors.RestErr) {

	containers, err := corev1.Container.GetContainers(ctx, client, namespace, podName)
	if err != nil {
		return nil, rest_errors.NewBadRequestError(err.Error())
	}
	result := make([]Kcontainer, 0, len(containers))
	for _, container := range containers {
		// fmt.Printf("%#v", container)
		result = append(result, Kcontainer{
			ContainerNamespace: namespace,
			ContainerName:      container.Name,
			LivelinessProbe:    container.LivenessProbe.String(),
			//Port:               container.Ports[0].ContainerPort,
			Image:          container.Image,
			ReadynessProbe: container.ReadinessProbe.String(),
			PodName:        podName,
		})
	}
	return result, nil
}

func (k *Kcontainer) GetContainerLogs(ctx context.Context, client *kubernetes.Clientset, namespace, podName, container string) ([]string, rest_errors.RestErr) {

	logs, err := corev1.Container.GetContainerLogs(ctx, client, namespace, podName, container)
	if err != nil {
		return nil, rest_errors.NewBadRequestError(err.Error())
	}
	return logs, nil
}
