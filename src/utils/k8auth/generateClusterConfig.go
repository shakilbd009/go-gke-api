package k8auth

import (
	"context"
	"encoding/base64"
	"fmt"

	ctn "cloud.google.com/go/container/apiv1"
	"github.com/shakilbd009/go-utils-lib/rest_errors"
	"google.golang.org/genproto/googleapis/container/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func GetGkekubeConfig(ctx context.Context, projectID, region, clusterName string) (*kubernetes.Clientset, rest_errors.RestErr) {

	kubeConfig, rest_err := GetGKEClustersConfig(ctx, projectID, region)
	if rest_err != nil {
		return nil, rest_err
	}
	cluster := fmt.Sprintf("gke_%s_%s_%s", projectID, region, clusterName)
	kcfg, err := clientcmd.NewNonInteractiveClientConfig(*kubeConfig, cluster, &clientcmd.ConfigOverrides{
		CurrentContext: cluster,
	}, nil).ClientConfig()
	if err != nil {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("kubeConfig could not be generated for %s, err: %s", clusterName, err.Error()))
	}
	client, err := kubernetes.NewForConfig(kcfg)
	if rest_err != nil {
		return nil, rest_err
	}
	return client, nil
}

func ListGKEWorkload(ctx context.Context, projectID, region string) ([]v1.Namespace, rest_errors.RestErr) {
	kubeConfig, err := GetGKEClustersConfig(ctx, projectID, region)
	if err != nil {
		return nil, err
	}
	var namespaces = make([]v1.Namespace, 0)
	for clusterName := range kubeConfig.Clusters {

		kcfg, err := clientcmd.NewNonInteractiveClientConfig(*kubeConfig, clusterName, &clientcmd.ConfigOverrides{
			CurrentContext: clusterName,
		}, nil).ClientConfig()
		if err != nil {
			return nil, rest_errors.NewBadRequestError(fmt.Sprintf("failed to create Kubernetes configuration cluster=%s: %s", clusterName, err.Error()))
		}
		k8s, err := kubernetes.NewForConfig(kcfg)
		if err != nil {
			return nil, rest_errors.NewInternalServerError(fmt.Sprintf("failed to create Kubernetes client cluster=%s: %s", clusterName, err.Error()), err)
		}
		ns, err := k8s.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
		if err != nil {
			return nil, rest_errors.NewInternalServerError(fmt.Sprintf("failed to list namespaces cluster=%s: %s", clusterName, err.Error()), err)
		}
		for _, n := range ns.Items {
			namespaces = append(namespaces, n)
		}
	}
	return namespaces, nil
}

func GetGKEClustersConfig(ctx context.Context, projectID, region string) (*api.Config, rest_errors.RestErr) {
	client, err := ctn.NewClusterManagerClient(ctx)
	if err != nil {
		return nil, rest_errors.NewInternalServerError(err.Error(), err)
	}
	cfg := api.Config{
		APIVersion: "v1",
		Kind:       "config",
		Clusters:   map[string]*api.Cluster{},
		AuthInfos:  map[string]*api.AuthInfo{},
		Contexts:   map[string]*api.Context{},
	}
	parent := fmt.Sprintf("projects/%s/locations/%s", projectID, region)
	resp, err := client.ListClusters(ctx, &container.ListClustersRequest{Parent: parent})
	if err != nil {
		return nil, rest_errors.NewInternalServerError(err.Error(), err)
	}
	for _, cluster := range resp.Clusters {
		name := fmt.Sprintf("gke_%s_%s_%s", projectID, region, cluster.Name)
		cert, err := base64.StdEncoding.DecodeString(cluster.MasterAuth.ClusterCaCertificate)
		if err != nil {
			return nil, rest_errors.NewInternalServerError(fmt.Sprintf("invalid certificate cluster=%s cert=%s: %w", name, cluster.MasterAuth.ClusterCaCertificate, err), err)
		}
		cfg.Clusters[name] = &api.Cluster{
			CertificateAuthorityData: cert,
			Server:                   "https://" + cluster.Endpoint,
		}
		cfg.Contexts[name] = &api.Context{
			Cluster:  name,
			AuthInfo: name,
		}
		cfg.AuthInfos[name] = &api.AuthInfo{
			AuthProvider: &api.AuthProviderConfig{
				Name: "gcp",
				Config: map[string]string{
					"scopes": "https://www.googleapis.com/auth/cloud-platform",
				},
			},
		}
	}
	return &cfg, nil
}
