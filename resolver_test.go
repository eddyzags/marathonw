package resolver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/eddyzags/resolver/marathon"

	"github.com/stretchr/testify/assert"
	"github.com/unrolled/render"
)

func TestResolverInstantiationWithoutError(t *testing.T) {
	assert := assert.New(t)

	handler := http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
		render.New().JSON(rw, http.StatusOK, nil)
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	resolver, err := New(ts.URL)
	assert.NoError(err, "an unexpected error occured in resolver instantiation")

	assert.NotNil(resolver, "resolver shouldn't be nil")
}

func TestResolverInstantiationWithErrorOnMarathon(t *testing.T) {
	assert := assert.New(t)

	resolver, err := New("test-123")
	assert.Error(err, "an unexpected error occured in resolver instantiation")

	assert.Nil(resolver, "resolver should be nil")
}

func TestResolveWithoutError(t *testing.T) {
	assert := assert.New(t)

	key := "RESOLVER_2_NAME"
	val := "service-test"

	apps := []*marathon.Application{
		{
			ID: "/test",
			Labels: &map[string]string{
				key: val,
			},
		},
	}

	handler := http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
		render.New().JSON(rw, http.StatusOK, apps)
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	resolver, err := New(ts.URL)
	assert.NoError(err, "an unexpected error occured in resolver instantiation")

	watcher, err := resolver.Resolve(val)
	assert.NoError(err, "an unexpected errororccured in resolve")

	watcher.Close()
}

func TestResolveWithErrorOnLabelNotFound(t *testing.T) {
	assert := assert.New(t)

	val := "service-test"

	handler := http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
		render.New().JSON(rw, http.StatusOK, []*marathon.Application{})
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	resolver, err := New(ts.URL)
	assert.NoError(err, "an unexpected error occured in resolver instantiation")

	watcher, err := resolver.Resolve(val)
	assert.Error(err, "an unexpected errororccured in resolve")

	assert.Nil(watcher, "the watcher should be nil")
}
