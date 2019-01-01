package token

import "io/ioutil"

type Source interface {
	Name() string
	Code() []byte
}

type file struct {
	Path string
	code []byte
}

func ReadFile(path string) (Source, error) {
	code, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return &file{path, code}, nil
}

func (f *file) Name() string {
	return f.Path
}

func (f *file) Code() []byte {
	if f.code == nil {
	}
	return f.code
}
