package server

import (
	"context"
	"flexole/mods/cmd"
	"net"

	"github.com/xtaci/smux"
)

func (s *Server) handle(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	// Authenticate.
	auth, err := s.conf.AuthFN(conn)

	if err != nil {
		s.conf.Log.Errf("Authentication failed: (%s) : %s", conn.RemoteAddr(), err.Error())
		return
	}

	// Set up encryption if needed.
	useConn := conn

	/*if auth.Encrypt {
		algo, key, err := s.conf.EncFN(auth, conn)

		if err != nil {
			fmt.Println("Failed to get enc info:", err)
			return
		}

		useConn, err = uconn.New(conn, &uconn.Opts{
			Algo: algo,
			Key:  key,
		})

		if err != nil {
			fmt.Println("Failed to create enc conn:", err)
			return
		}
	}*/

	// Init mux.
	sess, err := smux.Server(useConn, smux.DefaultConfig())

	if err != nil {
		s.conf.Log.Errf("Failed to create mux session: (%s) : %s", conn.RemoteAddr(), err.Error())
		return
	}

	// Open control stream.
	ctrl, err := sess.OpenStream()

	if err != nil {
		s.conf.Log.Errf("Failed to open control stream: (%s) : %s", conn.RemoteAddr(), err.Error())
		return
	}

	// Send control command.
	if _, err := ctrl.Write(cmd.New(cmd.CMD_OPEN_CTRL_CHAN, nil).Pack()); err != nil {
		s.conf.Log.Errf("Failed to send control chan command: (%s) : %s", conn.RemoteAddr(), err.Error())
		return
	}

	// Add pipe.
	pipe := s.User(auth.UserID).pipes.add(&Pipe{
		userID: auth.UserID,
		id:     auth.PipeID,
		active: true,
		conn:   conn,
		sess:   sess,
		ctrl:   ctrl,
	})

	// Listen for control commands.
	go s.listenCtrl(ctx, pipe)
	defer s.User(auth.UserID).pipes.rem(auth.PipeID)

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
