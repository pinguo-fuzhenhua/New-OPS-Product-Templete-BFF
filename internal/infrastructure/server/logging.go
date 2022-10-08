package server

import (
	"fmt"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// extractError returns the string of the error
func extractError(err error) (log.Level, string) {
	if err != nil {
		if e := errors.FromError(err); e != nil {
			if e.Code < 500 {
				return log.LevelWarn, fmt.Sprintf("%+v", err)
			}
			return log.LevelError, fmt.Sprintf("%+v", err)
		}
		return log.LevelError, fmt.Sprintf("%+v", err)
	}
	return log.LevelInfo, ""
}
