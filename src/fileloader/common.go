package fileloader

type commonFileLoaderFactory int

const (
	CommonFileLoaderFactory = commonFileLoaderFactory(0)
)

func (this commonFileLoaderFactory) Valid(cfg map[string]interface{}) error {
	fac, _, err := GetFileLoaderFactoryByType(cfg)
	if err != nil {
		return err
	}
	return fac.Valid(cfg)
}

func (this commonFileLoaderFactory) Compare(cfg map[string]interface{}, old map[string]interface{}) bool {
	fac1, xt1, err1 := GetFileLoaderFactoryByType(cfg)
	if err1 != nil {
		return false
	}
	_, xt2, err2 := GetFileLoaderFactoryByType(old)
	if err2 != nil {
		return false
	}
	if xt1 != xt2 {
		return false
	}
	return fac1.Compare(cfg, old)
}

func (this commonFileLoaderFactory) Create(cfg map[string]interface{}) (FileLoader, error) {
	fac, _, err := GetFileLoaderFactoryByType(cfg)
	if err != nil {
		return nil, err
	}
	return fac.Create(cfg)
}
