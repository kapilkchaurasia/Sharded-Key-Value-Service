package test

import (
	"github.com/peterbourgon/raft"
	"bytes"
	"net/http"
	"sync"
	"net/url"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	mustParseURL := func(rawurl string) *url.URL {
		u, err := url.Parse(rawurl)
		if err != nil {
			panic(err)
		}
		u.Path = ""
		return u
	}

	// Helper function to construct HTTP Peers
	mustNewHTTPPeer := func(u *url.URL) raft.Peer {
		p, err := raft.NewHTTPPeer(u)
		if err != nil {
			panic(err)
		}
		return p
	}

	ponger := func(uint64, []byte) []byte { return []byte(`PONG`) }
	// Construct the server
	s := raft.NewServer(3, &bytes.Buffer{}, ponger)

	// Expose the server using a HTTP transport
	raft.HTTPTransport(http.DefaultServeMux, s)
	go http.ListenAndServe(":8082", nil)
	s.SetConfiguration(
		mustNewHTTPPeer(mustParseURL("http://127.0.0.1:8081")),
		mustNewHTTPPeer(mustParseURL("http://127.0.0.1:8082")),
	)
	s.Start()
	wg.Wait()
}
