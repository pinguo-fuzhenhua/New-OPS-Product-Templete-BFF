package render

import (
	"net/http"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var MarshalOptions = protojson.MarshalOptions{
	UseEnumNumbers:  true,
	EmitUnpopulated: true,
}

type ErrorJSON struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func RenderJSON(ctx khttp.Context, data interface{}) error {
	if pbmsg, ok := data.(proto.Message); ok {
		buf, err := MarshalOptions.Marshal(pbmsg)
		if err != nil {
			return err
		}
		ctx.Blob(200, "application/json", buf)
		return nil
	}

	ctx.JSON(200, data)
	return nil
}

func LoginRequired(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)
	w.Write(loginRequired)
}
