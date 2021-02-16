package corev1

import (
	"bytes"
	"context"
	"io"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	Container containerIinterface = &container{}
)

type container struct{}
type containerIinterface interface {
	GetContainerLogs(ctx context.Context, client *kubernetes.Clientset, namespace, podName, container string) ([]string, error)
	GetContainers(ctx context.Context, client *kubernetes.Clientset, namespace, podName string) ([]corev1.Container, error)
}

func (*container) GetContainers(ctx context.Context, client *kubernetes.Clientset, namespace, podName string) ([]corev1.Container, error) {

	pod, err := client.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pod.Spec.Containers, nil
}

func (*container) GetContainerLogs(ctx context.Context, client *kubernetes.Clientset, namespace, podName, container string) ([]string, error) {

	lines := int64(5)
	request := client.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
		TailLines: &lines,
		Container: container,
	})
	stream, err := request.Stream(ctx)
	if err != nil {
		return nil, err
	}
	defer stream.Close()
	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, stream)
	if err != nil {
		return nil, err
	}
	logs := strings.Split(buf.String(), "\n")
	return logs, nil
}
