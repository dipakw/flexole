package auth

import (
	"time"
)

const (
	CHALLENGE_SIZE   = 40
	MAX_ID_META_SIZE = 287
)

type Auth struct {
	ID   []byte
	Meta map[string]string
	Key  []byte

	err *Err
}

type Err struct {
	reason string
	err    error
}

type ServerOpts struct {
	Timeout     time.Duration
	MaxSigSize  uint16
	MinSigSize  uint16
	DelayOnAuth time.Duration
	VerifySig   func(auth *Auth, msg []byte, sig []byte) (bool, error)
}

type ClientOpts struct {
	ID      []byte
	Meta    map[string]string
	Timeout time.Duration
	SignMsg func(msg []byte) ([]byte, error)
}

type readopts struct {
	full    bool
	timeout time.Duration
}
