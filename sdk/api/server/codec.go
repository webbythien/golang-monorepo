package server

import (
	"bytes"
	"encoding/json"
	"fmt"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoiface"
)

func DefaultProtoJSONCodec() connect.Codec {
	return &protoJSONCodec{name: "json"}
}

type protoJSONCodec struct {
	name string
}

var _ connect.Codec = (*protoJSONCodec)(nil)

var protojsonMarshalOptions = protojson.MarshalOptions{
	EmitUnpopulated:   true,
	EmitDefaultValues: true,
}

var protojsonUnmarshalOptions = protojson.UnmarshalOptions{
	DiscardUnknown: true,
}

func (c *protoJSONCodec) Name() string { return c.name }

func (c *protoJSONCodec) Marshal(message any) ([]byte, error) {
	protoMessage, ok := message.(proto.Message)
	if !ok {
		return nil, errNotProto(message)
	}
	return protojsonMarshalOptions.Marshal(protoMessage)
}

func (c *protoJSONCodec) MarshalAppend(dst []byte, message any) ([]byte, error) {
	protoMessage, ok := message.(proto.Message)
	if !ok {
		return nil, errNotProto(message)
	}
	return protojsonMarshalOptions.MarshalAppend(dst, protoMessage)
}

func (c *protoJSONCodec) Unmarshal(binary []byte, message any) error {
	protoMessage, ok := message.(proto.Message)
	if !ok {
		return errNotProto(message)
	}
	if len(binary) == 0 {
		return ErrZeroLengthPayload
	}
	// Discard unknown fields so clients and servers aren't forced to always use
	// exactly the same version of the schema.
	options := protojsonUnmarshalOptions
	err := options.Unmarshal(binary, protoMessage)
	if err != nil {
		return fmt.Errorf("unmarshal into %T: %w", message, err)
	}
	return nil
}

func (c *protoJSONCodec) MarshalStable(message any) ([]byte, error) {
	// protojson does not offer a "deterministic" field ordering, but fields
	// are still ordered consistently by their index. However, protojson can
	// output inconsistent whitespace for some reason, therefore it is
	// suggested to use a formatter to ensure consistent formatting.
	// https://github.com/golang/protobuf/issues/1373
	messageJSON, err := c.Marshal(message)
	if err != nil {
		return nil, err
	}
	compactedJSON := bytes.NewBuffer(messageJSON[:0])
	if err = json.Compact(compactedJSON, messageJSON); err != nil {
		return nil, err
	}
	return compactedJSON.Bytes(), nil
}

func (c *protoJSONCodec) IsBinary() bool {
	return false
}

func errNotProto(message any) error {
	if _, ok := message.(protoiface.MessageV1); ok {
		return fmt.Errorf("%w: %v", ErrUsesOldProtobuf, message)
	}
	return fmt.Errorf("%w: %v", ErrNotProtoMessage, message)
}
