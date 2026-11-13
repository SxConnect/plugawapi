package main

import (
	"context"
	"log/slog"

	"github.com/emiago/sipgo"
	"github.com/emiago/sipgo/sip"
)

type sipGateway struct {
	ua     *sipgo.UserAgent
	server *sipgo.Server
	log    *slog.Logger
	addr   string
}

func newSipGateway(log *slog.Logger, addr string) (*sipGateway, error) {
	ua, err := sipgo.NewUA(sipgo.WithUserAgent("PlugaWAPI"))
	if err != nil {
		return nil, err
	}
	srv, err := sipgo.NewServer(ua)
	if err != nil {
		return nil, err
	}
	gw := &sipGateway{ua: ua, server: srv, log: log, addr: addr}
	// handle INVITE
	srv.OnInvite(func(req *sip.Request, tx sip.ServerTransaction) {
		log.Info("SIP INVITE received", "from", req.From().String())
		// respond 200 OK (placeholder)
		res := sip.NewResponseFromRequest(req, 200, "OK", nil)
		tx.Respond(res)
	})
	srv.OnRegister(func(req *sip.Request, tx sip.ServerTransaction) {
		log.Info("SIP REGISTER received", "from", req.From().String())
		res := sip.NewResponseFromRequest(req, 200, "OK", nil)
		tx.Respond(res)
	})
	return gw, nil
}

func (g *sipGateway) start(ctx context.Context) error {
	return g.server.ListenAndServe(ctx, "udp", g.addr)
}
