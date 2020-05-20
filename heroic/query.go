package heroic

import "context"

func (c *Client) QueryMetrics(ctx context.Context, query *QueryMetricsRequest) (*QueryMetricsResponse, error) {
	req, err := c.NewRequest("POST", "query/metrics", query)
	if err != nil {
		return nil, err
	}
	var qmr QueryMetricsResponse
	_, err = c.Do(ctx, req, &qmr)
	return &qmr, err
}

type QueryMetricsRequest struct {
	Range       TimeRange   `json:"range"`
	Filter      Filter      `json:"filter"`
	Aggregation Aggregation `json:"aggregation"`
}

type QueryMetricsResponse struct {
	Range      AbsoluteTimeRange    `json:"range"`
	Errors     []RequestError       `json:"errors"`
	Result     []ShardedResultGroup `json:"result"`
	Statistics Statistics           `json:"statistics"`
}

type RequestError struct {
	Type string `json:"type"`
}

type NodeError struct {
	Type   string            `json:"type"`
	Error  string            `json:"error"`
	NodeID string            `json:"nodeId"`
	Node   string            `json:"node"`
	Tags   map[string]string `json:"tags"`
}

type ShardError struct {
	Type  string            `json:"type"`
	Error string            `json:"error"`
	Nodes []string          `json:"nodes"`
	Shard map[string]string `json:"shard"`
}

type QueryError struct {
	Type  string `json:"type"`
	Error string `json:"error"`
}

type ShardedResultGroup struct {
	Type      string            `json:"type"`
	Hash      string            `json:"hash"`
	Shard     map[string]string `json:"shard"`
	Cadence   int64             `json:"cadence"`
	Values    [][]interface{}   `json:"values"`
	Tags      map[string]string `json:"tags"`
	TagCounts map[string]int    `json:"tagCounts"`
}

type Statistics struct {
	Counters map[string]int `json:"counters"`
}

type TimeRange interface{}

type Filter []interface{}

type Aggregation interface{}

type TimeUnit string

const (
	Milliseconds = TimeUnit("MILLISECONDS")
	Seconds      = TimeUnit("SECONDS")
	Minutes      = TimeUnit("MINUTES")
	Hours        = TimeUnit("HOURS")
	Days         = TimeUnit("DAYS")
	Weeks        = TimeUnit("WEEKS")
	Months       = TimeUnit("MONTHS")
)

type TimeRangeType string

const (
	Absolute = TimeRangeType("absolute")
	Relative = TimeRangeType("relative")
)

type AbsoluteTimeRange struct {
	Type  TimeRangeType `json:"type,omitempty"`
	Start int64         `json:"start"`
	End   int64         `json:"end"`
}

type RelativeTimeRange struct {
	Type  TimeRangeType `json:"type"`
	Unit  TimeUnit      `json:"unit"`
	Value int           `json:"value"`
}

type AggregationType string

const (
	// sampling
	Min     = AggregationType("min")
	Max     = AggregationType("max")
	Average = AggregationType("average")
	Sum     = AggregationType("sum")

	// chaining
	Chain = AggregationType("chain")

	// grouping
	Group = AggregationType("group")

	// filtering
	TopK    = AggregationType("topk")
	BottomK = AggregationType("bottomk")
	AboveK  = AggregationType("abovek")
	BelowK  = AggregationType("belowk")
)

type Sample struct {
	Unit   TimeUnit `json:"unit"`
	Value  int      `json:"value"`
	Size   int      `json:"size,omitempty"`
	Extent int      `json:"extent,omitempty"`
}

type SamplingAggregation struct {
	Type     AggregationType `json:"type"`
	Sampling Sample          `json:"sampling"`
}

type ChainingAggregation struct {
	Type  AggregationType `json:"type"`
	Chain []Aggregation   `json:"chain"`
}

type GroupingAggregation struct {
	Type AggregationType `json:"type"`
	Of   []string        `json:"of"`
	Each Aggregation     `json:"each"`
}

type FilteringAggregation struct {
	Type AggregationType `json:"type"`
	K    int             `json:"k"`
}
