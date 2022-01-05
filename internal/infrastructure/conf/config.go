package conf

import (
	"time"

	lconfig "github.com/pinguo-icc/kratos-library/v2/conf"
)

type Bootstrap struct {
	Http   *HTTP
	Params *Params
}

type Params struct {
	ArticleSvcAddr string
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
