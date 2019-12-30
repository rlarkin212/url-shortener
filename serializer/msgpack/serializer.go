package msgpack

import (
	"github.com/pkg/errors"
	"github.com/rlarkin212/url-shortener/shortener"
	"github.com/vmihailenco/msgpack"
)

//Redirect struct from shortener
type Redirect struct{}

//Decode incoming request
func (r *Redirect) Decode(input []byte) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}

	err := msgpack.Unmarshal(input, redirect)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}

	return redirect, nil
}

//Encode incoming request
func (r *Redirect) Encode(input *shortener.Redirect) ([]byte, error) {
	rawMsg, err := msgpack.Marshal(input)

	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}

	return rawMsg, nil
}
