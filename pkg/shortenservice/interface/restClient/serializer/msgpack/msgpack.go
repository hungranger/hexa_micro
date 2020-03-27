package msgpack

import (
	"hexa_micro/pkg/shortenservice/model"

	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack"
)

// implerment IRedirectSerializer
type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*model.Redirect, error) {
	redirect := &model.Redirect{}
	if err := msgpack.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}
	return redirect, nil
}

func (r *Redirect) Encode(input *model.Redirect) ([]byte, error) {
	rawMsg, err := msgpack.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}
	return rawMsg, nil
}
