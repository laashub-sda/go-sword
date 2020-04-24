package core

import (
	"database/sql"
	"encoding/json"
	"go-sword/config"
	"go-sword/controller/render"
	"go-sword/response"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// Engine
type Engine struct {
	Config *config.Config
}

// Create Engine
func Default() *Engine {
	return &Engine{
		Config: &config.Config{},
	}
}

// Set Config like Database
func (e *Engine) SetConfig(cfg *config.Config) {
	e.Config = cfg
}

func (e *Engine) Run() {

	// Cache Panic
	defer func() {
		if err := recover(); err != nil {
			log.Printf("%v", err)
		}
	}()

	//http.HandleFunc("/sword/api/model/create", e.modelCreate)
	http.HandleFunc("/sword/api/crud/create", e.crudCreate)
	http.HandleFunc("/sword/api/model/table_list", e.tableList)
	http.HandleFunc("/sword/api/model/preview", e.modelPreview)

	// home page
	http.Handle("/sword/", http.StripPrefix("/sword/", http.FileServer(http.Dir("resource/web/base/dist"))))

	// render vue component
	http.HandleFunc("/sword/render", render.Render)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("服务启动失败 %v", err)
	}
}

func (e *Engine) tableList(w http.ResponseWriter, r *http.Request) {
	// user:password@(localhost)/dbname?charset=utf8&parseTime=True&loc=Local
	c := e.Config.Database
	db, err := sql.Open("mysql", c.User+":"+c.Password+"@tcp("+c.Host+":"+strconv.Itoa(c.Port)+")/"+c.Database+"?&parseTime=True")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	rows, err := db.Query("show tables")
	if err != nil {
		panic(err.Error())
	}

	tables := []string{}

	for rows.Next() {
		var tableName string
		rows.Scan(&tableName)
		tables = append(tables, tableName)
	}

	jsonData, err := json.Marshal(response.Ret{
		Code: http.StatusOK,
		Data: response.List{
			List: tables,
		},
	})

	_, err = w.Write(jsonData)
	if err != nil {
		panic(err.Error())
	}
}

func (e *Engine) Preview(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}

	var data map[string]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err.Error())
	}

	if data["table_name"] == "" {
		panic("tableName is empty")
	}
}