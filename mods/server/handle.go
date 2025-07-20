package server

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"flexole/mods/auth"
	"flexole/mods/cmd"
	"net"
	"time"

	"github.com/dipakw/uconn"
	"github.com/xtaci/smux"
)

func (s *Server) handle(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	// Authenticate.
	auth := auth.Server(conn, &auth.ServerOpts{
		Timeout:    10 * time.Second,
		MaxSigSize: 60,
		MinSigSize: 60,

		VerifySig: func(a *auth.Auth, msg []byte, sig []byte) (bool, error) {
			key, err := s.conf.KeyFN(a.ID)

			if err != nil {
				return false, err
			}

			// Generate sha256 hash.
			data := make([]byte, len(a.ID)+len(msg)+len(key))
			copy(data, a.ID)
			copy(data[len(a.ID):], msg)
			copy(data[len(a.ID)+len(msg):], key)

			// Generate sha256 hash.
			hash := sha256.New()
			hash.Write(data)

			// Compare with the signature.
			return hmac.Equal(hash.Sum(nil), sig), nil
		},
	})

	if !auth.Ok() {
		s.conf.Log.Errf("Authentication failed:: addr: %s | reason: %s | error: %s", conn.RemoteAddr(), auth.Err().Reason(), auth.Err().Main())
		return
	}

	userId := string(auth.ID)
	pipeId := auth.Meta["pipe"]
	encrypt := auth.Meta["enc"] == "1"

	var err error

	// Set up encryption if needed.
	useConn := conn

	if encrypt {
		useConn, err = uconn.New(conn, &uconn.Opts{
			Algo: uconn.ALGO_AES256_GCM,
			Key:  auth.Key,
		})

		if err != nil {
			s.conf.Log.Errf("Failed to set up encryption:: user: %s | addr: %s | error: %s", userId, conn.RemoteAddr(), err.Error())
			return
		}
	}

	// Init mux.
	sess, err := smux.Server(useConn, smux.DefaultConfig())

	if err != nil {
		s.conf.Log.Errf("Failed to create mux session:: addr: %s | error: %s", conn.RemoteAddr(), err.Error())
		return
	}

	// Open control stream.
	ctrl, err := sess.OpenStream()

	if err != nil {
		s.conf.Log.Errf("Failed to open control stream:: addr: %s | error: %s", conn.RemoteAddr(), err.Error())
		return
	}

	// Send control command.
	if _, err := ctrl.Write(cmd.New(cmd.CMD_OPEN_CTRL_CHAN, nil).Pack()); err != nil {
		s.conf.Log.Errf("Failed to send control chan command:: addr: %s | error: %s", conn.RemoteAddr(), err.Error())
		return
	}

	// Add pipe.
	pipe := s.User(userId).pipes.add(&Pipe{
		userID: userId,
		id:     pipeId,
		conn:   conn,
		sess:   sess,
		ctrl:   ctrl,
	})

	// Listen for control commands.
	go s.listenCtrl(ctx, pipe)
	defer s.User(userId).pipes.rem(pipeId)

	// To detect broken pipes.
	for {
		select {
		case <-ctx.Done():
			return
		default:
			_, err := sess.AcceptStream()

			if err != nil {
				return
			}
		}
	}
}
