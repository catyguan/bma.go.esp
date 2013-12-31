package valutil

// CoreValue
var (
	CoreValue primitiveValueFactory
)

type primitiveValueFactory int

func (this *primitiveValueFactory) Nil() PrimitiveValue {
	return PrimitiveValue{nil}
}
func (this *primitiveValueFactory) Int(v int) PrimitiveValue {
	return PrimitiveValue{v}
}
func (this *primitiveValueFactory) Int8(v int8) PrimitiveValue {
	return PrimitiveValue{v}
}
func (this *primitiveValueFactory) Int16(v int16) PrimitiveValue {
	return PrimitiveValue{v}
}
func (this *primitiveValueFactory) Int32(v int32) PrimitiveValue {
	return PrimitiveValue{v}
}
func (this *primitiveValueFactory) Int64(v int64) PrimitiveValue {
	return PrimitiveValue{v}
}
func (this *primitiveValueFactory) Uint8(v uint8) PrimitiveValue {
	return PrimitiveValue{v}
}
func (this *primitiveValueFactory) Uint16(v uint16) PrimitiveValue {
	return PrimitiveValue{v}
}
func (this *primitiveValueFactory) Uint32(v uint32) PrimitiveValue {
	return PrimitiveValue{v}
}
func (this *primitiveValueFactory) Uint64(v uint64) PrimitiveValue {
	return PrimitiveValue{v}
}
func (this *primitiveValueFactory) Float32(v float32) PrimitiveValue {
	return PrimitiveValue{v}
}
func (this *primitiveValueFactory) Float64(v float64) PrimitiveValue {
	return PrimitiveValue{v}
}
func (this *primitiveValueFactory) Bool(v bool) PrimitiveValue {
	return PrimitiveValue{v}
}
func (this *primitiveValueFactory) String(v string) PrimitiveValue {
	return PrimitiveValue{v}
}

type PrimitiveValue struct {
	value interface{}
}

func (this *PrimitiveValue) IsNil() bool {
	return this.value == nil
}
func (this *PrimitiveValue) ToInt() int {
	return ToInt(this.value, 0)
}
func (this *PrimitiveValue) ToInt8() int8 {
	return ToInt8(this.value, 0)
}
func (this *PrimitiveValue) ToInt16() int16 {
	return ToInt16(this.value, 0)
}
func (this *PrimitiveValue) ToInt32() int32 {
	return ToInt32(this.value, 0)
}
func (this *PrimitiveValue) ToInt64() int64 {
	return ToInt64(this.value, 0)
}
func (this *PrimitiveValue) ToUint8() uint8 {
	return ToUint8(this.value, 0)
}
func (this *PrimitiveValue) ToUint16() uint16 {
	return ToUint16(this.value, 0)
}
func (this *PrimitiveValue) ToUint32() uint32 {
	return ToUint32(this.value, 0)
}
func (this *PrimitiveValue) ToUint64() uint64 {
	return ToUint64(this.value, 0)
}
func (this *PrimitiveValue) ToFloat32() float32 {
	return ToFloat32(this.value, 0)
}
func (this *PrimitiveValue) ToFloat64() float64 {
	return ToFloat64(this.value, 0)
}
func (this *PrimitiveValue) ToBool() bool {
	return ToBool(this.value, false)
}
func (this *PrimitiveValue) ToString() string {
	return ToString(this.value, "")
}
func (this *PrimitiveValue) String() string {
	return this.ToString()
}
func (this *PrimitiveValue) Any() interface{} {
	return this.value
}
