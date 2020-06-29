package server

import (
	"encoding/json"
	"strings"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

var (
	marshaler   = jsonpb.Marshaler{}
	unMarshaler = jsonpb.Unmarshaler{}
)

func (c *context) Pb(code int, v proto.Message) error {
	c.Response().Header().Set(HeaderContentType, MIMEApplicationProtobuf)
	c.Response().WriteHeader(code)

	var b []byte
	var err error
	if strings.HasPrefix(c.Request().Header.Get("Accept"), "application/json") {
		b, err = json.Marshal(v)
	} else {
		b, err = proto.Marshal(v)
	}

	if err != nil {
		return err
	}
	// TODO 修改源代码,以适应接口缓存业务需求,对原来框架无影响,选择这里的主要是因为protobuf返回方式是调用的这里
	c.Set(ctxWriteBackKey, b)
	_, err = c.Response().Write(b)
	return err
}
