package tbus

import (
	"errors"
	"esp/espnet"
	"fmt"
	"logger"
)

// Proxy
const (
	OGS_IDLE = iota
	OGS_REQ_ING
	OGS_REQ_ED
	OGS_REP_ING
)

type CFInfo struct {
	channel espnet.Channel
	method  *ThriftMethodInfo
}

type ChannelFinder func(name string) (*CFInfo, error)

type Proxy struct {
	main        espnet.Channel
	outgo       espnet.Channel
	method      *ThriftMethodInfo
	outgoStatus byte

	finder ChannelFinder
}

func NewProxy(ch espnet.Channel, finder ChannelFinder) *Proxy {
	this := new(Proxy)
	this.main = ch
	this.finder = finder

	cid := this.closerId()
	ch.SetMessageListner(this.MainSendMessage)
	ch.SetCloseListener(cid, this.OnMainClose)

	return this
}

func (this *Proxy) closerId() string {
	return fmt.Sprintf("PR_%p", this)
}

func (this *Proxy) MainSendMessage(msg *espnet.Message) error {
	switch this.outgoStatus {
	case OGS_IDLE:
		return this.FindAndSend(msg)
	case OGS_REP_ING:
		// error!!!
		logger.Error(tag, "receive main message(%s) when out[%s] responsing", msg.Dump(), this.method.Name)
		return this.Kill()
	case OGS_REQ_ED:
		switch this.method.Oneway {
		case B3V_TRUE:
		case B3V_FALSE:
			logger.Error(tag, "receive main message(%s) when out[%s] not response", msg.Dump(), this.method.Name)
			return this.Kill()
		case B3V_UNKNOW:
			logger.Debug(tag, "receive main message but out[%s] no response, oneway method?", this.method.Name)
		}
		return this.FindAndSend(msg)
	case OGS_REQ_ING:
		och := this.outgo
		if och == nil {
			// channel break
			logger.Info(tag, "receiving main message when out[%s] break", this.method.Name)
			return this.Kill()
		}
		seqno, seqmax := espnet.FrameCoders.SeqNO.Get(msg.ToPackage())
		logger.Debug(tag, "send request - %s, %d/%d", this.method.Name, seqno, seqmax)
		if espnet.FrameCoders.SeqNO.IsLast(seqno, seqmax) {
			this.outgoStatus = OGS_REQ_ED
		}
		return och.SendMessage(msg)
	}
	return fmt.Errorf("unknow status %d", this.outgoStatus)
}

func (this *Proxy) OutgoSendMessage(msg *espnet.Message) error {
	mch := this.main
	if mch != nil {
		seqno, seqmax := espnet.FrameCoders.SeqNO.Get(msg.ToPackage())
		logger.Debug(tag, "receive response - %s, %d/%d", this.method.Name, seqno, seqmax)
		if espnet.FrameCoders.SeqNO.IsLast(seqno, seqmax) {
			this.outgoStatus = OGS_IDLE
			och := this.outgo
			this.outgo = nil
			och.SetMessageListner(nil)
			och.AskClose()
		} else {
			this.outgoStatus = OGS_REP_ING
		}
		return mch.SendMessage(msg)
	}
	return errors.New("closed")
}

func (this *Proxy) doClose(ch espnet.Channel, out bool) {
	if ch != nil {
		ch.SetMessageListner(nil)
		ch.SetCloseListener(this.closerId(), nil)
		if out {
			if this.outgoStatus == OGS_IDLE {
				ch.AskClose()
			} else {
				espnet.CloseForce(ch)
			}
		} else {
			ch.AskClose()
		}
	}
}

func (this *Proxy) OnMainClose() {
	mch := this.main
	this.main = nil
	if mch != nil {
		mch.SetMessageListner(nil)
	}

	och := this.outgo
	this.outgo = nil
	this.doClose(och, true)
}

func (this *Proxy) Kill() error {
	mch := this.main
	this.main = nil

	och := this.outgo
	this.outgo = nil

	if logger.EnableDebug(tag) {
		os := ""
		if och != nil {
			os = och.String()
		}
		logger.Debug(tag, "killing %s -> %s", mch, os)
	}

	this.doClose(mch, false)
	this.doClose(och, true)
	return nil
}

func (this *Proxy) FindAndSend(msg *espnet.Message) error {
	addr := msg.GetAddress()
	if addr == nil {
		err := logger.Error(tag, "message address nil")
		this.Kill()
		return err
	}
	a := addr.Identity()
	cfi, err := this.finder(a)
	if err != nil {
		oerr := logger.Error(tag, "dispatch '%s' fail - %s", a, err)
		// write texception
		bs := SerializeReplyException(msg, oerr.Error())
		if bs != nil {
			rmsg := msg.ReplyMessage()
			rmsg.SetPayload(bs)
			espnet.CloseAfterSend(rmsg)
			this.main.SendMessage(rmsg)
		} else {
			this.Kill()
		}
		return err
	}

	this.method = cfi.method
	this.outgo = cfi.channel
	this.outgo.SetMessageListner(this.OutgoSendMessage)
	seqno, seqmax := espnet.FrameCoders.SeqNO.Get(msg.ToPackage())
	logger.Debug(tag, "send request - %s, %d/%d", this.method.Name, seqno, seqmax)
	if espnet.FrameCoders.SeqNO.IsLast(seqno, seqmax) {
		this.outgoStatus = OGS_REQ_ED
	} else {
		this.outgoStatus = OGS_REQ_ING
	}
	logger.Debug(tag, "dispatch '%s' -> %s", a, this.outgo)

	return this.outgo.SendMessage(msg)
}
