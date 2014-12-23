package servbox

import (
	"bmautil/valutil"
	"esp/espnet/esnp"
	"fmt"
)

const (
	op_Join   = "join"
	op_Active = "active"
)

type objJoinQ struct {
	NodeName     string
	Net          string
	Address      string
	SerivceNames []string
}

func (this *objJoinQ) String() string {
	return fmt.Sprintf("%s,%s,%s,%v", this.NodeName, this.Net, this.Address, this.SerivceNames)
}

func (this *objJoinQ) Valid() error {
	if this.NodeName == "" {
		return fmt.Errorf("NodeName empty")
	}
	if this.Net == "" {
		return fmt.Errorf("Net empty")
	}
	if this.Address == "" {
		return fmt.Errorf("Address empty")
	}
	return nil
}

func (this *objJoinQ) Encode(msg *esnp.Message) error {
	esnp.MessageLineCoders.XData.Add(msg, 1, this.NodeName, esnp.Coders.String)
	esnp.MessageLineCoders.XData.Add(msg, 2, this.Net, esnp.Coders.String)
	esnp.MessageLineCoders.XData.Add(msg, 3, this.Address, esnp.Coders.String)
	if this.SerivceNames != nil {
		data := make([]interface{}, len(this.SerivceNames))
		for i, n := range this.SerivceNames {
			data[i] = n
		}
		esnp.MessageLineCoders.XData.Add(msg, 4, data, nil)
	}
	return nil
}

func (this *objJoinQ) Decode(msg *esnp.Message) error {
	it := msg.XDataIterator()
	for ; !it.IsEnd(); it.Next() {
		switch it.Xid() {
		case 1:
			v, err := it.Value(esnp.Coders.String)
			if err != nil {
				return err
			}
			this.NodeName = valutil.ToString(v, "")
		case 2:
			v, err := it.Value(esnp.Coders.String)
			if err != nil {
				return err
			}
			this.Net = valutil.ToString(v, "")
		case 3:
			v, err := it.Value(esnp.Coders.String)
			if err != nil {
				return err
			}
			this.Address = valutil.ToString(v, "")
		case 4:
			v, err := it.Value(nil)
			if err != nil {
				return err
			}
			if list, ok := v.([]interface{}); ok {
				snlist := make([]string, len(list))
				for i, sn := range list {
					snlist[i] = valutil.ToString(sn, "")
				}
				this.SerivceNames = snlist
			}
		}
	}
	return nil
}
