package tcp

import (
	"net"
)

const tcpNetwork = "tcp"

func NewTcpServer(addr string) (net.Listener, error) {
	return net.Listen(tcpNetwork, addr)
}
