package knife

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type PartInfo struct {
	Filename string
	Parts    int64
}

func LoadPartInfo(filename string) (*PartInfo, error) {
	info := &PartInfo{}
	data, e := ioutil.ReadFile(filename)
	if e != nil {
		return nil, e
	}
	if e = json.Unmarshal(data, info); e != nil {
		return nil, e
	}
	return info, nil
}

func SavePartInfo(filename string, info PartInfo) error {
	data, e := json.Marshal(info)
	if e != nil {
		return e
	}
	if e = ioutil.WriteFile(filename, data, os.ModePerm); e != nil {
		return e
	}
	return nil
}
