package serializer

import "hexa_micro/pkg/shortenservice/model"

type IRedirectSerializer interface {
	Decode(input []byte) (*model.Redirect, error)
	Encode(input *model.Redirect) ([]byte, error)
}
