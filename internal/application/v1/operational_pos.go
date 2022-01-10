package v1

import (
	"strconv"
	"strings"

	kerr "github.com/go-kratos/kratos/v2/errors"
	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/pinguo-icc/Camera360/internal/domain"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/clientset"
	"github.com/pinguo-icc/Camera360/internal/infrastructure/cparam"
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
		return nil, kerr.BadRequest("invalid param", "invalid param")
	}

	cp := cparam.FromContext(ctx)

	cVer, err := pver.Parse(cp.AppVersion)
	if err != nil {
		return nil, kerr.BadRequest("invalid AppVersion", "invalid param")
	}

	in := &opapi.PlacingRequest{
		Scope:         "Camera360",
		PosCodes:      strings.Split(posCodes, ","),
		Platform:      cp.Platform,
		ClientVersion: int64(cVer),
		UserData: &opapi.UserData{
			UserId:   cp.UserID,
			DeviceId: cp.EID,
			Properties: map[string]string{
				"language": cp.Language,
				"locale":   cp.Locale,
			},
		},
	}

	if isNewUser := ctx.Form().Get("isNewUser"); isNewUser != "" {
		b, err := strconv.ParseBool(isNewUser)
		if err == nil {
			in.UserData.IsNewUser = wrapperspb.Bool(b)
		}
	}

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
