package remaps

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"

	util "github.com/koangel/grapeNet/Utils"
)

func Lookup(host string) (string, error) {
	addrs, err := net.LookupHost(host)
	if err != nil {
		return "", err
	}
	if len(addrs) < 1 {
		return "", errors.New("unknown host")
	}
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return addrs[rd.Intn(len(addrs))], nil
}

var Data = []byte("abcdefghijklmnopqrstuvwabcdefghi")

const MaxPing = 99999

type Reply struct {
	Src   string
	Addr  string
	Time  int64
	TTL   uint8
	Error error
}

func MarshalMsg(req int, data []byte) ([]byte, error) {
	xid, xseq := os.Getpid()&0xffff, req
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: xid, Seq: xseq,
			Data: data,
		},
	}
	return wm.Marshal(nil)
}

type ping struct {
	Src  string
	Addr string
	Timeout int
	Conn net.Conn
	Data []byte
}

func (self *ping) Dail() (err error) {
	self.Conn, err = net.Dial("ip4:icmp", self.Addr)
	if err != nil {
		return err
	}
	return nil
}

func (self *ping) SetDeadline(timeout int) error {
	return self.Conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
}

func (self *ping) Close() error {
	return self.Conn.Close()
}

func (self *ping) Ping(count int) {
	if err := self.Dail(); err != nil {
		fmt.Println("Not found remote host")
		return
	}
	fmt.Println("Start ping from ", self.Conn.LocalAddr())
	self.SetDeadline(self.Timeout)
	for i := 0; i < count; i++ {
		r := sendPingMsg(self.Src, self.Addr, self.Conn, self.Data)
		if r.Error != nil {
			if opt, ok := r.Error.(*net.OpError); ok && opt.Timeout() {
				fmt.Printf("From %s reply: TimeOut\n", self.Addr)
				if err := self.Dail(); err != nil {
					fmt.Println("Not found remote host")
					return
				}
			} else {
				fmt.Printf("From %s reply: %s\n", self.Addr, r.Error)
			}
		} else {
			fmt.Printf("From %s reply: time=%d ttl=%d\n", self.Addr, r.Time, r.TTL)
		}
		//time.Sleep(1e9)
	}
}

func (self *ping) PingCount(count int) (reply []Reply) {
	if err := self.Dail(); err != nil {
		fmt.Println("Not found remote host")
		return
	}
	self.SetDeadline(self.Timeout)
	for i := 0; i < count; i++ {
		r := sendPingMsg(self.Src, self.Addr, self.Conn, self.Data)
		reply = append(reply, r)
		//time.Sleep(1e9)
	}
	return
}

func Run(addr string, req int, data []byte) (*ping, error) {
	wb, err := MarshalMsg(req, data)
	if err != nil {
		return nil, err
	}
	src := addr
	addr, err = Lookup(addr)
	if err != nil {
		return nil, err
	}
	return &ping{Data: wb, Addr: addr, Src: src}, nil
}

func PingOnce(addr string, timeout int) (reply Reply) {
	reply.Src = addr
	reply.Time = MaxPing

	ping, err := Run(addr, 8, []byte("t"))
	if err != nil {
		reply.Error = err
		return
	}
	defer ping.Close()
	ping.Timeout = timeout // 一定要设置超时
	r := ping.PingCount(1)
	if len(r) == 0 {
		reply.Error = fmt.Errorf("Not found remote host")
		return
	}

	return r[0]
}

func PingArray(addrs ...string) (reply []Reply) {
	addrArr := []string(addrs)

	jobs := util.SyncJob{}

	for _, v := range addrArr {
		jobs.Append(func(addr string) {
			r := PingOnce(addr, 1)
			reply = append(reply, r)
		}, v)
	}

	jobs.StartWait()
	return
}

func sendPingMsg(src, addr string, c net.Conn, wb []byte) (reply Reply) {
	start := time.Now()

	if _, reply.Error = c.Write(wb); reply.Error != nil {
		return
	}

	rb := make([]byte, 1500)
	var n int
	n, reply.Error = c.Read(rb)
	if reply.Error != nil {
		return
	}

	duration := time.Now().Sub(start)
	ttl := uint8(rb[8])
	rb = func(b []byte) []byte {
		if len(b) < 20 {
			return b
		}
		hdrlen := int(b[0]&0x0f) << 2
		return b[hdrlen:]
	}(rb)
	var rm *icmp.Message
	rm, reply.Error = icmp.ParseMessage(1, rb[:n])
	if reply.Error != nil {
		return
	}

	reply.Time = MaxPing
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		t := int64(duration / time.Millisecond)
		reply = Reply{src, addr, t, ttl, nil}
	case ipv4.ICMPTypeDestinationUnreachable:
		reply.Error = errors.New("Destination Unreachable")
	default:
		reply.Error = fmt.Errorf("Not ICMPTypeEchoReply %v", rm)
	}
	return
}
