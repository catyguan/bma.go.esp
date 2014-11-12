package memblock

type ListOfRemoveListener struct {
	Collections []RemoveListener
}

func (this *ListOfRemoveListener) ListenerFunc(key string, item *MapItem, rt REMOVE_TYPE) {
	for _, lis := range this.Collections {
		lis(key, item, rt)
	}
}
