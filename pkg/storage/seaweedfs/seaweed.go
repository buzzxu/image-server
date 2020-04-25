package seaweedfs

import (
	workerpool "github.com/buzzxu/boys/worker-pool"
	"net/http"
	"net/url"
)

type Seaweed struct {
	master    *url.URL
	filers    []*Filer
	chunkSize int64
	client    *httpClient
	workers   *workerpool.Pool
	token     string
}

func newSeaweed(masterURL string, filers []string, chunkSize int64, token string, client *http.Client) (c *Seaweed, err error) {
	u, err := parseURI(masterURL)
	if err != nil {
		return
	}

	c = &Seaweed{
		master:    u,
		client:    newHTTPClient(client),
		chunkSize: chunkSize,
	}
	if len(filers) > 0 {
		c.filers = make([]*Filer, 0, len(filers))
		for i := range filers {
			var filer *Filer
			filer, err = newFiler(filers[i], token, c.client)
			if err != nil {
				c.close()
				return
			}
			c.filers = append(c.filers, filer)
		}
	}

	// start underlying workers
	c.workers = createWorkerPool()
	c.workers.Start()

	return
}

func (c *Seaweed) close() (err error) {
	if c.workers != nil {
		c.workers.Stop()
	}
	if c.client != nil {
		err = c.client.Close()
	}

	if len(c.filers) > 0 {
		for i := 0; i < len(c.filers); i++ {
			c.filers[i].Close()
		}
	}
	return
}
