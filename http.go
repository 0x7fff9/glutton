package glutton

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"strings"
)

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(request, "\n")
}

// HandleHTTP takes a net.Conn and does basic HTTP communication
func (g *Glutton) HandleHTTP(conn net.Conn) {
	defer conn.Close()
	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		g.logger.Errorf("[http    ] %v", err)
		return
	}
	g.logger.Infof("[http    ] %s", formatRequest(req))
	if req.ContentLength > 0 {
		defer req.Body.Close()
		buf := bytes.NewBuffer(make([]byte, 0, req.ContentLength))
		_, err = buf.ReadFrom(req.Body)
		if err != nil {
			g.logger.Errorf("[http    ] %v", err)
			return
		}
		body := buf.Bytes()
		g.logger.Infof("[http    ] http body:\n%s", hex.Dump(body[:]))
	}
	conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
}
