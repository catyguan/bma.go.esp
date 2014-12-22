package securep

import (
	"bmautil/valutil"
	"crypto/md5"
	"encoding/hex"
	"esp/espnet/esnp"
)

type BaseAuthRequest struct {
	Token string
}

func (this *BaseAuthRequest) Encode(msg *esnp.Message) error {
	if this.Token != "" {
		esnp.MessageLineCoders.XData.Add(msg, 1, this.Token, nil)
	}
	return nil
}

func (this *BaseAuthRequest) Decode(msg *esnp.Message) error {
	if true {
		v, err := esnp.MessageLineCoders.XData.Get(msg, 1, nil)
		if err != nil {
			return err
		}
		this.Token = valutil.ToString(v, "")
	}
	return nil
}

func (this *BaseAuthRequest) Valid() error {
	return nil
}

func (this *BaseAuthRequest) Reset() {
	this.Token = ""
}

type BaseAuthResponse struct {
	Token string
}

func (this *BaseAuthResponse) Encode(msg *esnp.Message) error {
	if this.Token != "" {
		esnp.MessageLineCoders.XData.Add(msg, 1, this.Token, nil)
	}
	return nil
}

func (this *BaseAuthResponse) Decode(msg *esnp.Message) error {
	if true {
		v, err := esnp.MessageLineCoders.XData.Get(msg, 1, nil)
		if err != nil {
			return err
		}
		this.Token = valutil.ToString(v, "")
	}
	return nil
}

func (this *BaseAuthResponse) Valid() error {
	return nil
}

func (this *BaseAuthResponse) Reset() {
	this.Token = ""
}

func CreateAuthToken(tk string, k string) string {
	h := md5.New()
	h.Write([]byte(tk))
	h.Write([]byte(k))
	return hex.EncodeToString(h.Sum(nil))
}
