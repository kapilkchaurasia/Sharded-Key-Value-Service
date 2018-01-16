package master

import (
	"fmt"
	"DistributedDatabase/src/common"
	"net/http"
	"log"
	"strconv"
)

var WorkerSrvAddr = make([]string, 0)
type MasterOperation int

func StartMasterSrv(port int) string{ //8080
	op := new(MasterOperation)
	common.StartRpcServer(op,port,HTTPHandler() ,"master-server")
	return fmt.Sprintf("%d",port)
}

func (m *MasterOperation) Registry(arg *common.Args, reply *string) error {
	fmt.Printf("Registing worker server %s\n", arg.Address)
	WorkerSrvAddr = append(WorkerSrvAddr, arg.Address)
	return nil
}


func HTTPHandler() []common.HttpHandlers{
	var handlers []common.HttpHandlers
	rfHandlerfunc :=  func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("inside RF handler\n")
		rf, ok := r.URL.Query()["rf"]
		if !ok{
			log.Println("Url Param RF is missing")
			return
		}
		rfValue,_ :=strconv.Atoi(rf[0])
		configs := setConfiguration(rfValue)
		setConfig(configs)
	}

	putHandlerfunc :=  func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("inside putValue handler\n")
		key, ok := r.URL.Query()["key"]
		if !ok{
			log.Println("Url Param Key is missing")
			return
		}
		value, ok := r.URL.Query()["value"]
		if !ok{
			log.Println("Url Param value is missing")
			return
		}
		putValue(key[0],value[0])
	}

	getHandlerfunc :=  func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("inside getValue handler\n")
		key, ok := r.URL.Query()["key"]
		if !ok{
			log.Println("Url Param key is missing")
			return
		}
		getValue(key[0])
	}

	return append(handlers, common.HttpHandlers{"/RF",rfHandlerfunc},
	common.HttpHandlers{"/putValue",putHandlerfunc},common.HttpHandlers{"/getValue",getHandlerfunc} )
}

func setConfiguration(rf int) [][]string{
	var rfConfig [][]string
 	for i,_:= range WorkerSrvAddr{
 		var rfList []string
	 	for j := 0; j <= rf -1 ; j++ {
        	  rfList = append(rfList,WorkerSrvAddr[(i+j)%len(WorkerSrvAddr)])
	 	}
	 	rfConfig = append(rfConfig,rfList)
 	}
 	return rfConfig
}


func hashingFunc(key string) int {
	asciiStr := []rune(key)
	var summnation int
	for _,ascii:= range asciiStr{
		summnation = summnation + int(ascii - '0')
	}
	return summnation%len(WorkerSrvAddr)
}