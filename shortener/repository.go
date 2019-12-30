package shortener

//RedirectRespository interface for hitting db
type RedirectRespository interface {
	Find(code string) (*Redirect, error)
	Store(redirect *Redirect) error
}
