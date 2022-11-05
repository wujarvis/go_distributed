package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
)

func RegisterService(r Registration) error {
	heartbeatURL, err := url.Parse(r.HeartbeatURL)
	if err != nil {
		return err
	}
	http.HandleFunc(heartbeatURL.Path, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	serviceUpdateURL, err := url.Parse(r.ServiceUpdateURL)
	if err != nil {
		return err
	}
	http.Handle(serviceUpdateURL.Path, &serviceUpdateHandler{})
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err = enc.Encode(r)
	if err != nil {
		return err
	}
	res, err := http.Post(ServicesURL, "application/json", buf)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to register service. Registry service responsed with code %v\n", res.StatusCode)
	}
	return nil
}

type serviceUpdateHandler struct{}

func (suh serviceUpdateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	dec := json.NewDecoder(r.Body)
	var p patch
	err := dec.Decode(&p)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Printf("Updated received [Added url]: %v, [Removed url]: %v\n", p.Added, p.Removed)
	prov.Update(p)
}

func ShutdownService(url string) error {
	req, err := http.NewRequest(http.MethodDelete, ServicesURL, bytes.NewBuffer([]byte(url)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to deregister service, Registry service responsed with code %v", res.StatusCode)
	}
	return nil
}

type providers struct {
	services map[ServiceName][]string
	mutex    *sync.RWMutex
}

var prov = providers{
	services: make(map[ServiceName][]string),
	mutex:    new(sync.RWMutex),
}

func (p *providers) Update(pat patch) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, patchEntry := range pat.Added {
		if _, ok := p.services[patchEntry.Name]; !ok {
			p.services[patchEntry.Name] = make([]string, 0)
		}
		p.services[patchEntry.Name] = append(p.services[patchEntry.Name], patchEntry.URL)
	}

	for _, patchRemove := range pat.Removed {
		if providerURLs, ok := p.services[patchRemove.Name]; ok {
			for i := range providerURLs {
				if providerURLs[i] == patchRemove.URL {
					p.services[patchRemove.Name] = append(providerURLs[:i], providerURLs[i+1:]...)
				}
			}
		}
	}
}

func (p *providers) get(name ServiceName) (string, error) {
	providers, ok := p.services[name]
	if !ok {
		return "", fmt.Errorf("No providers available for service %v", name)
	}
	idx := int(rand.Float32() * float32(len(providers)))
	return providers[idx], nil
}

func GetProvider(name ServiceName) (string, error) {
	return prov.get(name)
}
