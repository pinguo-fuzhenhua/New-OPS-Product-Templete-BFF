package v1

import (
	"encoding/json"
	kerr "github.com/go-kratos/kratos/v2/errors"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/cparam"
	opbasic "github.com/pinguo-icc/operational-basic-svc/api"
	"strings"
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
	json.Unmarshal([]byte(data.Content), &res)
	return res, nil
}
