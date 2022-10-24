package v1

import (
	"encoding/json"
	"strings"

	kerr "github.com/go-kratos/kratos/v2/errors"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	opbasic "github.com/pinguo-icc/operational-basic-svc/api"
	"github.com/pinguo-icc/template/internal/infrastructure/cparam"
)

type JsonConfig struct {
	Opbasic opbasic.OperationalBasicClient
}

func (j *JsonConfig) Show(ctx khttp.Context) (res interface{}, err error) {
	cp := cparam.FromContext(ctx)
	if cp == nil {
		return nil, kerr.BadRequest("parse common params failed", "解析公共请求参数发生错误")
	}
	appName := cp.AppID
	code := strings.TrimSpace(ctx.Form().Get("code"))
	if code == "" {
		return nil, kerr.BadRequest("invalid param, code required", "invalid param")
	}
	req := &opbasic.JsonConfigShowRequest{
		AppName: appName,
		Code:    code,
	}
	data, err := j.Opbasic.JsonConfigShow(ctx, req)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(data.Content), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
