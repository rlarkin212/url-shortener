package shortener

import (
	"errors"
	"time"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

var (
	//ErrRedirectNotFound not found error
	ErrRedirectNotFound = errors.New("Redirect Not Found")

	//ErrRedirectInvalid invalid error
	ErrRedirectInvalid = errors.New("Redirect Invalid")
)

type redirectService struct {
	redirectRepo RedirectRespository
}

//NewRedirectRespository new instance of RedirectRespository
func NewRedirectRespository(redirectRepo RedirectRespository) RedirectService {
	return &redirectService{
		redirectRepo,
	}
}

func (r *redirectService) Find(code string) (*Redirect, error) {
	return r.redirectRepo.Find(code)
}

func (r *redirectService) Store(redirect *Redirect) error {

	err := validate.Validate(redirect)

	if err != nil {
		return errs.Wrap(ErrRedirectInvalid, "service.Redirect.Store")
	}

	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC().Unix()

	result := r.redirectRepo.Store(redirect)

	return result
}
