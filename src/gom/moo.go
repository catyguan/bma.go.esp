package gom

type moo struct {
}

func (this *moo) Rawget(key string) interface{} {
	return nil
}

func (this *moo) Rawset(key string, val interface{}) {

}

func (this *moo) Delete(key string) {

}

func (this *moo) Len() int {
	return 0
}

func (this *moo) ToMap() map[string]interface{} {
	return map[string]interface{}{}
}
