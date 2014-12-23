package servbox

import (
	"bmautil/conndialpool"
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"fmt"
	"logger"
	"time"
)

func (this *Service) checkNode(node *nodeInfo) (*nodeInfo, error) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	old := this.nodes[node.name]
	for _, sn := range node.services {
		sinfo := this.servs[sn]
		if old == nil {
			if sinfo != nil && sinfo.node != nil {
				return nil, fmt.Errorf("service[%s] exists(%s)", sn, sinfo.node.name)
			}
		} else {
			if sinfo != nil && sinfo.node != old {
				return nil, fmt.Errorf("service[%s] engaged(%s)", sn, sinfo.node.name)
			}
		}
	}
	return old, nil
}

func (this *Service) doJoin(sock espsocket.Socket, msg *esnp.Message) error {
	var q objJoinQ
	err0 := q.Decode(msg)
	if err0 != nil {
		return err0
	}
	err0 = q.Valid()
	if err0 != nil {
		return err0
	}
	logger.Debug(tag, "join request(%s)", q)
	node := new(nodeInfo)
	node.name = q.NodeName
	node.net = q.Net
	node.address = q.Address
	node.services = q.SerivceNames
	// check services
	_, err0 = this.checkNode(node)
	if err0 != nil {
		return err0
	}
	if true {
		rmsg := msg.ReplyMessage()
		err0 = sock.WriteMessage(rmsg)
		if err0 != nil {
			return err0
		}
	}

	msg1, err1 := sock.ReadMessage()
	if err1 != nil {
		return err1
	}
	op, err2 := esnp.MessageLineCoders.Address.Get(msg1, esnp.ADDRESS_OP)
	if err2 != nil {
		return err2
	}
	if op != op_Active {
		return fmt.Errorf("invalid op(%s) after op(join)", op)
	}
	return this.doActive(sock, node, msg1)
}

func (this *Service) doReplace(node *nodeInfo, old *nodeInfo) (bool, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	old2 := this.nodes[node.name]
	if old2 != old {
		return false, nil
	}
	this.nodes[node.name] = node
	if old != nil {
		for _, sn := range old.services {
			delete(this.servs, sn)
		}
	}
	for _, sn := range node.services {
		si := new(servItem)
		si.node = node
		this.servs[sn] = si
	}
	return true, nil
}

func (this *Service) doActive(sock espsocket.Socket, node *nodeInfo, msg *esnp.Message) error {
	cfg := new(conndialpool.DialPoolConfig)
	cfg.Net = node.net
	cfg.Address = node.address
	cfg.InitSize = 1
	cfg.MaxSize = this.config.MaxConnSize
	cfg.TimeoutMS = this.config.TimeoutMS
	cfg.Valid()
	node.pool = conndialpool.NewDialPool(fmt.Sprintf("node_%s", node.name), cfg)
	if !node.pool.StartAndRun() {
		return fmt.Errorf("start dialpool fail")
	}
	conn, errC := node.pool.GetConn(time.Duration(this.config.TimeoutMS)*time.Second, true)
	if errC != nil {
		node.pool.Close()
		return errC
	}
	node.pool.ReturnConn(conn)

	for {
		old, err0 := this.checkNode(node)
		if err0 != nil {
			return err0
		}
		ok, err1 := this.doReplace(node, old)
		if err1 != nil {
			return err1
		}
		if ok {
			break
		}
	}
	if true {
		rmsg := msg.ReplyMessage()
		errR := sock.WriteMessage(rmsg)
		if errR != nil {
			return errR
		}
	}
	logger.Info(tag, "active(%s, %s, %s, %v)", node.name, node.net, node.address, node.services)
	// keep it
	for {
		if sock.IsBreak() {
			break
		}
		time.Sleep(1 * time.Second)
	}
	return nil
}
