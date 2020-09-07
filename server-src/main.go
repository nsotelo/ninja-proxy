// This is a proxy service that lets you generate temporary URLs for giving third-parties controlled access.
// Unlike other solutions this is completely stateless, reducing complexity and improving scalability.
package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	b64 "encoding/base64"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
	"strings"
	"time"
)

const BUFFER_SIZE = 4096
const TIMEOUT = 20 * time.Second

func main() {
	var encKey string
	var host string
	var port string

	flag.StringVar(&encKey, "key", "", "The AES key used to decrypt requests. Should be 16, 24 or 32 bytes encoded in base64.")
	flag.StringVar(&host, "host", "", "The host to bind this service to. Defaults to 0.0.0.0")
	flag.StringVar(&port, "port", "7777", "The port to bind this service to. Defaults to 7777")
	flag.Parse()
	key, _ := base64.StdEncoding.DecodeString(encKey)
	checkKey(key)

	service := host + ":" + port
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Warning: %s", err)
			continue
		}
		go handleClient([]byte(key), conn)
	}
}

// Get the URL and expiry time from the proxy authorisation header
func unpackRequest(key []byte, req http.Request, client net.Conn) (upstream string, modReq http.Request) {
	normalAuth, isNormalAuth := req.Header["Authorization"]
	proxyAuth, isProxyAuth := req.Header["Proxy-Authorization"]

	var authData string
	var authHeader string
	if isNormalAuth {
		authData = normalAuth[0]
		authHeader = "Authorization"
	} else if isProxyAuth {
		authData = proxyAuth[0]
		authHeader = "Proxy-Authorization"
	} else {
		realm := []string{"Basic realm: \"Access to this resource.\""}
		handleConnectFailure(
			req,
			client,
			"Authorization required.",
			"401 Unauthorized",
			map[string][]string{"WWW-Authenticate": realm})
	}

	decoded, _ := b64.URLEncoding.DecodeString(strings.Replace(authData, "Basic ", "", 1))
	splitString := strings.Split(string(decoded), ":")
	nonce, ciphertext := splitString[0], splitString[1]
	u, expiry, headers := decrypt(key, nonce, ciphertext)
	if time.Now().After(expiry) {
		handleConnectFailure(
			req,
			client,
			"This link has expired.",
			"400 Bad Request",
			map[string][]string{},
		)
	}

	username := u.User.Username()
	password, _ := u.User.Password()
	u.User = nil
	delete(req.Header, authHeader)

	if username != "" && password != "" {
		auth := username + ":" + password
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
		req.Header[authHeader] = []string{basicAuth}
	}

	if isNormalAuth {
		req.Host = u.Host
	}

	for key, value := range headers {
		req.Header[key] = []string{value}
	}

	return u.Host, req

}

// Parse and decrypt the proxy authorisation header, then manage a proxy channel
func handleClient(key []byte, client net.Conn) {
	defer trapError()
	defer client.Close()
	upstream := handleConnect(key, client)
	defer upstream.Close()

	request := make(chan []byte)
	response := make(chan []byte)
	go listen(client, request)
	go listen(upstream, response)
	for {
		select {
		case sent := <-request:
			isClosed := transfer(upstream, sent)
			if isClosed {
				break
			}
			go listen(client, request)
		case recv := <-response:
			isClosed := transfer(client, recv)
			if isClosed {
				break
			}
			go listen(upstream, response)
		}

	}
}

// Decrypt the upstream connection information, check expiry and open connection to upstream host
func handleConnect(key []byte, client net.Conn) net.Conn {
	reader := bufio.NewReader(client)
	req, err := http.ReadRequest(reader)
	checkError(err)

	upstreamHost, modReq := unpackRequest(key, *req, client)
	validReq, err := httputil.DumpRequest(&modReq, true)
	checkError(err)

	dialer := net.Dialer{Timeout: TIMEOUT}
	conn, err := dialer.Dial("tcp", upstreamHost)
	checkError(err)
	transfer(conn, validReq)
	return conn
}

// Kill the connection when something goes wrong
func handleConnectFailure(req http.Request, client net.Conn, body string, status string, header map[string][]string) {
	statusCode, err := strconv.Atoi(strings.Split(status, " ")[0])
	checkError(err)
	response := &http.Response{
		Status:        status,
		StatusCode:    statusCode,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
		ContentLength: int64(len(body)),
		Request:       &req,
		Header:        make(http.Header, 0),
	}
	rawResponse, err := httputil.DumpResponse(response, true)
	checkError(err)
	client.Write(rawResponse)
	panic(body)
}

// Wait for new data from a socket.
func listen(conn net.Conn, data chan []byte) {
	buf := make([]byte, BUFFER_SIZE)
	read_len, err := conn.Read(buf)
	if err != nil {
		data <- make([]byte, 0)
	} else {
		data <- buf[:read_len]
	}
}

// Transfer data from one scoket to another and report if the sending socket is closed.
func transfer(conn net.Conn, data []byte) bool {
	if len(data) == 0 {
		return true
	} else {
		conn.Write(data)
		return false
	}
}
