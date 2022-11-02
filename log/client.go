package log

import (
	"bytes"
	"fmt"
	stlog "log"
	"net/http"

	"distributed/registry"
)

func SetClientLogger(serviceURL string, clientService registry.ServiceName) {
	stlog.SetPrefix(fmt.Sprintf("[%v] - ", clientService))
	stlog.SetFlags(0)
	stlog.SetOutput(&clientLogger{url: serviceURL})
}

type clientLogger struct {
	url string
}

func (cl clientLogger) Write(data []byte) (int, error) {
	b := bytes.NewBuffer([]byte(data))
	r, err := http.Post(cl.url+"/log", "text/plain", b)
	if err != nil{
		return 0, err
	}
	if r.StatusCode != http.StatusOK{
		return 0, fmt.Errorf("Failed to send message, Service responsed with code %v\n", r.StatusCode)
	}
	return len(data), nil
}
