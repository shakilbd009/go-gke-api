package utility

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/shakilbd009/go-utils-lib/rest_errors"
)

func CheckQueryParam(c *gin.Context, queryParam ...string) (map[string]string, rest_errors.RestErr) {

	params := make(map[string]string, len(queryParam))
	for _, q := range queryParam {
		query, ok := c.GetQuery(q)
		if !ok {
			return nil, rest_errors.NewBadRequestError(fmt.Sprintf("%q is missing as query param", q))
		}
		params[q] = query
	}
	return params, nil
}
