package v1

import (
	"strconv"
	"strings"

	kerr "github.com/go-kratos/kratos/v2/errors"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	pver "github.com/pinguo-icc/go-base/v2/version"
	opmapi "github.com/pinguo-icc/operations-material-svc/api"
	"github.com/pinguo-icc/template/internal/infrastructure/cparam"
	"github.com/pinguo-icc/template/internal/infrastructure/server"
)

type MaterialPositions struct {
	MP opmapi.MaterialPositionsClient
}

type CategoryResponse struct {
	Items []*opmapi.PlaceResponse_Plan `json:"items"`
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

	in := &opmapi.PlaceCategoriesRequest{
		Scope:         cp.AppID,
		PosCode:       ctx.Vars().Get("position"),
		Platform:      cp.Platform,
		ClientVersion: int64(clientVer),
		UserData: &opmapi.UserData{
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
		return &CategoryResponse{Items: []*opmapi.PlaceResponse_Plan{}}, nil
	}

	return &CategoryResponse{Items: placeResp.List}, nil
}

type MaterialRequest struct {
	VIPStatus string `json:"vipStatus"`
	Cates     string `json:"cates"`
	PageType  string `json:"pageType"`
	ScrollID  string `json:"scrollID"`
	PageSize  int    `json:"pageSize"`
}
type MaterialResponse struct {
	ScrollID string                       `json:"scrollID"`
	Items    []*opmapi.PlaceResponse_Plan `json:"items"`
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

	req := new(MaterialRequest)
	if err := ctx.BindQuery(req); err != nil {
		return nil, err
	}

	var cates []string
	if _cates := req.Cates; _cates != "" {
		cates = strings.Split(_cates, ",")
	}

	in := &opmapi.PlaceMaterialsRequest{
		Scope:         cp.AppID,
		PosCode:       ctx.Vars().Get("position"),
		Platform:      cp.Platform,
		ClientVersion: int64(clientVer),
		UserData: &opmapi.UserData{
			UserId:   cp.UserID,
			DeviceId: cp.EID,
			Properties: map[string]string{
				"language":  cp.Language,
				"locale":    cp.Locale,
				"vipstatus": req.VIPStatus,
			},
			Language:  cp.Language,
			UtcOffset: int32(cp.UtcOffset),
		},
		Prefetch: 72,
		Cates:    cates,
		PageType: req.PageType,
		ScrollID: req.ScrollID,
		PageSize: int64(req.PageSize),
	}

	m.rewriteForUserData(ctx, in.UserData, cp)

	placeResp, err := m.MP.PlaceMaterialsV2(ctx, in)
	if err != nil {
		return nil, kerr.InternalServer(err.Error(), "服务器请求发生错误")
	}

	if len(placeResp.List) == 0 {
		return &MaterialResponse{ScrollID: "", Items: []*opmapi.PlaceResponse_Plan{}}, nil
	}

	return &MaterialResponse{ScrollID: placeResp.ScrollID, Items: placeResp.List}, nil
}

func (m *MaterialPositions) rewriteForUserData(ctx khttp.Context, in *opmapi.UserData, cp *cparam.Params) {
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

func (m *MaterialPositions) MaterialDetail(ctx khttp.Context) (any, error) {
	type MaterialDetailRequest struct {
		Ids      string `json:"ids"`
		WithRefs bool   `json:"withRefs"`
		Filter   bool   `json:"filter"`
	}

	req := new(MaterialDetailRequest)
	err := ctx.BindQuery(req)
	if err != nil {
		return nil, err
	}
	if req.Ids == "" {
		return nil, kerr.BadRequest("ids is empty", "ids 不能为空")
	}

	cp := cparam.FromContext(ctx)
	if cp == nil {
		return nil, kerr.BadRequest("parse common params failed", "解析公共请求参数发生错误")
	}

	clientVer, err := pver.Parse(cp.AppVersion)
	if err != nil {
		return nil, kerr.BadRequest(err.Error(), "请求参数版本号有误")
	}

	materialIDs := strings.Split(req.Ids, ",")
	mdres, err := m.MP.MaterialDetail(ctx, &opmapi.MaterialDetailRequest{
		Scope:         cp.AppID,
		Platform:      cp.Platform,
		ClientVersion: int64(clientVer),
		WithRefes:     req.WithRefs,
		Ids:           materialIDs,
		Language:      cp.Language,
		Filter:        req.Filter,
	})
	if err != nil {
		return nil, err
	}
	if len(mdres.Materials) == 0 {
		return materialPosEmptyResp, nil
	}
	return mdres.Materials, nil
}
