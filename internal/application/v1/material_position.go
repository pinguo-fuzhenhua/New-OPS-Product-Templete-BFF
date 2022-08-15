package v1

import (
	"strconv"
	"strings"

	kerr "github.com/go-kratos/kratos/v2/errors"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/cparam"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/server"
	pver "github.com/pinguo-icc/go-base/v2/version"
	mpsvc "github.com/pinguo-icc/material-positions-svc/api"
)

type MaterialPositions struct {
	MP mpsvc.MaterialPositionsClient
}

// /v2/material-positions/{position}/categories?vipStatus=
func (m *MaterialPositions) Categories(ctx khttp.Context) (interface{}, error) {
	cp := cparam.FromContext(ctx)
	if cp == nil {
		return nil, kerr.BadRequest("parse common params failed", "解析公共请求参数发生错误")
	}

	clientVer, err := pver.Parse(cp.AppVersion)
	if err != nil {
		return nil, kerr.BadRequest(err.Error(), "请求参数版本号有误")
	}

	in := &mpsvc.PlaceCategoriesRequest{
		Scope:         cp.AppID,
		PosCode:       ctx.Vars().Get("position"),
		Platform:      cp.Platform,
		ClientVersion: int64(clientVer),
		UserData: &mpsvc.UserData{
			UserId:   cp.UserID,
			DeviceId: cp.EID,
			Properties: map[string]string{
				"language":  cp.Language,
				"locale":    cp.Locale,
				"vipstatus": ctx.Form().Get("vipStatus"),
			},
			UtcOffset: int32(cp.UtcOffset),
			Language:  cp.Language,
		},
		Prefetch: 72,
	}

	m.rewriteForUserData(ctx, in.UserData, cp)

	placeResp, err := m.MP.PlaceCategories(ctx, in)
	if err != nil {
		return nil, kerr.InternalServer(err.Error(), "服务器请求发生错误")
	}

	if len(placeResp.List) == 0 {
		return materialPosEmptyResp, nil
	}

	return placeResp.List, nil
}

// /v2/material-positions/{position}/materials?vipStatus=&cates=
func (m *MaterialPositions) Materials(ctx khttp.Context) (interface{}, error) {
	cp := cparam.FromContext(ctx)
	if cp == nil {
		return nil, kerr.BadRequest("parse common params failed", "解析公共请求参数发生错误")
	}

	clientVer, err := pver.Parse(cp.AppVersion)
	if err != nil {
		return nil, kerr.BadRequest(err.Error(), "请求参数版本号有误")
	}

	var cates []string
	if _cates := ctx.Query().Get("cates"); _cates != "" {
		cates = strings.Split(_cates, ",")
	}

	in := &mpsvc.PlaceMaterialsRequest{
		Scope:         cp.AppID,
		PosCode:       ctx.Vars().Get("position"),
		Platform:      cp.Platform,
		ClientVersion: int64(clientVer),
		UserData: &mpsvc.UserData{
			UserId:   cp.UserID,
			DeviceId: cp.EID,
			Properties: map[string]string{
				"language":  cp.Language,
				"locale":    cp.Locale,
				"vipstatus": ctx.Form().Get("vipStatus"),
			},
			Language:  cp.Language,
			UtcOffset: int32(cp.UtcOffset),
		},
		Prefetch: 72,
		Cates:    cates,
	}

	m.rewriteForUserData(ctx, in.UserData, cp)

	placeResp, err := m.MP.PlaceMaterialsV2(ctx, in)
	if err != nil {
		return nil, kerr.InternalServer(err.Error(), "服务器请求发生错误")
	}

	if len(placeResp.List) == 0 {
		return materialPosEmptyResp, nil
	}

	return placeResp.List, nil
}

func (m *MaterialPositions) rewriteForUserData(ctx khttp.Context, in *mpsvc.UserData, cp *cparam.Params) {
	header := ctx.Header()

	forNewUser := "0"
	if v := server.IsNewUser(cp, ctx.Request()); v {
		forNewUser = "1"
	}
	in.Properties["fornewuser"] = forNewUser

	if v := header.Get("Pg-Mock-Grayratio"); v != "" {
		if d, err := strconv.Atoi(v); err == nil && d > 0 {
			in.ForceGrayGroup = uint32(d)
		}
	}

	if v := header.Get("Pg-Mock-Vipstatus"); v != "" {
		in.Properties["vipstatus"] = v
	}

	// TODO  Pg-Mock-Usergroupid 精准用户群
}

var materialPosEmptyResp = []struct{}{}
