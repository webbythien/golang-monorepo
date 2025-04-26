package server

import "errors"

var (
	ErrZeroLengthPayload = errors.New("zero-length payload is not a valid JSON object")
	ErrUsesOldProtobuf   = errors.New("%T uses github.com/golang/protobuf, but connect-go only supports google.golang.org/protobuf: see https://go.dev/blog/protobuf-apiv2")
	ErrNotProtoMessage   = errors.New("%T doesn't implement proto.Message")
)
