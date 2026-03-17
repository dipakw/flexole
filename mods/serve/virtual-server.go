package serve

import (
	"net"
	"net/http"
)

type singleConnListener struct {
	conn net.Conn
	done bool
}

func (l *singleConnListener) Accept() (net.Conn, error) {
	if l.done {
		return nil, net.ErrClosed
	}
	l.done = true
	return l.conn, nil
}

func (l *singleConnListener) Close() error {
	return nil
}

func (l *singleConnListener) Addr() net.Addr {
	return l.conn.LocalAddr()
}

type virtualServer struct {
	server *http.Server
}

func (s *Serve) newVirtualServer() *virtualServer {
	handler := http.FileServer(http.FS(s.fs))

	return &virtualServer{
		server: &http.Server{
			Handler: handler,
		},
	}
}

func (v *virtualServer) serve(conn net.Conn) {
	v.server.Serve(&singleConnListener{conn: conn})
}
