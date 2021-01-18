package cache

import (
	"db/Util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)
var f interface {
}

var cacheMap = map[string]interface{}{}



func ReadData() {

	err := filepath.Walk("data",
		func(path string, info os.FileInfo, err error) error {

			Util.Check(err)

			fmt.Println(path, info.Size())
			if strings.HasSuffix(path, ".json") {

				file, err := ioutil.ReadFile(path)
				err = json.Unmarshal(file, &f)
				m := f.(map[string]interface{})

				i := m["_id"]
				j := i.(string)

				Util.Check(err)
				cacheMap[j] = m


			}
			return nil
		})

	Util.Check(err)
}

func IsInCache(id string) interface{} {

	s := cacheMap[id]
	if cacheMap[id] != "" {
		return s
	} else {
		return ""
	}

}

func GetCache() *map[string]interface{} {
	return &cacheMap
}
