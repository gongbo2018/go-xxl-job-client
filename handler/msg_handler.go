package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/dubbogo/getty"

	"github.com/feixiaobo/go-xxl-job-client/v2/transport"
)

const (
	cronTime         = 20e9
	writePkg_Timeout = 5 * time.Second
)

type MessageHandler struct {
	GettyClient *transport.GettyRPCClient
	MsgHandle   func(pkg interface{}) (res []byte, err error)
}

func (h *MessageHandler) OnOpen(session getty.Session) error {
	log.Print("OnOpen session:", session.Stat())
	h.GettyClient.AddSession(session)
	return nil
}

func (h *MessageHandler) OnError(session getty.Session, err error) {
	log.Printf("OnError session{%s} got error{%v}, will be closed.", session.Stat(), err)
}

func (h *MessageHandler) OnClose(session getty.Session) {
	log.Printf("OnClose session{%s} is closing......", session.Stat())
	h.GettyClient.RemoveSession(session)
}

func (h *MessageHandler) OnMessage(session getty.Session, pkg interface{}) {
	s, ok := pkg.([]interface{})
	if !ok {
		log.Printf("illegal package{%#v}", pkg)
		return
	}

	for _, v := range s {
		if v != nil {
			res, err := h.MsgHandle(v)
			reply(session, res, err)
		}
	}
}

func (h *MessageHandler) OnCron(session getty.Session) {
	active := session.GetActive()
	if cronTime < time.Since(active).Nanoseconds() {
		log.Printf("OnCorn session{%s} timeout{%s}", session.Stat(), time.Since(active).String())
		session.Close()
		h.GettyClient.RemoveSession(session)
	}
}

func reply(session getty.Session, resBy []byte, err error) {
	pkg := transport.NewHttpResponsePkg(http.StatusOK, resBy)
	if err != nil || resBy == nil {
		pkg = transport.NewHttpResponsePkg(http.StatusInternalServerError, resBy)
	}

	if err := session.WritePkg(pkg, writePkg_Timeout); err != nil {
		log.Printf("WritePkg error: %#v, %#v", pkg, err)
	}
}
