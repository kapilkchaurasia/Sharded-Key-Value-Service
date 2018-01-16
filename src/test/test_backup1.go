package test

import (
	"github.com/peterbourgon/raft"
	"bytes"
	"net/http"
	"sync"
	"net/url"
	"os"
	"bufio"
	"fmt"
)

func main() {
	var wg sync.WaitGroup
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

	wg.Add(1)
	ponger := func(uint64, []byte) []byte { return []byte(`PONG`) }
	// Construct the server
	s := raft.NewServer(2, &bytes.Buffer{}, ponger)

	// Expose the server using a HTTP transport
	raft.HTTPTransport(http.DefaultServeMux, s)

	http.HandleFunc("/getValue", func(w http.ResponseWriter, r *http.Request) {
		fo, err := os.Open("/Users/k0c00nc/output.txt")
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(fo)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		fo.Close()
	})

	go http.ListenAndServe(":8081", nil)
	s.SetConfiguration(
		mustNewHTTPPeer(mustParseURL("http://127.0.0.1:8081")),
	)
	s.Start()
wg.Wait()
}
