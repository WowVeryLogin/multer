package multer

import (
	"context"
	"fmt"
	"github/WowVeryLogin/multer/pkg/httpclient"
	"github/WowVeryLogin/multer/pkg/pool"
	"io/ioutil"
	"sync"
)

type Request struct {
	Urls []string `json:"urls"`
}

type Response struct {
	Results [][]byte `json:"data"`
}

type Multer struct {
	pool       pool.Pool
	httpClient httpclient.Client
}

func New(cl httpclient.Client, pool pool.Pool) Multer {
	return Multer{
		pool:       pool,
		httpClient: cl,
	}
}

func (m *Multer) Close() {
	m.pool.Close()
}

func (m *Multer) HandleRequest(
	ctx context.Context,
	req Request,
) (*Response, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var (
		mtx  sync.Mutex
		resp Response
	)

	for _, url := range req.Urls {
		url := url
		m.pool.Put(ctx, pool.Task(func(ctx context.Context) error {
			res, err := m.httpClient.Get(ctx, url)
			if err != nil {
				return fmt.Errorf("making url request to %s: %w", url, err)
			}
			defer res.Body.Close()
			data, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return fmt.Errorf("reading response body: %w", err)
			}
			mtx.Lock()
			resp.Results = append(resp.Results, data)
			mtx.Unlock()

			return nil
		}))
	}

	var (
		err  error
		once sync.Once
	)
	m.pool.WaitBarrier(func(e error) {
		cancel()
		once.Do(func() { err = e })
	})

	if err != nil {
		return nil, fmt.Errorf("pool worker: %w", err)
	}
	return &resp, nil
}
