package service

import "context"

type LogService interface {
	Write(ctx context.Context, msg string) error
}

type LogServiceImpl struct {
	LogService
}

func NewLogService() *LogServiceImpl {
	return &LogServiceImpl{}
}
