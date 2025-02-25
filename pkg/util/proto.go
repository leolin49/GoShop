package util

import "google.golang.org/protobuf/proto"

func Serialize(message proto.Message) ([]byte, error) {
	return proto.Marshal(message)
}

// func SerializeM(message []proto.Message) ([]byte, error) {
// 	var data []byte
// 	for _, msg := range message {
// 		d, err := Serialize(msg)
// 		if err != nil {
// 			return nil, err
// 		}
// 		data = append(data, d...)
// 	}
// 	return data, nil
// }

func Deserialize(b []byte, message proto.Message) error {
	return proto.Unmarshal(b, message)
}

// func DeserializeM(b []byte, message []proto.Message) error {
// 	var err error
// 	for len(b) > 0 {
// 		var d []byte
// 		var msg proto.Message
// 		err = Deserialize(d, msg)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }
