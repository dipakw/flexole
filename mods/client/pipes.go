package client

import (
	"net"

	"github.com/dipakw/uconn"
	"github.com/xtaci/smux"
)

func (p *Pipes) Add(id string, encrypt bool) error {
	conn, err := net.Dial(p.c.conf.Server.Net, p.c.conf.Server.Addr)

	if err != nil {
		return err
	}

	useConn := conn

	if encrypt {
		useConn, err = uconn.New(conn, &uconn.Opts{
			Algo: p.c.conf.EncAlgo,
			Key:  p.c.conf.EncKey,
		})

		if err != nil {
			return err
		}
	}

	sess, err := smux.Client(useConn, smux.DefaultConfig())

	// Set up control channel.
	ctrl, err := p.c.setupCtrlChan(sess)

	if err != nil {
		return err
	}

	p.c.mu.Lock()
	defer p.c.mu.Unlock()

	p.c.pipes[id] = &aPipe{
		id:     id,
		active: true,
		conn:   useConn,
		sess:   sess,
		ctrl:   ctrl,
	}

	p.c.wg.Add(1)

	go p.c.listen(id)

	return nil
}

func (p *Pipes) Rem(id string) error {
	p.c.mu.Lock()
	defer p.c.mu.Unlock()

	if pipe, ok := p.c.pipes[id]; ok {
		pipe.sess.Close()
		pipe.conn.Close()

		delete(p.c.pipes, id)
	}

	return nil
}
