package lib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
	"testing"
	"time"
)

var (
	addr string = "127.0.0.1:6111"
	initOnce sync.Once = sync.Once{}
	serverOnce sync.Once = sync.Once{}
)


type HttpConf struct {
	ServerAddr     string   `toml:"server_addr"`
	ReadTimeout    int      `toml:"read_timeout"`
	WriteTimeout   int      `toml:"write_timeout"`
	MaxHeaderBytes int      `toml:"max_header_bytes"`
	AllowHost      []string `toml:"allow_host"`
}

//获取 程序运行环境 dev prod
func Test_GetConfEnv(t *testing.T) {
	InitTest()
	fmt.Println(GetConfEnv())
	DestroyTest()
}

// 加载自定义配置文件
func Test_ParseLocalConfig(t *testing.T) {
	InitTest()
	httpProfile := &HttpConf{}
	err:=ParseLocalConfig("http.toml",httpProfile)
	if err!=nil{
		t.Fatal(err)
	}
	fmt.Println(httpProfile)
	DestroyTest()
}

//测试服务器解析
func TestParseServerAddr(t *testing.T) {
	serverAddr := "10.94.64.101:6379"

	host, port := ParseServerAddr(serverAddr)
	if host != "10.94.64.101" {
		t.Fatalf("parse failure, wanted %s, result %s", "10.94.64.101", host)
	}
	if port != "6379" {
		t.Fatalf("parse failure, wanted %s, result %s", "6379", port)
	}

	fmt.Printf("Host: %s\nPort: %s\n", host, port)
}

//测试PostJson请求
func TestJson(t *testing.T) {
	InitTestServer()
	//首次scrollsId不传递
	jsonStr := "{\"source\":\"control\",\"cityId\":\"12\",\"trailNum\":10,\"dayTime\":\"2018-11-21 16:08:00\",\"limit\":2,\"andOperations\":{\"cityId\":\"eq\",\"trailNum\":\"gt\",\"dayTime\":\"eq\"}}"
	url := "http://"+addr+"/json"
	_, res, err := HttpJSON(NewTrace(), url, jsonStr, 1000, nil)
	fmt.Println(string(res))
	if err != nil {
		fmt.Println(err.Error())
	}
}

//测试Get请求
func TestGet(t *testing.T) {
	InitTestServer()
	a := url.Values{
		"city_id": {"12"},
	}
	url := "http://"+addr+"/get"
	_, res, err := HttpGET(NewTrace(), url, a, 1000, nil)
	fmt.Println("city_id="+string(res))
	if err != nil {
		fmt.Println(err.Error())
	}
}

//测试Post请求
func TestPost(t *testing.T) {
	InitTestServer()
	a := url.Values{
		"city_id": {"12"},
	}
	url := "http://"+addr+"/post"
	_, res, err := HttpPOST(NewTrace(), url, a, 1000, nil, "")
	fmt.Println("city_id="+string(res))
	if err != nil {
		fmt.Println(err.Error())
	}
}

//初始化测试用例
func InitTest()  {
	initOnce.Do(func() {
		if err:=Init("../conf/dev/");err!=nil{
			log.Fatal(err)
		}
	})
}

//销毁测试用例
func DestroyTest()  {
	//Destroy()
}

//只运行一次服务器
func InitTestServer() {
	serverOnce.Do(func() {
		http.HandleFunc("/postjson", func(writer http.ResponseWriter, request *http.Request) {
			bodyBytes, _ := ioutil.ReadAll(request.Body)
			request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes)) // Write body back
			writer.Write([]byte(bodyBytes))
		})
		http.HandleFunc("/get", func(writer http.ResponseWriter, request *http.Request) {
			request.ParseForm()
			cityID := request.FormValue("city_id")
			writer.Write([]byte(cityID))
		})
		http.HandleFunc("/post", func(writer http.ResponseWriter, request *http.Request) {
			request.ParseForm()
			cityID := request.FormValue("city_id")
			writer.Write([]byte(cityID))
		})
		go func() {
			log.Println("ListenAndServe ", addr)
			err := http.ListenAndServe(addr, nil) //设置监听的端口
			if err != nil {
				log.Fatal("ListenAndServe: ", err)
			}
		}()
		time.Sleep(time.Second)
	})
}