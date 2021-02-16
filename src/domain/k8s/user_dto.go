package k8s

import (
	"fmt"

	"github.com/shakilbd009/go-utils-lib/rest_errors"
)

var (
	msgEmpty = "%s field is empty"
	userName = "user_name"
	email    = "email"
)

type K8sUser struct {
	UserName  string `json:"user_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Token     string `json:"token,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

func (k *K8sUser) Validate() rest_errors.RestErr {

	if k.UserName == "" {
		return rest_errors.NewBadRequestError(fmt.Sprintf(msgEmpty, userName))
	}
	if k.Email == "" {
		return rest_errors.NewBadRequestError(fmt.Sprintf(msgEmpty, email))
	}
	if k.Namespace == "" {
		return rest_errors.NewBadRequestError(fmt.Sprintf(msgEmpty, namespace))
	}
	return nil
}
