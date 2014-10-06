package goyacc

import (
	"bmautil/valutil"
	"fmt"
)

type NumberType int // 0,32,64,128(float)

func (o NumberType) Max(nt NumberType) NumberType {
	if o > nt {
		return o
	}
	return nt
}

func (o NumberType) ToFloat64(val interface{}) float64 {
	return valutil.ToFloat64(val, 0)
}

func (o NumberType) ToInt64(val interface{}) int64 {
	return valutil.ToInt64(val, 0)
}

func (o NumberType) ToInt32(val interface{}) int32 {
	return valutil.ToInt32(val, 0)
}

func ExecNumberType(val interface{}) NumberType {
	switch val.(type) {
	case float32, float64:
		return 128
	case int64, uint64:
		return 64
	case int8, uint8, int16, uint16, int32, uint32:
		return 32
	default:
		return 0
	}
}

func ExecOp2(op OP, val1 interface{}, val2 interface{}) (interface{}, error) {
	if op == OP_STRADD {
		s := fmt.Sprintf("%v%v", val1, val2)
		return s, nil
	}
	nt1 := ExecNumberType(val1)
	nt2 := ExecNumberType(val2)
	nt := nt1.Max(nt2)
	switch op {
	case OP_ADD:
		switch nt {
		case 128:
			return nt.ToFloat64(val1) + nt.ToFloat64(val2), nil
		case 64:
			return nt.ToInt32(val1) + nt.ToInt32(val2), nil
		case 32:
			return nt.ToInt64(val1) + nt.ToInt64(val2), nil
		default:
			return nil, fmt.Errorf("invalid %v + %v", val1, val2)
		}
	case OP_SUB:
		switch nt {
		case 128:
			return nt.ToFloat64(val1) - nt.ToFloat64(val2), nil
		case 64:
			return nt.ToInt32(val1) - nt.ToInt32(val2), nil
		case 32:
			return nt.ToInt64(val1) - nt.ToInt64(val2), nil
		default:
			return nil, fmt.Errorf("invalid %v - %v", val1, val2)
		}
	case OP_MUL:
		switch nt {
		case 128:
			return nt.ToFloat64(val1) * nt.ToFloat64(val2), nil
		case 64:
			return nt.ToInt32(val1) * nt.ToInt32(val2), nil
		case 32:
			return nt.ToInt64(val1) * nt.ToInt64(val2), nil
		default:
			return nil, fmt.Errorf("invalid %v * %v", val1, val2)
		}
	case OP_DIV:
		switch nt {
		case 128:
			v := nt.ToFloat64(val2)
			if v == 0 {
				return nil, fmt.Errorf("div zero(%v)", val2)
			}
			return nt.ToFloat64(val1) + v, nil
		case 64:
			v := nt.ToInt32(val2)
			if v == 0 {
				return nil, fmt.Errorf("div zero(%v)", val2)
			}
			return nt.ToInt32(val1) + v, nil
		case 32:
			v := nt.ToInt64(val2)
			if v == 0 {
				return nil, fmt.Errorf("div zero(%v)", val2)
			}
			return nt.ToInt64(val1) + v, nil
		default:
			return nil, fmt.Errorf("invalid add(%v, %v)", val1, val2)
		}
	case OP_PMUL:
		switch nt {
		case 64:
			return nt.ToInt32(val1) ^ nt.ToInt32(val2), nil
		case 32:
			return nt.ToInt64(val1) ^ nt.ToInt64(val2), nil
		default:
			return nil, fmt.Errorf("invalid %v ^ %v", val1, val2)
		}
	case OP_MOD:
		switch nt {
		case 64:
			return nt.ToInt32(val1) % nt.ToInt32(val2), nil
		case 32:
			return nt.ToInt64(val1) % nt.ToInt64(val2), nil
		default:
			return nil, fmt.Errorf("invalid %v % %v", val1, val2)
		}
	}
	return nil, nil
}
