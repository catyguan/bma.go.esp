package election

import (
	"bmautil/byteutil"
	"bmautil/coder"
	"esp/cluster/nodeid"
	"esp/espnet"
)

// // string
// type stringCoder int

// func (this stringCoder) DoEncode(w *byteutil.BytesBufferWriter, v string) {
// 	w.Write([]byte(v))
// }

// func (this stringCoder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
// 	this.DoEncode(w, v.(string))
// 	return nil
// }

// func (this stringCoder) DoDecode(r *byteutil.BytesBufferReader) string {
// 	return string(r.ReadAll())
// }

// func (this stringCoder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
// 	s := this.DoDecode(r)
// 	return s, nil
// }

var (
	EpochIdCoder = epochIdCoder(0)
	StatusCoder  = statusCoder(0)
)

type epochIdCoder int

func (O epochIdCoder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	coder.Uint64.DoEncode(w, uint64(v.(EpochId)))
	return nil
}

func (O epochIdCoder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	v, err := coder.Uint64.DoDecode(r)
	if err != nil {
		return nil, err
	}
	return EpochId(v), nil
}

type statusCoder int

func (O statusCoder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	coder.Uint8.DoEncode(w, uint8(v.(Status)))
	return nil
}

func (O statusCoder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	v, err := coder.Uint8.DoDecode(r)
	if err != nil {
		return nil, err
	}
	return Status(v), nil
}

const (
	OP_VOTE_REQ      = "vr"
	OP_VOTE_RESP     = "vp"
	OP_ANNOUNCE_REQ  = "ar"
	OP_ANNOUNCE_RESP = "ap"
)

// CandidateState
func (this *CandidateState) WriteState(xd *espnet.MessageXData) error {
	xd.Add(1, this.Id, nodeid.Coder)
	xd.Add(2, this.Epoch, EpochIdCoder)
	xd.Add(3, this.Status, StatusCoder)
	xd.Add(4, this.Leader, nodeid.Coder)
	return nil
}

func (this *CandidateState) ReadState(it *espnet.XDataIterator) (bool, error) {
	switch it.Xid() {
	case 1:
		v, err := it.Value(nodeid.Coder)
		if err != nil {
			return true, err
		}
		this.Id = v.(nodeid.NodeId)
	case 2:
		v, err := it.Value(EpochIdCoder)
		if err != nil {
			return true, err
		}
		this.Epoch = v.(EpochId)
	case 3:
		v, err := it.Value(StatusCoder)
		if err != nil {
			return true, err
		}
		this.Status = v.(Status)
	case 4:
		v, err := it.Value(nodeid.Coder)
		if err != nil {
			return true, err
		}
		this.Leader = v.(nodeid.NodeId)
	default:
		return false, nil
	}
	return true, nil
}

// VoteReq
func (this *VoteReq) Write(msg *espnet.Message) error {
	msg.GetAddress().Set(espnet.ADDRESS_OP, OP_VOTE_REQ)
	xd := msg.XDatas()
	this.WriteState(xd)
	xd.Add(5, this.Proposal, nodeid.Coder)
	xd.Add(6, this.Renew, coder.Bool)
	return nil
}

func (this *VoteReq) Read(msg *espnet.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		ok, err0 := this.ReadState(it)
		if err0 != nil {
			return err0
		}
		if ok {
			continue
		}
		switch it.Xid() {
		case 5:
			v, err := it.Value(nodeid.Coder)
			if err != nil {
				return err
			}
			this.Leader = v.(nodeid.NodeId)
		case 6:
			v, err := it.Value(coder.Bool)
			if err != nil {
				return err
			}
			this.Renew = v.(bool)
		}
	}
	return nil
}

// VoteResp
func (this *VoteResp) Write(msg *espnet.Message) error {
	msg.GetAddress().Set(espnet.ADDRESS_OP, OP_VOTE_RESP)
	xd := msg.XDatas()
	if this.State != nil {
		this.State.WriteState(xd)
	}
	xd.Add(5, this.Id, nodeid.Coder)
	xd.Add(6, this.Epoch, EpochIdCoder)
	xd.Add(7, this.Accept, coder.Bool)
	return nil
}

func (this *VoteResp) Read(msg *espnet.Message) error {
	var st CandidateState
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		ok, err0 := st.ReadState(it)
		if err0 != nil {
			return err0
		}
		if ok {
			continue
		}
		switch it.Xid() {
		case 5:
			v, err := it.Value(nodeid.Coder)
			if err != nil {
				return err
			}
			this.Id = v.(nodeid.NodeId)
		case 6:
			v, err := it.Value(EpochIdCoder)
			if err != nil {
				return err
			}
			this.Epoch = v.(EpochId)
		case 7:
			v, err := it.Value(coder.Bool)
			if err != nil {
				return err
			}
			this.Accept = v.(bool)
		}
	}
	if st.Id != nodeid.INVALID {
		this.State = &st
	}
	return nil
}

// AnnounceReq
func (this *AnnounceReq) Write(msg *espnet.Message) error {
	msg.GetAddress().Set(espnet.ADDRESS_OP, OP_ANNOUNCE_REQ)
	xd := msg.XDatas()
	this.WriteState(xd)
	return nil
}

func (this *AnnounceReq) Read(msg *espnet.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		ok, err0 := this.ReadState(it)
		if err0 != nil {
			return err0
		}
		if ok {
			continue
		}
		switch it.Xid() {
		}
	}
	return nil
}

// AnnounceResp
func (this *AnnounceResp) Write(msg *espnet.Message) error {
	msg.GetAddress().Set(espnet.ADDRESS_OP, OP_ANNOUNCE_RESP)
	xd := msg.XDatas()
	if this.State != nil {
		this.State.WriteState(xd)
	}
	xd.Add(5, this.Id, nodeid.Coder)
	xd.Add(6, this.Epoch, EpochIdCoder)
	xd.Add(7, this.Accept, coder.Bool)
	return nil
}

func (this *AnnounceResp) Read(msg *espnet.Message) error {
	var st CandidateState
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		ok, err0 := st.ReadState(it)
		if err0 != nil {
			return err0
		}
		if ok {
			continue
		}
		switch it.Xid() {
		case 5:
			v, err := it.Value(nodeid.Coder)
			if err != nil {
				return err
			}
			this.Id = v.(nodeid.NodeId)
		case 6:
			v, err := it.Value(EpochIdCoder)
			if err != nil {
				return err
			}
			this.Epoch = v.(EpochId)
		case 7:
			v, err := it.Value(coder.Bool)
			if err != nil {
				return err
			}
			this.Accept = v.(bool)
		}
	}
	if st.Id != nodeid.INVALID {
		this.State = &st
	}
	return nil
}
