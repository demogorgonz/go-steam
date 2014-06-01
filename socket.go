package steam

import (
	"errors"
	"net"

	"github.com/golang/glog"
)

type socket struct {
	conn  *net.UDPConn
	raddr *net.UDPAddr
}

func newSocket(addr string) (*socket, error) {
	raddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return nil, err
	}

	return &socket{conn, raddr}, nil
}

func (s *socket) close() {
	s.conn.Close()
}

func (s *socket) send(payload []byte) error {
	glog.V(1).Infof("steam: sending %v bytes payload to %v", len(payload), s.raddr)
	glog.V(2).Infof("steam: sending payload to %v: %X", s.raddr, payload)
	n, err := s.conn.WriteToUDP(payload, s.raddr)
	if err != nil {
		return err
	}
	if n != len(payload) {
		return errors.New("steam: could not send full request to server")
	}

	return nil
}

func (s *socket) receivePacket() ([]byte, error) {
	var buf [1500]byte
	n, _, err := s.conn.ReadFromUDP(buf[:])
	if err != nil {
		return nil, err
	}
	glog.V(1).Infof("steam: received %v bytes from %v", n, s.raddr)
	glog.V(2).Infof("steam: received payload %v: %X", s.raddr, buf[:n])

	return buf[:n], nil
}

func (s *socket) receive() ([]byte, error) {
	buf, err := s.receivePacket()
	if err != nil {
		return nil, err
	}
	if buf[0] == 0xFE {
		return nil, errors.New("steam: cannot handle split packets")
	}
	return buf[4:], nil
}