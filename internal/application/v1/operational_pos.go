package v1

import (
	"strconv"
	"strings"

	kerr "github.com/go-kratos/kratos/v2/errors"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/pinguo-icc/Camera360/internal/domain"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/clientset"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/cparam"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/server"
	fdpkg "github.com/pinguo-icc/field-definitions/pkg"
	pver "github.com/pinguo-icc/go-base/v2/version"
	opapi "github.com/pinguo-icc/operational-positions-svc/api"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// 运营位

type OperationalPos struct {
	*clientset.ClientSet

	Parser *domain.ActivitiesParser
}

func (o *OperationalPos) PullByCodes(ctx khttp.Context) (interface{}, error) {
	posCodes := strings.TrimSpace(ctx.Form().Get("codes"))
	if posCodes == "" {
		return nil, kerr.BadRequest("invalid param, codes required", "invalid param")
	}

	cp := cparam.FromContext(ctx)
	cVer, err := pver.Parse(cp.AppVersion)
	if err != nil {
		return nil, kerr.BadRequest("invalid AppVersion", "invalid param")
	}

	in := &opapi.PlacingRequest{
		Prefetch:      72,
		Scope:         cp.AppID,
		PosCodes:      strings.Split(posCodes, ","),
		Platform:      cp.Platform,
		ClientVersion: int64(cVer),
		UserData: &opapi.UserData{
			UserId:    cp.UserID,
			DeviceId:  cp.EID,
			UtcOffset: int32(cp.UtcOffset),
			Properties: map[string]string{
				"language":  cp.Language,
				"locale":    cp.Locale,
				"vipstatus": ctx.Form().Get("vipStatus"),
			},
		},
	}

	v := server.IsNewUser(cp, ctx.Request())
	in.UserData.IsNewUser = wrapperspb.Bool(v)
	in.UserData.Properties["fornewuser"] = "0"
	if v {
		in.UserData.Properties["fornewuser"] = "1"
	}
	o.rewriteForPreview(ctx, in)
	langMatcher, err := fdpkg.NewLanguageMatcher(cp.Language, cp.Locale)
	if err != nil {
		return nil, kerr.BadRequest(err.Error(), "client language, locale invalid")
	}

	res, err := o.OperationalPositionsClient.Placing(ctx, in)
	if err != nil {
		return nil, kerr.InternalServer(err.Error(), "call service failed")
	}

	ret, err := o.Parser.Parse(ctx, langMatcher, res.Payload)
	if err != nil {
		return nil, kerr.InternalServer(err.Error(), "parse content failed")
	}
	return ret, nil
}

func (o *OperationalPos) rewriteForPreview(ctx khttp.Context, in *opapi.PlacingRequest) {
	header := ctx.Request().Header
	if v := header.Get("Pg-Mock-Grayratio"); v != "" {
		if d, err := strconv.Atoi(v); err == nil && d > 0 {
			in.UserData.ForceGrayGroup = uint32(d)
		}
	}

	if v := header.Get("Pg-Mock-Vipstatus"); v != "" {
		in.UserData.Properties["vipstatus"] = v
	}

	// TODO Pg-Mock-Usergroupid 精准用户群组
}
