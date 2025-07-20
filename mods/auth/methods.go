package auth

func (a *Auth) Encrypt(b []byte) ([]byte, error) {
	return encrypt(a.Key, b)
}

func (a *Auth) Decrypt(b []byte) ([]byte, error) {
	return decrypt(a.Key, b)
}

func (a *Auth) Err() *Err {
	return a.err
}

func (a *Auth) Ok() bool {
	return a.err == nil
}

func (a *Auth) re(err *Err) *Auth {
	a.err = err
	return a
}

func (e *Err) Reason() string {
	return e.reason
}

func (err *Err) Main() error {
	return err.err
}
