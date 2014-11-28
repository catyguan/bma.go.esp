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
	case int8, uint8, int16, uint16, int32, uint32, int, uint:
		return 32
	case string:
		return 1
	default:
		return 0
	}
}

func ExecOp2(op OP, val1 interface{}, val2 interface{}) (bool, interface{}, error) {
	if op == OP_STRADD {
		s := fmt.Sprintf("%v%v", val1, val2)
		return true, s, nil
	}
	nt1 := ExecNumberType(val1)
	nt2 := ExecNumberType(val2)
	nt := nt1.Max(nt2)
	switch op {
	case OP_ADD:
		switch nt {
		case 128:
			return true, nt.ToFloat64(val1) + nt.ToFloat64(val2), nil
		case 32:
			return true, nt.ToInt32(val1) + nt.ToInt32(val2), nil
		case 64:
			return true, nt.ToInt64(val1) + nt.ToInt64(val2), nil
		default:
			return false, nil, fmt.Errorf("invalid %v + %v", val1, val2)
		}
	case OP_SUB:
		switch nt {
		case 128:
			return true, nt.ToFloat64(val1) - nt.ToFloat64(val2), nil
		case 32:
			return true, nt.ToInt32(val1) - nt.ToInt32(val2), nil
		case 64:
			return true, nt.ToInt64(val1) - nt.ToInt64(val2), nil
		default:
			return false, nil, fmt.Errorf("invalid %v - %v", val1, val2)
		}
	case OP_MUL:
		switch nt {
		case 128:
			return true, nt.ToFloat64(val1) * nt.ToFloat64(val2), nil
		case 32:
			return true, nt.ToInt32(val1) * nt.ToInt32(val2), nil
		case 64:
			return true, nt.ToInt64(val1) * nt.ToInt64(val2), nil
		default:
			return false, nil, fmt.Errorf("invalid %v * %v", val1, val2)
		}
	case OP_DIV:
		switch nt {
		case 128:
			v := nt.ToFloat64(val2)
			if v == 0 {
				return false, nil, fmt.Errorf("div zero(%v)", val2)
			}
			return true, nt.ToFloat64(val1) + v, nil
		case 32:
			v := nt.ToInt32(val2)
			if v == 0 {
				return false, nil, fmt.Errorf("div zero(%v)", val2)
			}
			return true, nt.ToInt32(val1) + v, nil
		case 64:
			v := nt.ToInt64(val2)
			if v == 0 {
				return true, nil, fmt.Errorf("div zero(%v)", val2)
			}
			return true, nt.ToInt64(val1) + v, nil
		default:
			return false, nil, fmt.Errorf("invalid add(%v, %v)", val1, val2)
		}
	case OP_PMUL:
		switch nt {
		case 32:
			return true, nt.ToInt32(val1) ^ nt.ToInt32(val2), nil
		case 64:
			return true, nt.ToInt64(val1) ^ nt.ToInt64(val2), nil
		default:
			return false, nil, fmt.Errorf("invalid %v ^ %v", val1, val2)
		}
	case OP_MOD:
		switch nt {
		case 32:
			return true, nt.ToInt32(val1) % nt.ToInt32(val2), nil
		case 64:
			return true, nt.ToInt64(val1) % nt.ToInt64(val2), nil
		default:
			return false, nil, fmt.Errorf("invalid %v % %v", val1, val2)
		}
	case OP_AND:
		if val1 == nil || val2 == nil {
			return true, false, nil
		}
		v1, ok1 := val1.(bool)
		if !ok1 {
			v1 = valutil.ToBool(val1, false)
		}
		v2, ok2 := val2.(bool)
		if !ok2 {
			v2 = valutil.ToBool(val2, false)
		}
		return true, v1 && v2, nil
	case OP_OR:
		if val1 == nil || val2 == nil {
			return true, false, nil
		}
		v1, ok1 := val1.(bool)
		if !ok1 {
			v1 = valutil.ToBool(val1, false)
		}
		v2, ok2 := val2.(bool)
		if !ok2 {
			v2 = valutil.ToBool(val2, false)
		}
		return true, v1 || v2, nil
	case OP_LT:
		switch nt {
		case 128:
			return true, nt.ToFloat64(val1) < nt.ToFloat64(val2), nil
		case 32:
			return true, nt.ToInt32(val1) < nt.ToInt32(val2), nil
		case 64:
			return true, nt.ToInt64(val1) < nt.ToInt64(val2), nil
		default:
			if nt1 == 1 && nt2 == 1 {
				s1 := val1.(string)
				s2 := val2.(string)
				return true, s1 < s2, nil
			}
			return false, nil, fmt.Errorf("invalid %v < %v", val1, val2)
		}
	case OP_LTEQ:
		switch nt {
		case 128:
			return true, nt.ToFloat64(val1) <= nt.ToFloat64(val2), nil
		case 32:
			return true, nt.ToInt32(val1) <= nt.ToInt32(val2), nil
		case 64:
			return true, nt.ToInt64(val1) <= nt.ToInt64(val2), nil
		default:
			if nt1 == 1 && nt2 == 1 {
				s1 := val1.(string)
				s2 := val2.(string)
				return true, s1 <= s2, nil
			}
			return false, nil, fmt.Errorf("invalid %v <= %v", val1, val2)
		}
	case OP_GT:
		switch nt {
		case 128:
			return true, nt.ToFloat64(val1) > nt.ToFloat64(val2), nil
		case 32:
			return true, nt.ToInt32(val1) > nt.ToInt32(val2), nil
		case 64:
			return true, nt.ToInt64(val1) > nt.ToInt64(val2), nil
		default:
			if nt1 == 1 && nt2 == 1 {
				s1 := val1.(string)
				s2 := val2.(string)
				return true, s1 > s2, nil
			}
			return false, nil, fmt.Errorf("invalid %v > %v", val1, val2)
		}
	case OP_GTEQ:
		switch nt {
		case 128:
			return true, nt.ToFloat64(val1) >= nt.ToFloat64(val2), nil
		case 32:
			return true, nt.ToInt32(val1) >= nt.ToInt32(val2), nil
		case 64:
			return true, nt.ToInt64(val1) >= nt.ToInt64(val2), nil
		default:
			if nt1 == 1 && nt2 == 1 {
				s1 := val1.(string)
				s2 := val2.(string)
				return true, s1 >= s2, nil
			}
			return false, nil, fmt.Errorf("invalid %v >= %v", val1, val2)
		}
	case OP_EQ, OP_NOTEQ:
		r := false
		if s1, ok1 := val1.(string); ok1 {
			if s2, ok2 := val2.(string); ok2 {
				r = s1 == s2
				if op == OP_NOTEQ {
					r = !r
				}
				return true, r, nil
			}
		}
		if b1, ok1 := val1.(bool); ok1 {
			if b2, ok2 := val2.(bool); ok2 {
				r = b1 == b2
				if op == OP_NOTEQ {
					r = !r
				}
				return true, r, nil
			}
		}
		if val1 == nil {
			r = val2 == nil
			if op == OP_NOTEQ {
				r = !r
			}
			return true, r, nil
		}
		if val2 == nil {
			r = val1 == nil
			if op == OP_NOTEQ {
				r = !r
			}
			return true, r, nil
		}
		if nt1 == 0 || nt2 == 0 {
			return false, nil, fmt.Errorf("invalid %v,%T == %v,%T", val1, val1, val2, val2)
		}
		switch nt {
		case 128:
			r = nt.ToFloat64(val1) == nt.ToFloat64(val2)
		case 32:
			r = nt.ToInt32(val1) == nt.ToInt32(val2)
		case 64:
			r = nt.ToInt64(val1) == nt.ToInt64(val2)
		default:
			return false, nil, fmt.Errorf("invalid %v == %v", val1, val2)
		}
		if op == OP_NOTEQ {
			r = !r
		}
		return true, r, nil
	}
	return false, nil, nil
}
