package client

import (
	"net"

	"github.com/dipakw/byrate/handle"
	"github.com/xtaci/smux"
)

var byrateConfig = &handle.Config{
	Version: "0.2.0",
}

func (c *Client) vSpeed(stream *smux.Stream) {
	conn1, conn2 := net.Pipe()

	defer conn1.Close()
	defer conn2.Close()

	go handle.Handle(conn1, byrateConfig)

	relay(c.ctx, conn2, stream)
}
