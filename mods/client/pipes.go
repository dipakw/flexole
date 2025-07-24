package client

import (
	"crypto/sha256"
	"flexole/mods/auth"
	"fmt"
	"net"
	"time"

	"github.com/dipakw/uconn"
	"github.com/xtaci/smux"
)

func (pp *Pipes) Add(id string, encrypt bool) error {
	conn, err := net.Dial(pp.c.conf.Server.Net, pp.c.conf.Server.Addr)

	if err != nil {
		return err
	}

	pp.c.conf.Log.Inff("Adding pipe => id: %s | encrypted: %v", id, encrypt)

	// Authenticate.
	auth := auth.Client(conn, &auth.ClientOpts{
		ID:      pp.c.conf.ID,
		Timeout: 10 * time.Second,

		Meta: map[string]string{
			"pipe": id,
			"enc":  map[bool]string{true: "1", false: "0"}[encrypt],
		},

		SignMsg: func(msg []byte) ([]byte, error) {
			data := make([]byte, len(pp.c.conf.ID)+len(msg)+len(pp.c.conf.Key))
			copy(data, pp.c.conf.ID)
			copy(data[len(pp.c.conf.ID):], msg)
			copy(data[len(pp.c.conf.ID)+len(msg):], pp.c.conf.Key)

			// Generate sha256 hash.
			hash := sha256.New()
			hash.Write(data)

			// Return the signature.
			return hash.Sum(nil), nil
		},
	})

	if !auth.Ok() {
		return fmt.Errorf("Authentication error: %s : %s", auth.Err().Reason(), auth.Err().Main().Error())
	}

	useConn := conn

	if encrypt {
		useConn, err = uconn.New(conn, &uconn.Opts{
			Algo: uconn.ALGO_AES256_GCM,
			Key:  auth.Key,
		})

		if err != nil {
			return fmt.Errorf("Failed to create encrypted connection => pipe: %s | error: %s", id, err.Error())
		}
	}

	// Send OK
	if _, err := useConn.Write([]byte("OK")); err != nil {
		return fmt.Errorf("Failed to send OK => pipe: %s | error: %s", id, err.Error())
	}

	sess, err := smux.Client(useConn, smux.DefaultConfig())

	// Set up control channel.
	ctrl, err := pp.c.setupCtrlChan(sess)

	if err != nil {
		return fmt.Errorf("Failed to setup control channel => pipe: %s | error: %s", id, err.Error())
	}

	pp.c.mu.Lock()
	defer pp.c.mu.Unlock()

	pp.c.pipesList[id] = &connPipe{
		id:     id,
		active: true,
		conn:   useConn,
		sess:   sess,
		ctrl:   ctrl,
	}

	pp.c.conf.Log.Inff("Pipe added => id: %s | encrypted: %v", id, encrypt)

	go pp.c.listen(id)

	return nil
}

func (pp *Pipes) Rem(id string) error {
	pp.c.mu.Lock()
	defer pp.c.mu.Unlock()

	if pipe, ok := pp.c.pipesList[id]; ok {
		pipe.conn.Close()
		delete(pp.c.pipesList, id)
	}

	return nil
}
