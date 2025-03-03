package util

import (
	"crypto/sha256"
	"encoding/hex"

	"google.golang.org/protobuf/proto"
)

func Serialize(message proto.Message) ([]byte, error) {
	return proto.Marshal(message)
}

func Deserialize(b []byte, message proto.Message) error {
	return proto.Unmarshal(b, message)
}

func GenerateProtoHash(message proto.Message) (string, error) {
	msgBytes, err := Serialize(message)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(msgBytes)
	return hex.EncodeToString(hash[:]), nil
}
