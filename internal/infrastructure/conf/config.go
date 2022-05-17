package conf

import (
	"time"

	"github.com/pinguo-icc/go-lib/v2/uploader/qiniu"
	lconfig "github.com/pinguo-icc/kratos-library/v2/conf"
	"github.com/pinguo-icc/kratos-library/v2/trace"
)

const Scope = "Camera360"

type Bootstrap struct {
	App       *App
	Http      *HTTP
	Trace     *trace.Config
	Clientset *Clientset
	Qiniu     *qiniu.Config
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
<<<<<<< HEAD
	FieldDef                string
	OperationalPos          string
	Material                string
	DataEnv                 string
=======
	FieldDef       string
	OperationalPos string
	Material       string
	DataEnv        string
>>>>>>> ad8536911d77e6c95fe895c44208b20776dd7743
	OperationalBasicSvcAddr string
}
