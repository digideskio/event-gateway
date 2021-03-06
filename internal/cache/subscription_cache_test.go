package cache

import (
	"testing"

	"github.com/serverless/event-gateway/functions"
	"github.com/serverless/event-gateway/internal/cors"
	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

func TestSubscriptionCacheModified(t *testing.T) {
	scache := newSubscriptionCache(zap.NewNop())

	scache.Modified("testsub1", []byte(`{"subscriptionId":"testsub1", "event": "test.event", "functionId": "testfunc1", "path": "/"}`))
	scache.Modified("testsub2", []byte(`{"subscriptionId":"testsub2", "event": "test.event", "functionId": "testfunc2", "path": "/"}`))

	assert.Equal(
		t,
		[]functions.FunctionID{functions.FunctionID("testfunc1"), functions.FunctionID("testfunc2")},
		scache.eventToFunctions["/"]["test.event"],
	)
}

func TestSubscriptionCacheModifiedHTTPSubscription(t *testing.T) {
	scache := newSubscriptionCache(zap.NewNop())

	scache.Modified("testsub1", []byte(`{"subscriptionId":"testsub1", "event": "http", "functionId": "testfunc1", "path": "/a", "method": "GET"}`))
	scache.Modified("testsub2", []byte(`{"subscriptionId":"testsub2", "event": "http", "functionId": "testfunc2", "path": "/b", "method": "GET"}`))

	id, _, _ := scache.endpoints["GET"].Resolve("/a")
	assert.Equal(t, functions.FunctionID("testfunc1"), *id)
}

func TestSubscriptionCacheModifiedCORSConfiguration(t *testing.T) {
	scache := newSubscriptionCache(zap.NewNop())

	scache.Modified("testsub1", []byte(`{
		"subscriptionId":"testsub1",
		"event": "http",
		"functionId": "testfunc1",
		"path": "/a",
		"method": "GET",
		"cors": {
			"origins": ["http://example.com"]
		}
	}`))

	_, _, corsConfig := scache.endpoints["GET"].Resolve("/a")
	assert.Equal(t, &cors.CORS{Origins: []string{"http://example.com"}}, corsConfig)
}

func TestSubscriptionCacheModified_WrongPayload(t *testing.T) {
	scache := newSubscriptionCache(zap.NewNop())

	scache.Modified("testsub", []byte(`not json`))

	assert.Equal(t, []functions.FunctionID(nil), scache.eventToFunctions["/"]["test.event"])
}

func TestSubscriptionCacheModifiedDeleted(t *testing.T) {
	scache := newSubscriptionCache(zap.NewNop())

	scache.Modified("testsub1", []byte(`{"subscriptionId":"testsub1", "event": "test.event", "functionId": "testfunc1", "path": "/"}`))
	scache.Modified("testsub2", []byte(`{"subscriptionId":"testsub2", "event": "test.event", "functionId": "testfunc2", "path": "/"}`))
	scache.Deleted("testsub1", []byte(`{"subscriptionId":"testsub1", "event": "test.event", "functionId": "testfunc1", "path": "/"}`))

	assert.Equal(t, []functions.FunctionID{functions.FunctionID("testfunc2")}, scache.eventToFunctions["/"]["test.event"])
}

func TestSubscriptionCacheModifiedDeletedHTTPSubscription(t *testing.T) {
	scache := newSubscriptionCache(zap.NewNop())

	scache.Modified("testsub1", []byte(`{"subscriptionId":"testsub1", "event": "http", "functionId": "testfunc1", "path": "/", "method": "GET"}`))
	scache.Deleted("testsub1", []byte(`{"subscriptionId":"testsub1", "event": "http", "functionId": "testfunc1", "path": "/", "method": "GET"}`))

	id, _, _ := scache.endpoints["GET"].Resolve("/")
	assert.Nil(t, id)
}

func TestSubscriptionCacheModifiedDeletedLast(t *testing.T) {
	scache := newSubscriptionCache(zap.NewNop())

	scache.Modified("testsub", []byte(`{"subscriptionId":"testsub", "event": "test.event", "functionId": "testfunc", "path": "/"}`))
	scache.Deleted("testsub", []byte(`{"subscriptionId":"testsub", "event": "test.event", "functionId": "testfunc", "path": "/"}`))

	assert.Equal(t, []functions.FunctionID(nil), scache.eventToFunctions["/"]["test.event"])
}
