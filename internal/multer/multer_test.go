package multer

import (
	"bytes"
	"context"
	"fmt"
	"github/WowVeryLogin/multer/pkg/httpclient"
	"github/WowVeryLogin/multer/pkg/pool/flexibale"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMulter(t *testing.T) {
	t.Parallel()
	ts := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "Some data")
			},
		),
	)
	defer ts.Close()

	p := flexibale.New(&flexibale.Config{
		MaxWorkers: 2,
	})
	m := New(httpclient.New(&httpclient.Config{}), p)
	defer m.Close()
	resp, err := m.HandleRequest(context.Background(), Request{
		Urls: []string{
			ts.URL,
			ts.URL,
			ts.URL,
		},
	})
	if err != nil {
		t.Error(err)
	}

	expected := [][]byte{
		[]byte("Some data\n"),
		[]byte("Some data\n"),
		[]byte("Some data\n"),
	}
	for i := range resp.Results {
		if !bytes.Equal(resp.Results[i], expected[i]) {
			t.FailNow()
		}
	}
}
