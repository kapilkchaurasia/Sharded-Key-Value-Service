package worker

import (
	"fmt"
	"net/url"
	"github.com/peterbourgon/raft"
	"bytes"
	"net/http"
	"strconv"
	"strings"
)

var magicNo int = 1000
// Helper function to parse URLs
var mustParseURL = func(rawurl string) *url.URL {
	u, err := url.Parse(rawurl)
	if err != nil {
		panic(err)
	}
	u.Path = ""
	return u
}

// Helper function to construct HTTP Peers
var mustNewHTTPPeer = func(u *url.URL) raft.Peer {
	p, err := raft.NewHTTPPeer(u)
	if err != nil {
		panic(err)
	}
	return p
}

func RaftServer(raftId uint64, port int) *raft.Server{
	var applyValue = func(value string) func(uint64, []byte) []byte {
		return func(index uint64, cmd []byte) []byte {
			inputCmd := string(cmd[:])
			fmt.Printf("cmd is %s", inputCmd)
			switch cmdType := strings.Split(inputCmd, ":")[0]; cmdType {
			case "PUT":
				writeToDb(strings.Split(inputCmd, ":")[1])
				return []byte("success")
			case "GET":
				return []byte(fetchFromDb(strings.Split(inputCmd, ":")[1]))

			}
			return []byte("noOp")
		}
	}
	fmt.Printf("raftId is %d on port %d",raftId, port+magicNo)
	s := raft.NewServer(raftId, &bytes.Buffer{}, applyValue("key"))
	r := http.NewServeMux()
	raft.HTTPTransport(r, s)
	go http.ListenAndServe(fmt.Sprintf(":%d",port+magicNo), r)
	return s
}


func SetConfig(s *raft.Server,rawPeers[]string){
	var peers []raft.Peer
	for _,peer := range rawPeers{
		port,_ :=strconv.Atoi(peer)
		peers = append(peers, mustNewHTTPPeer(mustParseURL(fmt.Sprintf("http://127.0.0.1:%d",port+magicNo)))) //"http://127.0.0.1:9081"
	}
	s.SetConfiguration(peers...)
	s.Start()
}

func writeToDb(kv string) {

}

func fetchFromDb(key string) string {
return ""
}