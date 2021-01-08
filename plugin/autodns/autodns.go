package autodns

import (
	"context"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"net"
	"strconv"
)

type AutoDNS struct{}

func (ad AutoDNS) Name() string {
	return "autodns"
}

func (ad AutoDNS) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	// get request information from parameters
	state := request.Request{W: w, Req: r}

	// setup reply message
	replyMsg := &dns.Msg{}
	replyMsg.SetReply(r)
	replyMsg.Authoritative = true

	ip := state.IP()
	var rr dns.RR

	// write response msg to client
	switch state.Family() {
	case 1:
		rr = &dns.A{}
		rr.(*dns.A).Hdr = dns.RR_Header{
			Name:   state.QName(),
			Rrtype: dns.TypeA,
			Class:  state.QClass(),
		}
	case 2:
		rr = &dns.AAAA{}
		rr.(*dns.AAAA).Hdr = dns.RR_Header{
			Name:   state.QName(),
			Rrtype: dns.TypeAAAA,
			Class:  state.QClass(),
		}
		rr.(*dns.AAAA).AAAA = net.ParseIP(ip)
	}

	srv := &dns.SRV{}
	srv.Hdr = dns.RR_Header{
		Name:   "_" + state.Proto() + "." + state.QName(),
		Rrtype: dns.TypeSRV,
		Class:  state.QClass(),
	}
	port, _ := strconv.Atoi(state.Port())
	srv.Port = uint16(port)
	srv.Target = "."

	replyMsg.Extra = []dns.RR{rr, srv}
	w.WriteMsg(replyMsg)
	return 0, nil
}
