package find

import (
	"db/Util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

/*

 */
func By(what string, terms string, collection string) map[string]interface{} {

	return readDataCollection(collection, terms, what)

}

func readDataCollection(collection string, terms string, what string) map[string]interface{} {

	dir, err := ioutil.ReadDir(Util.Data)
	Util.Check(err)
	data := map[string]interface{}{}
	for _, info := range dir {

		log.Println("info: ", info.Name())
		if info.Name() != collection {
			log.Fatal("collection does not exit")
		} else {
			err := filepath.Walk("data",
				func(path string, info os.FileInfo, err error) error {

					Util.Check(err)

					fmt.Println(path, info.Size())
					if strings.HasSuffix(path, ".json") {

						file, err := ioutil.ReadFile(path)
						err = json.Unmarshal(file, &Util.GenericType)
						m := Util.GenericType.(map[string]interface{})
						i := m["_id"]
						j := i.(string)
						if m[what] != nil && m[what] == terms {

							data[j] = m
						}

						Util.Check(err)

					}
					return nil
				})
			Util.Check(err)

			//ioutil.ReadAll(Util.GetDataDir() + "/" +collection)
		}
	}
	return data
}
