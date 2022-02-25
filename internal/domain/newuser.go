package domain

import (
	"net/http"
	"time"

	"github.com/pinguo-icc/Camera360/internal/infrastructure/cparam"
)

// 各产品新老用户判断规则

// 判断规则定义，不同产品可定义不同的规则
var newUserRules = map[string]func(cp *cparam.Params, r *http.Request) bool{
	"Camera360": func(cp *cparam.Params, _ *http.Request) bool {
		return newUserDependOnInitTimestamp(cp.InitStamp, 1)
	},
}

// IsNewUser 判断请求是否为新用户
// 若给定产品未定义规则将使用默认的规则，即：安装时间（Header 参数：InitStamp）在一天内视作新用户
func IsNewUser(cp *cparam.Params, r *http.Request) bool {
	if fn, ok := newUserRules[cp.AppID]; ok {
		return fn(cp, r)
	}
	return newUserDependOnInitTimestamp(cp.InitStamp, 1)
}

func newUserDependOnInitTimestamp(initTime, diffDay int64) bool {
	return initTime+diffDay*24*3600 < time.Now().Unix()
}
