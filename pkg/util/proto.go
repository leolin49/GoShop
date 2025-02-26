package util

import "google.golang.org/protobuf/proto"

func Serialize(message proto.Message) ([]byte, error) {
	return proto.Marshal(message)
}

func Deserialize(b []byte, message proto.Message) error {
	return proto.Unmarshal(b, message)
}
