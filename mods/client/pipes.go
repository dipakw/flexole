package client

import (
	"net"

	"github.com/dipakw/uconn"
	"github.com/xtaci/smux"
)

func (pp *Pipes) Add(id string, encrypt bool) error {
	conn, err := net.Dial(pp.c.conf.Server.Net, pp.c.conf.Server.Addr)

	if err != nil {
		return err
	}

	// Send the PIPE ID to the server.
	if _, err := conn.Write([]byte(id)); err != nil {
		conn.Close()
		return err
	}

	useConn := conn

	if encrypt {
		useConn, err = uconn.New(conn, &uconn.Opts{
			Algo: pp.c.conf.EncAlgo,
			Key:  pp.c.conf.EncKey,
		})

		if err != nil {
			return err
		}
	}

	sess, err := smux.Client(useConn, smux.DefaultConfig())

	// Set up control channel.
	ctrl, err := pp.c.setupCtrlChan(sess)

	if err != nil {
		return err
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

	pp.c.wg.Add(1)

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
