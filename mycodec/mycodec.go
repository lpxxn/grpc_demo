package mycodec

import (
	"encoding/json"
	"fmt"

	"google.golang.org/grpc/encoding"
)

const Name = "myCodec"

type MyCodec struct {
}

func init() {
	encoding.RegisterCodec(MyCodec{})
}

func (MyCodec) Marshal(v interface{}) ([]byte, error) {
	fmt.Println("myCodec Marshal")
	return json.Marshal(v)
}

func (MyCodec) Unmarshal(data []byte, v interface{}) error {
	fmt.Println("myCodec Unmarshal")
	return json.Unmarshal(data, v)
}

func (MyCodec) Name() string {
	return Name
}
