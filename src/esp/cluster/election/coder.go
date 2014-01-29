package election

import (
	"bmautil/byteutil"
	"bmautil/coder"
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
