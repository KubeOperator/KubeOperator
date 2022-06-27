package migrations

import (
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type SourceDriver struct {
	Dir string
}

func NewSourceDriver(dir string) (*SourceDriver, error) {
	_, err := os.Stat(dir)
	if err != nil {
		return nil, err
	}
	sd := &SourceDriver{
		Dir: dir,
	}
	return sd, nil
}

func (s SourceDriver) ReadUp(version int) (*Migrations, error) {
	files, err := ioutil.ReadDir(s.Dir)
	if err != nil {
		return nil, err
	}
	ms := &Migrations{
		Index:      uintSlice{},
		Migrations: make(map[int]Migration, len(files)),
	}
	for _, f := range files {
		vstr := f.Name()[:strings.Index(f.Name(), "_")]
		v, err := strconv.Atoi(vstr)
		if err != nil {
			return nil, errors.New("file is not valid")
		}
		if v > version {
			ms.Migrations[v] = Migration{
				Version:  v,
				FileName: f.Name(),
			}
			ms.Index = append(ms.Index, v)
		}
	}
	ms.buildIndex()
	return ms, nil
}
