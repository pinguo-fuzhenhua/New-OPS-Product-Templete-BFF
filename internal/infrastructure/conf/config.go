package conf

import (
	"time"

	lconfig "github.com/pinguo-icc/kratos-library/v2/conf"
	"github.com/pinguo-icc/kratos-library/v2/trace"
)

type Bootstrap struct {
	App       *App
	Http      *HTTP
	Trace     *trace.Config
	Clientset *Clientset
	Params    *Params
}

type Params struct{}

func Load(env string) (*Bootstrap, error) {
	out := new(Bootstrap)
	err := lconfig.Load(env, out, nil)
	return out, err
}

type HTTP struct {
	Address string
	Timeout time.Duration
}

type App struct {
	Name string
	Env  string
}

type Clientset struct {
	FieldDef       string
	OperationalPos string
	Material       string
}
