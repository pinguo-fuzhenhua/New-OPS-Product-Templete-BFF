package render

import (
	kerr "github.com/go-kratos/kratos/v2/errors"
)

var (
	ErrLoginRequired = kerr.Forbidden("", "Login Required")

	loginRequired, _ = MarshalOptions.Marshal(ErrLoginRequired)
)
