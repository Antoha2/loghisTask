package transport

import (
	"net/http"

	"github.com/antoha2/loghis/service"
)

type MQ interface {
	Init() error
	//Write(ctx context.Context, log string) error
}

type MQImpl struct {
	service service.LogServiceImpl
	server  *http.Server
}

func NewMQ(s *service.LogServiceImpl) *MQImpl {

	return &MQImpl{service: *s}
}

type LoggerMsg struct {
	Log string `json:"log"`
}
