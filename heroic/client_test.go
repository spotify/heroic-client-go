package heroic

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestNewRequest(t *testing.T) {
	heroicUrl, _ := url.Parse("http://localhost/")
	clientId := "test-client-id"
	client := NewClient(heroicUrl, nil, &clientId)
	request, _ := client.NewRequest("POST", "test/path", []interface{}{})
	if request.Method != "POST" ||
		request.Host != "localhost" ||
		request.URL.Path != "/test/path" ||
		request.Header.Values("Content-Type")[0] != "application/json" ||
		request.Header.Values("X-Client-Id")[0] != clientId {
		log.Fatal("Unexpected request")
	}
}

func TestDoRequest(t *testing.T) {
	port, closeF := createHeroicServer(func(r *http.Request) {}, 1)
	defer closeF()

	heroicUrl, _ := url.Parse("http://localhost:" + port + "/")
	clientId := "test-client-id"
	client := NewClient(heroicUrl, nil, &clientId)

	request, _ := client.NewRequest("POST", "test/path", []interface{}{})
	response, err := client.Do(context.Background(), request, nil)
	if err != nil || response.StatusCode != 200 {
		log.Fatal("Error while executing request", err)
	}
}

func TestQueryMetrics(t *testing.T) {
	clientId := "test-client-id"
	port, closeF := createHeroicServer(func(r *http.Request) {
		if r.Method != "POST" {
			log.Fatalf("Unexpected request method %s != %s", r.Method, "POST")
		}
		if r.RequestURI != "/query/metrics" {
			log.Fatalf("Unexpected client id %s != %s", r.RequestURI, "/query/metrics")
		}
		if r.Header.Values("X-Client-Id")[0] != clientId {
			log.Fatalf("Unexpected client id %s != %s",
				r.Header.Values("X-Client-Id")[0], clientId)
		}
	}, 10, 5, 2, 10)
	defer closeF()

	heroicUrl, _ := url.Parse("http://localhost:" + port + "/")
	client := NewClient(heroicUrl, nil, &clientId)
	metrics, err := client.QueryMetrics(context.Background(), &QueryMetricsRequest{})
	if err != nil {
		log.Fatal("Error while querying for metrics", err)
	}

	if len(metrics.Errors) > 0 {
		log.Fatal("Error while querying for metrics", metrics.Errors)
	}

	if metrics.Result[0].Values[0][1].(float64) != 10 ||
		metrics.Result[0].Values[1][1].(float64) != 5 ||
		metrics.Result[0].Values[2][1].(float64) != 2 ||
		metrics.Result[0].Values[3][1].(float64) != 10 {
		log.Fatalf("Unexpected value returned %s", metrics.Result[0].Values)
	}
}

func createHeroicServer(
	assertRequest func(r *http.Request), responseValues ...float64,
) (string, func()) {
	var vals = make([][]interface{}, len(responseValues))
	for i, v := range responseValues {
		vals[i] = []interface{}{float64(i), v}
	}
	heroicHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertRequest(r)
		heroicResponse := QueryMetricsResponse{
			Range:      AbsoluteTimeRange{},
			Errors:     nil,
			Result:     []ShardedResultGroup{{Values: vals}},
			Statistics: Statistics{},
		}
		hpaEncoded, _ := json.Marshal(heroicResponse)
		w.Header().Add("Content-Type", "application/json")
		_, _ = w.Write(hpaEncoded)
	})

	metricsServer := httptest.NewServer(heroicHandler)
	splitUrl := strings.Split(metricsServer.URL, ":")
	metricsPort := splitUrl[len(splitUrl)-1]
	closeFunc := metricsServer.Close
	return metricsPort, closeFunc
}
