package resolver

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/eddyzags/resolver/marathon"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/unrolled/render"
	"google.golang.org/grpc/naming"
)

func TestPollInstantiationInstantiationWithoutError(t *testing.T) {
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

	marathonClient := marathon.NewClient(&marathon.Config{
		URI: ts.URL,
	})

	poller, err := newPoll(val, marathonClient)
	assert.NoError(err, "an unexpected error occured in poller instantiation")

	assert.Equal(int64(2), poller.portIndex, "the port index should be equals")
	assert.Equal(apps[0].ID, poller.appID, "the app id should be equals")
	assert.Equal(val, poller.label, "the label should be equals")
}

func TestPollInstantiationWithErrorOnGetApplications(t *testing.T) {
	assert := assert.New(t)

	val := "service-test"

	handler := http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
		render.New().JSON(rw, http.StatusBadRequest, nil)
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	marathonClient := marathon.NewClient(&marathon.Config{
		URI: ts.URL,
	})

	poller, err := newPoll(val, marathonClient)
	assert.Error(err, "an error was expected in poller instantiation")

	assert.Nil(poller, "the poller should be nil")
}

func TestPollInstantiationWithErrorOnAppDuplicate(t *testing.T) {
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
		{
			ID: "/test-2",
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

	marathonClient := marathon.NewClient(&marathon.Config{
		URI: ts.URL,
	})

	poller, err := newPoll(val, marathonClient)
	assert.Error(err, "an error was expected in poller instantiation")

	assert.Nil(poller, "the poller should be nil")
}

func TestPollInstantiationWithErrorOnAppNotFound(t *testing.T) {
	assert := assert.New(t)

	val := "service-test"

	handler := http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
		render.New().JSON(rw, http.StatusOK, []*marathon.Application{})
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	marathonClient := marathon.NewClient(&marathon.Config{
		URI: ts.URL,
	})

	poller, err := newPoll(val, marathonClient)
	assert.Error(err, "an error was expected in poller instantiation")

	assert.Nil(poller, "the poller should be nil")
}

func TestPollInstantiationWithErrorOnLabelSyntax(t *testing.T) {
	assert := assert.New(t)

	val := "service-test"

	apps := []*marathon.Application{
		{
			ID: "/test",
			Labels: &map[string]string{
				"": val,
			},
		},
	}

	handler := http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
		render.New().JSON(rw, http.StatusOK, apps)
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	marathonClient := marathon.NewClient(&marathon.Config{
		URI: ts.URL,
	})

	poller, err := newPoll(val, marathonClient)
	assert.Error(err, "an error was expected in poller instantiation")

	assert.Nil(poller, "the poller should be nil")
}

func TestPollInstantiationWithErrorOnLabelPortIndexSyntax(t *testing.T) {
	assert := assert.New(t)

	key := "RESOLVER_N_NAME"
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

	marathonClient := marathon.NewClient(&marathon.Config{
		URI: ts.URL,
	})

	poller, err := newPoll(val, marathonClient)
	assert.Error(err, "an error was expected in poller instantiation")

	assert.Nil(poller, "the poller should be nil")
}

func TestPollInstantiationWithErrorOnLabelPortIndexSyntax2(t *testing.T) {
	assert := assert.New(t)

	key := "RESOLVER_-1_NAME"
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
		_ = render.New().JSON(rw, http.StatusOK, apps)
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	marathonClient := marathon.NewClient(&marathon.Config{
		URI: ts.URL,
	})

	poller, err := newPoll(val, marathonClient)
	assert.Error(err, "an error was expected in poller instantiation")

	assert.Nil(poller, "the poller should be nil")
}

func TestPollNextAddWithoutError(t *testing.T) {
	assert := assert.New(t)

	grpcServer, addr, err := newGRPCServer()
	assert.NoError(err, "an unexpected error occured in grpc server instantiation")

	key := "RESOLVER_1_NAME"
	val := "service-test"

	port, err := strconv.ParseInt(strings.Split(addr, ":")[1], 10, 32)
	assert.NoError(err, "an unexpected error occured in grpc server port parsing")

	apps := []*marathon.Application{
		{
			ID: "/test",
			Labels: &map[string]string{
				key: val,
			},
		},
	}

	tasks := []*marathon.Task{
		{
			ID:    uuid.Must(uuid.NewV4()).String(),
			AppID: apps[0].ID,
			Host:  "127.0.0.1",
			Ports: []int{
				2221,
				int(port),
				2223,
			},
		},
	}

	handler := http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
		switch rq.URL.Path {
		case "/v2/apps":
			render.New().JSON(rw, http.StatusOK, apps)
		case "/v2/tasks":
			render.New().JSON(rw, http.StatusOK, tasks)
		}
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	marathonClient := marathon.NewClient(&marathon.Config{
		URI: ts.URL,
	})

	poller, err := newPoll(val, marathonClient)
	assert.NoError(err, "an unexpected error occured in poller instantiation")

	poller.run()

	ups, err := poller.Next()
	assert.NoError(err, "an unexpected error occured in poller next")

	assert.Equal(1, len(ups), "The number of updates should be 1")
	assert.Equal(naming.Add, ups[0].Op, "the operations should be equals")
	assert.Equal(addr, ups[0].Addr, "the app id should be equals")

	grpcServer.Stop()
	poller.Close()
}

func TestPollNextDeleteWithoutError(t *testing.T) {
	assert := assert.New(t)

	grpcServer, addr, err := newGRPCServer()
	assert.NoError(err, "an unexpected error occured in grpc server instantiation")

	defer grpcServer.Stop()

	key := "RESOLVER_1_NAME"
	val := "service-test"

	port, err := strconv.ParseInt(strings.Split(addr, ":")[1], 10, 32)
	assert.NoError(err, "an unexpected error occured in grpc server port parsing")

	apps := []*marathon.Application{
		{
			ID: "/test",
			Labels: &map[string]string{
				key: val,
			},
		},
	}

	tasks := []*marathon.Task{
		{
			ID:    uuid.Must(uuid.NewV4()).String(),
			AppID: apps[0].ID,
			Host:  "127.0.0.1",
			Ports: []int{
				2221,
				int(port),
				2223,
			},
		},
	}

	handler := http.HandlerFunc(func(rw http.ResponseWriter, rq *http.Request) {
		switch rq.URL.Path {
		case "/v2/apps":
			render.New().JSON(rw, http.StatusOK, apps)
		case "/v2/tasks":
			render.New().JSON(rw, http.StatusOK, tasks)
		}
	})

	ts := httptest.NewServer(handler)
	defer ts.Close()

	marathonClient := marathon.NewClient(&marathon.Config{
		URI: ts.URL,
	})

	poller, err := newPoll(val, marathonClient)
	assert.NoError(err, "an unexpected error occured in poller instantiation")
	defer poller.Close()

	go poller.run()

	_, err = poller.Next()
	assert.NoError(err, "an unexpected error occured in poller next")

	grpcServer.Stop()

	ups, err := poller.Next()
	assert.NoError(err, "an unexpected error occured in poller next")

	assert.Equal(1, len(ups), "The number of updates should be 1")
	assert.Equal(naming.Delete, ups[0].Op, "the operations should be equals")
	assert.Equal(addr, ups[0].Addr, "the app id should be equals")
}
