package election

import (
	"bmautil/byteutil"
	"bmautil/coder"
	"esp/cluster/nodeinfo"
	"esp/espnet/esnp"
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
	EpochIdCoder        = epochIdCoder(0)
	StatusCoder         = statusCoder(0)
	CandidateStateCoder = candidateStateCoder(0)
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

type candidateStateCoder int

func (O candidateStateCoder) Encode(w *byteutil.BytesBufferWriter, v interface{}) error {
	if v == nil {
		coder.Int32.DoEncode(w, 0)
		return nil
	}
	cs := v.(*CandidateState)
	coder.Int32.DoEncode(w, 1)
	if err := nodeinfo.NodeIdCoder.Encode(w, cs.Id); err != nil {
		return err
	}
	if err := EpochIdCoder.Encode(w, cs.Epoch); err != nil {
		return err
	}
	if err := StatusCoder.Encode(w, cs.Status); err != nil {
		return err
	}
	if err := nodeinfo.NodeIdCoder.Encode(w, cs.Leader); err != nil {
		return err
	}
	return nil
}

func (O candidateStateCoder) Decode(r *byteutil.BytesBufferReader) (interface{}, error) {
	v0, err0 := coder.Int32.DoDecode(r)
	if err0 != nil {
		return nil, err0
	}
	if v0 == 0 {
		return nil, nil
	}
	cs := new(CandidateState)
	v1, err1 := nodeinfo.NodeIdCoder.Decode(r)
	if err1 != nil {
		return nil, err1
	}
	cs.Id = v1.(nodeinfo.NodeId)

	v2, err2 := EpochIdCoder.Decode(r)
	if err2 != nil {
		return nil, err2
	}
	cs.Epoch = v2.(EpochId)

	v3, err3 := StatusCoder.Decode(r)
	if err3 != nil {
		return nil, err3
	}
	cs.Status = v3.(Status)

	v4, err4 := nodeinfo.NodeIdCoder.Decode(r)
	if err4 != nil {
		return nil, err4
	}
	cs.Leader = v4.(nodeinfo.NodeId)

	return cs, nil
}

const (
	OP_VOTE_REQ      = "vr"
	OP_VOTE_RESP     = "vp"
	OP_ANNOUNCE_REQ  = "ar"
	OP_ANNOUNCE_RESP = "ap"
)

// VoteReq
func (this *VoteReq) Write(msg *esnp.Message) error {
	msg.GetAddress().SetOp(OP_VOTE_REQ)
	xd := msg.XDatas()
	xd.Add(1, this.State, CandidateStateCoder)
	xd.Add(2, this.Proposal, nodeinfo.NodeIdCoder)
	xd.Add(3, this.Renew, coder.Bool)
	return nil
}

func (this *VoteReq) Read(msg *esnp.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 1:
			v, err := it.Value(CandidateStateCoder)
			if err != nil {
				return err
			}
			if v != nil {
				this.State = v.(*CandidateState)
			}
		case 2:
			v, err := it.Value(nodeinfo.NodeIdCoder)
			if err != nil {
				return err
			}
			this.Proposal = v.(nodeinfo.NodeId)
		case 3:
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
func (this *VoteResp) Write(msg *esnp.Message) error {
	msg.GetAddress().SetOp(OP_VOTE_RESP)
	xd := msg.XDatas()
	xd.Add(1, this.Id, nodeinfo.NodeIdCoder)
	xd.Add(2, this.Epoch, EpochIdCoder)
	xd.Add(3, this.Accept, coder.Bool)
	xd.Add(4, this.State, CandidateStateCoder)
	return nil
}

func (this *VoteResp) Read(msg *esnp.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 1:
			v, err := it.Value(nodeinfo.NodeIdCoder)
			if err != nil {
				return err
			}
			this.Id = v.(nodeinfo.NodeId)
		case 2:
			v, err := it.Value(EpochIdCoder)
			if err != nil {
				return err
			}
			this.Epoch = v.(EpochId)
		case 3:
			v, err := it.Value(coder.Bool)
			if err != nil {
				return err
			}
			this.Accept = v.(bool)
		case 4:
			v, err := it.Value(CandidateStateCoder)
			if err != nil {
				return err
			}
			if v != nil {
				this.State = v.(*CandidateState)
			}
		}
	}
	return nil
}

// AnnounceReq
func (this *AnnounceReq) Write(msg *esnp.Message) error {
	msg.GetAddress().SetOp(OP_ANNOUNCE_REQ)
	xd := msg.XDatas()
	xd.Add(1, this.State, CandidateStateCoder)
	return nil
}

func (this *AnnounceReq) Read(msg *esnp.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 1:
			v, err := it.Value(CandidateStateCoder)
			if err != nil {
				return err
			}
			if v != nil {
				this.State = v.(*CandidateState)
			}
		}
	}
	return nil
}

// AnnounceResp
func (this *AnnounceResp) Write(msg *esnp.Message) error {
	msg.GetAddress().Set(esnp.ADDRESS_OP, OP_ANNOUNCE_RESP)
	xd := msg.XDatas()
	xd.Add(1, this.Id, nodeinfo.NodeIdCoder)
	xd.Add(2, this.Epoch, EpochIdCoder)
	xd.Add(3, this.Accept, coder.Bool)
	xd.Add(4, this.State, CandidateStateCoder)
	return nil
}

func (this *AnnounceResp) Read(msg *esnp.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 1:
			v, err := it.Value(nodeinfo.NodeIdCoder)
			if err != nil {
				return err
			}
			this.Id = v.(nodeinfo.NodeId)
		case 2:
			v, err := it.Value(EpochIdCoder)
			if err != nil {
				return err
			}
			this.Epoch = v.(EpochId)
		case 3:
			v, err := it.Value(coder.Bool)
			if err != nil {
				return err
			}
			this.Accept = v.(bool)
		case 4:
			v, err := it.Value(CandidateStateCoder)
			if err != nil {
				return err
			}
			if v != nil {
				this.State = v.(*CandidateState)
			}
		}
	}
	return nil
}
