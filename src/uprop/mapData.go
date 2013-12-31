package uprop

type MapData interface {
	ToMap() map[string]interface{}

	FromMap(data map[string]interface{}) error
}

func Copy(des MapData, src MapData) error {
	m := src.ToMap()
	return des.FromMap(m)
}
