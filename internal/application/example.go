package application

import (
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/clientset"
)

type Example struct {
	*clientset.ClientSet
}

func (e *Example) Get(ctx Context) (interface{}, error) {
	return []string{"example", "get"}, nil
}

func (e *Example) Post(ctx Context) error {
	name, _ := PathParam(ctx, "name")
	if name == "error" {
		return errors.BadRequest("reason", "response message")
	}

	ctx.JSON(200, []string{"name: ", name})
	return nil
}
