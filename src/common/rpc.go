package common

import (
	"net/rpc"
	"log"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type Args struct{
	Address string
}

type HttpHandlers struct {
	Path    string
	Handler func(http.ResponseWriter, *http.Request)
}

type SetConfigArg struct {
	Idx int
	PeerTuple []string
}
type PutValueArg struct{
	Key string
	Value string
	Idx int
}
type GetValueArg struct{
	Key string
	Idx int
}

func StartRpcServer(rcvr interface{}, port int, handlers []HttpHandlers , kind string){
	fmt.Printf("starting %s on port %d\n", kind, port)
	r := mux.NewRouter()
	rpcServer := rpc.NewServer()
	rpcServer.Register(rcvr)
	r.Handle("/rpc", rpcServer)
	if(handlers!=nil) {
		for _, handler := range handlers {
			r.HandleFunc(handler.Path, handler.Handler)
		}
	}
	address := fmt.Sprintf(":%d", port)
	go http.ListenAndServe(address, r)
	fmt.Printf("%s server started\n", kind)
}

func StartRpcCli(location string, kind string) *rpc.Client {
	fmt.Printf("starting client %s for srv location %s\n" ,kind, location )
	conn, err := rpc.DialHTTPPath("tcp", fmt.Sprintf("127.0.0.1:%s", location),"/rpc")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	fmt.Print("client started listen'g\n")
	return conn
}