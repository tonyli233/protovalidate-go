package connect

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/protobuf/proto"
)

func New(options ...InterceptorOption) (connect.Interceptor, error) {
	cfg := config{}
	for _, opt := range options {
		opt(&cfg)
	}
	if cfg.validator == nil {
		validator, err := protovalidate.New()
		if err != nil {
			return nil, err
		}
		cfg.validator = validator
	}
	return interceptor{
		cfg: cfg,
	}, nil
}

type interceptor struct {
	cfg config
}

func (i interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
		message, ok := request.Any().(proto.Message)
		if !ok {
			return nil, errors.New("invalid request message")
		}
		if err := i.cfg.validator.Validate(message); err != nil {
			return nil, err
		}
		return next(ctx, request)
	}
}

func (i interceptor) WrapStreamingClient(clientFunc connect.StreamingClientFunc) connect.StreamingClientFunc {
	return clientFunc
}

func (i interceptor) WrapStreamingHandler(handlerFunc connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return handlerFunc
}

type config struct {
	validator *protovalidate.Validator
}

type InterceptorOption func(*config)

func WithValidator(validator *protovalidate.Validator) InterceptorOption {
	return func(c *config) {
		c.validator = validator
	}
}
