package service

import (
	"github.com/artm2000/urlbook/internal/core/entity/code"
	"github.com/artm2000/urlbook/internal/core/entity/message"
	"github.com/artm2000/urlbook/internal/core/model/response"
)

func createSuccessResponse[T any](data *T, trackId string) *response.Response[T] {
	return &response.Response[T]{
		Error:      false,
		Message:    "",
		StatusCode: code.SUCCESS,
		Data:       data,
		TrackId:    trackId,
	}
}

func createFailResponse[T any](message message.Message, trackId string, statusCode int) *response.Response[T] {
	return &response.Response[T]{
		Error:      true,
		Message:    message,
		StatusCode: statusCode,
		Data:       nil,
		TrackId:    trackId,
	}
}
