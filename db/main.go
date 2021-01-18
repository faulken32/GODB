package main

import (
	"db/Util"
	"db/cache"
	"db/find"
	"encoding/json"
	"github.com/bradhe/stopwatch"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	data = "data/"
)

func hello(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	marshal, _ := json.Marshal("GenericType")
	_, _ = w.Write(marshal)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func create(w http.ResponseWriter, r *http.Request) {

	log.Print("create")

	body := r.Body
	vars := mux.Vars(r)
	o := vars["object"]

	all, _ := ioutil.ReadAll(body)

	_ = json.Unmarshal(all, &Util.GenericType)
	u1 := uuid.Must(uuid.NewV4())
	id := u1.String()

	m := Util.GenericType.(map[string]interface{})
	m["_id"] = id

	marshal, _ := json.Marshal(m)
	log.Print(marshal)
	if _, err := os.Stat(data + o); os.IsNotExist(err) {
		err := os.Mkdir(data+o, 0755)
		check(err)
	}

	log.Print(data + o + "/" + id + ".json")
	err := ioutil.WriteFile(data+o+"/"+id+".json", marshal, 0644)
	check(err)
	body = nil

}

func getById(w http.ResponseWriter, r *http.Request) {

	log.Print("get")

	watch := stopwatch.Start()
	vars := mux.Vars(r)

	s := vars["id"]
	o := vars["object"]

	inCache := cache.IsInCache(s)

	var marshal []byte

	if inCache != "" {
		log.Println("find in cache", inCache)

		res, err := json.Marshal(inCache)
		Util.Check(err)
		marshal = res
		watch.Stop()
		log.Println(watch.Milliseconds())
	} else {
		file, err := ioutil.ReadFile(data + o + "/" + s + ".json")
		err = json.Unmarshal(file, &Util.GenericType)
		m := Util.GenericType.(map[string]interface{})
		check(err)
		res, _ := json.Marshal(m)
		marshal = res
		watch.Stop()
		log.Println(watch.Milliseconds())
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(marshal)
	s = ""
	vars = nil

}

func findBy(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	o := vars["object"]
	what := vars["what"]
	t := vars["term"]

	by := find.By(what, t, o)
	marshal, err := json.Marshal(by)
	Util.Check(err)

	w.Header().Set("Content-Type", "application/json")
	watch := stopwatch.Start()
	_, _ = w.Write(marshal)
	log.Println(by)
	watch.Stop()
	log.Println(watch.Milliseconds())
}

func main() {

	log.Print("stating server")
	go cache.ReadData()

	r := mux.NewRouter()
	r.HandleFunc("/", hello)
	r.HandleFunc("/create/{object}", create)
	r.HandleFunc("/get/{object}/{id}", getById)
	r.HandleFunc("/find/{object}/{what}/{term}", findBy)

	log.Fatal(http.ListenAndServe(":8080", r))

}
