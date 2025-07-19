package server

import (
	"flexole/mods/cmd"
	"fmt"
	"net"

	"github.com/dipakw/uconn"
	"github.com/xtaci/smux"
)

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	// Authenticate.
	auth, err := s.conf.AuthFN(conn)

	if err != nil {
		fmt.Println("Auth failed:", err)
		return
	}

	// Set up encryption if needed.
	useConn := conn

	if auth.Encrypt {
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
	}

	// Init mux.
	sess, err := smux.Server(useConn, smux.DefaultConfig())

	if err != nil {
		fmt.Println("Failed to create mux session:", err)
		return
	}

	// Open control stream.
	ctrl, err := sess.OpenStream()

	if err != nil {
		fmt.Println("Failed to open control stream:", err)
		return
	}

	// Send control command.
	if _, err := ctrl.Write(cmd.New(cmd.CMD_OPEN_CTRL_CHAN, nil).Pack()); err != nil {
		fmt.Println("Failed to send control chan command:", err)
		return
	}

	// Add pipe.
	pipe := s.AddPipe(&Pipe{
		userID: auth.UserID,
		id:     auth.PipeID,
		active: true,
		conn:   conn,
		sess:   sess,
		ctrl:   ctrl,
	})

	// Listen for control commands.
	go s.listenCtrl(pipe)

	defer s.RemPipe(auth.UserID, auth.PipeID)
	defer sess.Close()

	// To detect broken pipes.
	for {
		_, err := sess.AcceptStream()

		if err != nil {
			break
		}
	}
}
