package build

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/gopacket"
)

const (
	Port    = 80
	Version = "0.1"
)

const (
	CmdPort = "-p"
	CmdURL  = "-u"
)

type H struct {
	port    int
	version string
	url     string
}

var hp *H

func NewInstance() *H {
	if hp == nil {
		hp = &H{
			port:    Port,
			version: Version,
		}
	}
	return hp
}

func (m *H) ResolveStream(net, transport gopacket.Flow, buf io.Reader) {

	bio := bufio.NewReader(buf)
	buff := bytes.NewBufferString("")
	for {
		req, err := http.ReadRequest(bio)
		buff.Reset()

		if err == io.EOF {
			return
		} else if err != nil {
			continue
		} else {

			url := req.URL.String()
			if len(m.url) > 0 && !strings.Contains(url, m.url) {
				req.Body.Close()
				continue
			}

			_, _ = buff.WriteString(fmt.Sprintf("[%v] ", req.Method))
			_, _ = buff.WriteString(fmt.Sprintf("[%v%v] ", req.Host, url))

			_ = buff.WriteByte('[')
			for k, v := range req.Header {
				_, _ = buff.WriteString(fmt.Sprintf("%v=%v,", k, v[0]))
			}
			_ = buff.WriteByte(']')

			_ = buff.WriteByte('[')
			_, _ = buff.ReadFrom(req.Body)
			_ = buff.WriteByte(']')

			log.Println(buff.String())

			req.Body.Close()
		}
	}
}

func (m *H) BPFFilter() string {
	return "tcp and port " + strconv.Itoa(m.port)
}

func (m *H) Version() string {
	return Version
}

func (m *H) SetFlag(flg []string) {

	c := len(flg)
	if c == 0 {
		return
	}
	if c>>1 == 0 {
		fmt.Println("ERR : Http Number of parameters")
		os.Exit(1)
	}

	for i := 0; i < c; i = i + 2 {
		key := flg[i]
		val := flg[i+1]

		switch key {
		case CmdPort:
			port, err := strconv.Atoi(val)
			m.port = port
			if err != nil {
				panic("ERR : port")
			}
			if port < 0 || port > 65535 {
				panic("ERR : port(0-65535)")
			}
		case CmdURL:
			m.url = strings.TrimSpace(val)
			if len(m.url) == 0 {
				panic("ERR : url(no space)")
			}
		default:
			panic("ERR : http's params")
		}
	}
}
