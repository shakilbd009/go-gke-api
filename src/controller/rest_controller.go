package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shakilbd009/go-gke-api/src/domain/entity"
	"github.com/shakilbd009/go-gke-api/src/domain/k8s"
	"github.com/shakilbd009/go-gke-api/src/services"
	"github.com/shakilbd009/go-gke-api/src/utils/k8auth"
	"github.com/shakilbd009/go-gke-api/src/utils/utility"
	"github.com/shakilbd009/go-utils-lib/rest_errors"
)

var (
	Kcontroller   kcontrollerInterface = &kcontroller{}
	invalidJSON                        = "invalid json body"
	namespace                          = "namespace"
	projectID                          = "project"
	region                             = "region"
	clusterName                        = "clusterName"
	podName                            = "podName"
	containerName                      = "containerName"
)

type kcontroller struct{}

type kcontrollerInterface interface {
	CreateDeployment(*gin.Context)
	CreateMultiContainerDeployment(c *gin.Context)
	DeleteDeployment(*gin.Context)
	GetPods(*gin.Context)
	ListGKE(c *gin.Context)
	GetClusters(c *gin.Context)
	GetClustersDeployments(c *gin.Context)
	GetNamespaces(c *gin.Context)
	GetPodLogs(*gin.Context)
	GetContainers(c *gin.Context)
	GetContainerLogs(c *gin.Context)
}

func (k *kcontroller) GetClustersDeployments(c *gin.Context) {

	params, err := utility.CheckQueryParam(c, projectID, region)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	var e entity.Request
	e.Project = params[projectID]
	e.Region = params[region]
	result, err := services.KCommonServices.GetAllDeployments(c.Request.Context(), &e)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (k *kcontroller) GetNamespaces(c *gin.Context) {

	params, err := utility.CheckQueryParam(c, projectID, region, clusterName)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	var e entity.Request
	e.Project = params[projectID]
	e.Region = params[region]
	result, err := services.KnamespaceServices.GetAll(c.Request.Context(), &e, params[clusterName])
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (k *kcontroller) GetClusters(c *gin.Context) {

	params, err := utility.CheckQueryParam(c, projectID, region)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	var e entity.Request
	e.Project = params[projectID]
	e.Region = params[region]
	result, err := services.KgkeServices.GetClusters(c.Request.Context(), &e)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (k *kcontroller) ListGKE(c *gin.Context) {

	params, err := utility.CheckQueryParam(c, projectID, region)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	namespaces, err := k8auth.ListGKEWorkload(c.Request.Context(), params[projectID], params[region])
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, namespaces)
}

func (k *kcontroller) CreateDeployment(c *gin.Context) {
	var deployment k8s.K8sDeployment
	if err := c.ShouldBindJSON(&deployment); err != nil {
		fmt.Println(err)
		restErr := rest_errors.NewBadRequestError(invalidJSON)
		c.JSON(restErr.Status(), restErr)
		return
	}
	result, err := services.KDeploymentServices.CreateDeployment(c.Request.Context(), &deployment)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusCreated, *result)
}

func (k *kcontroller) CreateMultiContainerDeployment(c *gin.Context) {
	var deployment k8s.K8sDeployments
	if err := c.ShouldBindJSON(&deployment); err != nil {
		fmt.Println(err)
		restErr := rest_errors.NewBadRequestError(invalidJSON)
		c.JSON(restErr.Status(), restErr)
		return
	}
	result, err := services.KDeploymentServices.CreateMultiContainerDeployment(c.Request.Context(), &deployment)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusCreated, *result)
}

func (k *kcontroller) DeleteDeployment(c *gin.Context) {

	var deployment k8s.K8sDeployment
	if err := c.ShouldBindJSON(&deployment); err != nil {
		restErr := rest_errors.NewBadRequestError(invalidJSON)
		c.JSON(restErr.Status(), restErr)
		return
	}

	result, err := services.KDeploymentServices.DeleteDeployment(c.Request.Context(), &deployment)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusCreated, *result)
}

func (k *kcontroller) GetPods(c *gin.Context) {

	params, err := utility.CheckQueryParam(c, projectID, region, clusterName, namespace)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	result, err := services.KPodServices.GetPods(c.Request.Context(), params[projectID], params[region], params[clusterName], params[namespace])
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	if len(result) == 0 {
		restErr := rest_errors.NewNotFoundError(fmt.Sprintf("No resources found in %s namespace", namespace))
		c.JSON(restErr.Status(), restErr)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (k *kcontroller) GetPodLogs(c *gin.Context) {

	params, err := utility.CheckQueryParam(c,
		projectID,
		region,
		clusterName,
		namespace,
		podName)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	logs, err := services.KPodServices.GetPodLogs(c.Request.Context(),
		params[projectID],
		params[region],
		params[clusterName],
		params[namespace],
		params[podName])
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, logs)
}

func (k *kcontroller) GetContainers(c *gin.Context) {

	params, err := utility.CheckQueryParam(c,
		projectID,
		region,
		clusterName,
		namespace,
		podName)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	result, err := services.KContainerServices.GetContainers(c.Request.Context(), params[projectID],
		params[region],
		params[clusterName],
		params[namespace],
		params[podName])
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	if len(result) == 0 {
		restErr := rest_errors.NewNotFoundError(fmt.Sprintf("No resources found in %s namespace", namespace))
		c.JSON(restErr.Status(), restErr)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (k *kcontroller) GetContainerLogs(c *gin.Context) {

	params, err := utility.CheckQueryParam(c,
		projectID,
		region,
		clusterName,
		namespace,
		podName,
		containerName)
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	logs, err := services.KContainerServices.GetContainerLogs(c.Request.Context(), params[projectID],
		params[region],
		params[clusterName],
		params[namespace],
		params[podName],
		params[containerName])
	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, logs)
}
