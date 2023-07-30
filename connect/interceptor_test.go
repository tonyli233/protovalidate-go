package connect

import (
	"context"
	"testing"

	"github.com/bufbuild/connect-go"
	examplev1 "github.com/bufbuild/protovalidate-go/internal/gen/tests/example/v1"
	"github.com/stretchr/testify/assert"
)

func Test_interceptor_WrapUnary(t *testing.T) {
	t.Parallel()

	interceptor, err := New()
	assert.NoError(t, err)
	noopUnaryFunc := func(context context.Context, request connect.AnyRequest) (connect.AnyResponse, error) {
		return nil, nil
	}
	coordinates := &examplev1.Coordinates{
		Lat: -100,
		Lng: 200,
	}
	wrappedUnaryFunc := interceptor.WrapUnary(noopUnaryFunc)
	request := connect.NewRequest(coordinates)
	_, err = wrappedUnaryFunc(context.Background(), request)
	assert.Error(t, err)

	coordinates = &examplev1.Coordinates{
		Lat: 0,
		Lng: 0,
	}
	request = connect.NewRequest(coordinates)
	wrappedUnaryFunc = interceptor.WrapUnary(noopUnaryFunc)
	_, err = wrappedUnaryFunc(context.Background(), request)
	assert.NoError(t, err)
}
