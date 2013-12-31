package valutil

import (
	"errors"
	"fmt"
	"math"
)

type SizeUnit int8

const (
	SizeB = 0
	SizeK = 1
	SizeM = 2
	SizeG = 3
)

func (level SizeUnit) CompareLevelTo(other SizeUnit) int {
	return int(level) - int(other)
}

func (level SizeUnit) String() string {
	switch level {
	case SizeB:
		return ""
	case SizeK:
		return "K"
	case SizeM:
		return "M"
	case SizeG:
		return "G"
	}
	return "UNKNOW"
}

func NewSizeUnit(s uint8) (SizeUnit, bool) {
	switch s {
	case 'B', 'b':
		return SizeB, true
	case 'K', 'k':
		return SizeK, true
	case 'M', 'm':
		return SizeM, true
	case 'G', 'g':
		return SizeG, true
	}
	return SizeB, false
}

func SizeString(size uint64, base int, sizeUnit SizeUnit) string {
	lv := math.Pow(float64(base), float64(sizeUnit))
	v := float64(size) / lv
	if math.Floor(v) == v {
		return fmt.Sprintf("%d%s", uint64(v), sizeUnit.String())
	} else {
		return fmt.Sprintf("%.3f%s", v, sizeUnit.String())
	}
}

func ToSize(origin string, base int, sizeUnit SizeUnit) (r uint64, err error) {
	if base < 1 {
		return 0, errors.New("Error objectUnit values")
	}

	length := len(origin)
	originUnit, ok := NewSizeUnit(origin[length-1])

	gap := originUnit.CompareLevelTo(sizeUnit)

	vstr := origin
	if ok {
		vstr = origin[0 : length-1]
	}
	originData := ToFloat64(vstr, -1)
	if originData < 0 {
		return 0, errors.New("origin data can't lower than 0")
	}
	return uint64(originData * math.Pow(float64(base), float64(gap))), nil
}
