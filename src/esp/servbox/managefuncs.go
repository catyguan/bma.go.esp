package servbox

import (
	"bmautil/conndialpool"
	"esp/espnet/esnp"
	"esp/espnet/espsocket"
	"esp/services/servboot"
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
	logger.Debug(tag, "join request(%v)", q)
	node := new(nodeInfo)
	node.name = q.NodeName
	node.net = q.Net
	node.address = q.Address
	node.services = q.SerivceNames
	node.info = q.Info
	node.skipKill = q.SkipKill
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

	msg1, err1 := sock.ReadMessage(true)
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

func (this *Service) doRemove(node *nodeInfo) bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	old := this.nodes[node.name]
	if old != node {
		return false
	}
	delete(this.nodes, node.name)
	for _, sn := range node.services {
		delete(this.servs, sn)
	}
	node.pool.Close()
	return true
}

func (this *Service) doReplace(node *nodeInfo, old *nodeInfo) (bool, bool, error) {
	this.lock.Lock()
	defer this.lock.Unlock()
	old2 := this.nodes[node.name]
	if old2 != old {
		return false, false, nil
	}
	if node.isSame(old) {
		old.ReplaceProp(node)
		if node.pool != nil {
			node.pool.Close()
		}
		logger.Debug(tag, "same nodeInfo(%s, %s), skip", node.name, node.info)
		return true, false, nil
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
	return true, true, nil
}

func (this *Service) doKill(old *nodeInfo) error {
	pool := old.pool
	if pool == nil {
		return nil
	}
	if pool.ActiveConn() == 0 {
		return nil
	}
	conn, err := pool.GetConn(1*time.Second, true)
	if err != nil {
		return err
	}
	defer pool.CloseConn(conn)
	sock := espsocket.NewConnSocket(conn, 0)
	msg := esnp.NewRequestMessage()
	msg.GetAddress().SetCall(servboot.NAME_SERVICE, servboot.NAME_OP_SHUTDOWN)
	err1 := sock.WriteMessage(msg)
	if err1 != nil {
		return err1
	}
	_, err2 := sock.ReadMessage(true)
	return err2
}

func (this *Service) onReplace(old *nodeInfo) {
	if old == nil {
		return
	}
	if old.skipKill {
		return
	}
	if old.pool == nil {
		return
	}
	if old.pool.ActiveConn() == 0 {
		logger.Info(tag, "kill replaced(%s, %s) skip", old, old.pool)
		old.pool.Close()
		return
	}
	logger.Debug(tag, "kill replaced(%s) after %d sec", old.name, this.config.KillDelaySec)
	time.AfterFunc(time.Duration(this.config.KillDelaySec)*time.Second, func() {
		defer old.pool.Close()
		err := this.doKill(old)
		if err != nil {
			logger.Info(tag, "kill replaced(%s) fail - %s", old, err)
			return
		}
		logger.Info(tag, "kill replaced(%s) done", old)
	})
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

	replaced := false
	for {
		old, err0 := this.checkNode(node)
		if err0 != nil {
			return err0
		}
		// if old != nil {
		// 	logger.Info(tag, "%v -- %v", node.pool, old.pool)
		// }
		ok, done, err1 := this.doReplace(node, old)
		if err1 != nil {
			return err1
		}
		if ok {
			if done {
				replaced = true
				this.onReplace(old)
			}
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
	if !replaced {
		logger.Info(tag, "same node(%s,%s) active", node.name, node.info)
		return nil
	}
	logger.Info(tag, "active(%s,%s)(%s, %s) - %v", node.name, node.info, node.net, node.address, node.services)
	for {
		_, err := sock.ReadMessage(false)
		if err != nil {
			if this.doRemove(node) {
				logger.Info(tag, "node(%s,%s) remove on connection break", node.name, node.info)
				break
			}
		}
	}
	return nil
}
