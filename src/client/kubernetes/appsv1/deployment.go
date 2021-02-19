package appsv1

import (
	"context"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	appName = "app"
	http    = "http"
	port80  = 80
)

var (
	Deployment deploymentInterface = &deployment{}
)

type deploymentInterface interface {
	Get(ctx context.Context, client *kubernetes.Clientset, namespace, deploymentName string) (*appsv1.Deployment, error)
	GetAll(ctx context.Context, client *kubernetes.Clientset, namespace string) ([]appsv1.Deployment, error)
	Create(ctx context.Context, client *kubernetes.Clientset, namespace, deploymentName, containerName, image string, replicas *int32) (string, error)
	Delete(ctx context.Context, client *kubernetes.Clientset, namespace, deploymentName string) error
	CreateMultiContainer(ctx context.Context, client *kubernetes.Clientset, namespace, deploymentName string, containers []Container, replicas *int32) (string, error)
}
type Container struct {
	Name    string   `json:"name,omitempty"`
	Image   string   `json:"image,omitempty"`
	Command []string `json:"commands,omitempty"`
}
type deployment struct{}

func (*deployment) Get(ctx context.Context, client *kubernetes.Clientset, namespace, deploymentName string) (*appsv1.Deployment, error) {
	resp, err := client.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (*deployment) GetAll(ctx context.Context, client *kubernetes.Clientset, namespace string) ([]appsv1.Deployment, error) {
	resp, err := client.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func (*deployment) Create(ctx context.Context, client *kubernetes.Clientset, namespace, deploymentName, containerName, image string, replicas *int32) (string, error) {

	param := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					appName: containerName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						appName: containerName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  containerName,
							Image: image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          http,
									ContainerPort: port80,
								},
							},
						},
					},
				},
			},
		},
	}
	result, err := client.AppsV1().Deployments(namespace).Create(ctx, &param, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	return result.CreationTimestamp.String(), nil
	//return result.Status.String(), nil
}

func (*deployment) CreateMultiContainer(ctx context.Context, client *kubernetes.Clientset, namespace, deploymentName string, containers []Container, replicas *int32) (string, error) {

	var containerName string
	apiContainers := make([]apiv1.Container, 0)

	for _, c := range containers {
		containerName += strings.Title(c.Name)

	}
	for _, c := range containers {
		cont := apiv1.Container{
			Name:    c.Name,
			Image:   c.Image,
			Command: c.Command,
			// Ports: []apiv1.ContainerPort{
			// 	{
			// 		Name:          http,
			// 		ContainerPort: port80,
			// 	},
			// },

		}
		apiContainers = append(apiContainers, cont)
	}

	param := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					appName: containerName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						appName: containerName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: apiContainers,
					//RestartPolicy: apiv1.RestartPolicy("Never"),
				},
			},
		},
	}
	result, err := client.AppsV1().Deployments(namespace).Create(ctx, &param, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	return result.CreationTimestamp.String(), nil
	//return result.Status.String(), nil
}

func (*deployment) Delete(ctx context.Context, client *kubernetes.Clientset, namespace, deploymentName string) error {
	return client.AppsV1().Deployments(namespace).Delete(ctx, deploymentName, metav1.DeleteOptions{})
}
