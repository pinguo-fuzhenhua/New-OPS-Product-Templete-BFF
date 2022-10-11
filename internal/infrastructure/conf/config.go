package conf

import (
	"time"

	"github.com/pinguo-icc/go-lib/v2/uploader/qiniu"
	lconfig "github.com/pinguo-icc/kratos-library/v2/conf"
	"github.com/pinguo-icc/kratos-library/v2/trace"
)

const Scope = "April"

type Bootstrap struct {
	App       *App
	Http      *HTTP
	Trace     *trace.Config
	Clientset *Clientset
	Qiniu     *qiniu.Config
	Params    *Params
	Recorder  *Recorder
	HTML5     *HTML5Config
}

type Params struct{}

type HTML5Config struct {
	HTML5URLPrefix string
}

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
	FieldDef                string
	OperationalPos          string
	Material                string
	DataEnv                 string
	OperationalBasicSvcAddr string
}

type Recorder struct {
	FilePath             string
	MaxSize              int
	MaxAge               int
	MaxBackups           int
	RegexpRecordeRouters []string
}
