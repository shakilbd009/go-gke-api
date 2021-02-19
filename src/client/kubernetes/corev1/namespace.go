package corev1

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

var (
	Namespace namespaceInterface = &namespace{}
)

type namespaceInterface interface {
	Get(context.Context, *kubernetes.Clientset, string) error
	GetAll(ctx context.Context, client *kubernetes.Clientset) ([]v1.Namespace, error)
}
type namespace struct{}

func (*namespace) Get(ctx context.Context, client *kubernetes.Clientset, namespaceName string) error {

	result, err := client.CoreV1().Namespaces().Get(ctx, namespaceName, metav1.GetOptions{})
	if err != nil {
		return err
	}
	if result.Name == metav1.NamespaceDefault || result.Name == metav1.NamespaceSystem {
		return fmt.Errorf(fmt.Sprintf("users are not permitted to use %s namespace", namespaceName))
	}
	return nil
}

func (*namespace) GetAll(ctx context.Context, client *kubernetes.Clientset) ([]v1.Namespace, error) {

	result, err := client.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return result.Items, nil
}
