package test

import (
	"github.com/peterbourgon/raft"
	"net/url"
	"bytes"
	"net/http"
	"fmt"
	"time"
	"os"
	"bufio"
)

var s *raft.Server
func main() {
	// A no-op ApplyFunc
	applyValue := func(value string) func(uint64, []byte) []byte {
		return func(index uint64, cmd []byte) []byte {
			fo, err := os.Create("/Users/k0c00nc/output.txt")
			if err != nil {
				panic(err)
			}
			fo.Write(cmd)
			fo.Close()
			return []byte(`PONG`)
		}
	}
	// Helper function to parse URLs
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
	// Construct the server
	s := raft.NewServer(1, &bytes.Buffer{}, applyValue("key"))

	// Expose the server using a HTTP transport
	raft.HTTPTransport(http.DefaultServeMux, s)

	http.HandleFunc("/putValue", func(w http.ResponseWriter, r *http.Request) {
		response := make(chan []byte)
		if err := s.Command([]byte(`PUT:`), response); err != nil {
			panic(err) // command not accepted
		}
		// After the command is replicated, we'll receive the response
		fmt.Printf("-----------------%s------------------\n", <-response)
	})

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

	go http.ListenAndServe(":8080", nil)

	// Set the initial server configuration
	s.SetConfiguration(
		mustNewHTTPPeer(mustParseURL("http://127.0.0.1:8080")), // this server
		mustNewHTTPPeer(mustParseURL("http://127.0.0.1:8081")),
		mustNewHTTPPeer(mustParseURL("http://127.0.0.1:8082")),
	)

	// Start the server
	s.Start()
	time.Sleep(time.Minute*10)
}
