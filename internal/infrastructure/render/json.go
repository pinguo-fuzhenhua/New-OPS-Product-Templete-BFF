package render

import (
	"encoding/json"
	"net/http"
	"text/template"

	khttp "github.com/go-kratos/kratos/v2/transport/http"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var MarshalOptions = protojson.MarshalOptions{
	UseEnumNumbers:  true,
	EmitUnpopulated: true,
}

type JSON struct {
	Data    interface{} `json:"data"`
	Status  int         `json:"status"`
	Message string      `json:"message"`
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

func RenderJSONP(ctx khttp.Context, data interface{}) (err error) {
	cb := ctx.Query().Get("callback")
	if cb == "" {
		return RenderJSON(ctx, data)
	}

	var buf []byte
	if pbmsg, ok := data.(proto.Message); ok {
		buf, err = MarshalOptions.Marshal(pbmsg)
	} else {
		buf, err = json.Marshal(data)
	}
	if err != nil {
		return err
	}

	callback := template.JSEscapeString(cb)
	w := ctx.Response()
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")

	if _, err = w.Write([]byte(callback)); err != nil {
		return err
	}

	if _, err = w.Write([]byte("(")); err != nil {
		return err
	}

	if _, err = w.Write(buf); err != nil {
		return err
	}

	if _, err = w.Write([]byte(");")); err != nil {
		return err
	}

	return nil
}

func LoginRequired(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)
	w.Write(loginRequired)
}
