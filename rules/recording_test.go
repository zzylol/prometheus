// Copyright 2013 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rules

import (
	// "context"
	// "html/template"
	"testing"
	"time"
	"path/filepath"
	"os"

	"github.com/stretchr/testify/require"

	"github.com/prometheus/prometheus/model/labels"
	// "github.com/prometheus/prometheus/model/timestamp"
	"github.com/prometheus/prometheus/promql"
	"github.com/prometheus/prometheus/promql/parser"
	// "github.com/prometheus/prometheus/util/teststorage"
	"github.com/prometheus/prometheus/util/testutil"
)

/*
func TestRuleEval(t *testing.T) {
	t.Log("Hi! in TestRuleEval!")	

	storage := teststorage.New(t)
	defer storage.Close()

	opts := promql.EngineOpts{
		Logger:     nil,
		Reg:        nil,
		MaxSamples: 10,
		Timeout:    10 * time.Second,
	}

	engine := promql.NewEngine(opts)
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	now := time.Now()

	suite := []struct {
		name   string
		expr   parser.Expr
		labels labels.Labels
		result promql.Vector
		err    string
	}{
		{
			name:   "nolabels",
			expr:   &parser.NumberLiteral{Val: 1},
			labels: labels.Labels{},
			result: promql.Vector{promql.Sample{
				Metric: labels.FromStrings("__name__", "nolabels"),
				Point:  promql.Point{V: 1, T: timestamp.FromTime(now)},
			}},
		},
		{
			name:   "labels",
			expr:   &parser.NumberLiteral{Val: 1},
			labels: labels.FromStrings("foo", "bar"),
			result: promql.Vector{promql.Sample{
				Metric: labels.FromStrings("__name__", "labels", "foo", "bar"),
				Point:  promql.Point{V: 1, T: timestamp.FromTime(now)},
			}},
		},
	}

	for _, test := range suite {
		rule := NewRecordingRule(test.name, test.expr, test.labels)
		start_time := time.Now()
		result, err := rule.Eval(ctx, now, EngineQueryFunc(engine, storage), nil, 0)
		elapsed := time.Since(start_time)
		t.Log("time used is:", elapsed)
		if test.err == "" {
			require.NoError(t, err)
		} else {
			require.Equal(t, test.err, err.Error())
		}
		require.Equal(t, test.result, result)
	}
}
*/

func newTestFromFile(t testutil.T, filename string) (*promql.Test, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return promql.NewTest(t, string(content))
}

// zeying added function
func TestRuleEvalPressure(t *testing.T) {
	t.Log("Hi! in TestRuleEvalPressure!")
	files, err := filepath.Glob("testdata/*.test")
	require.NoError(t, err)	
	fn := files[0]
	suite, err := newTestFromFile(t, fn)
	require.NoError(t, err)
	defer suite.Close()

	err = suite.Run() // just load data 
	require.NoError(t, err)

	expr, err := parser.ParseExpr(`sum by(instance) (avg_over_time(http_requests[10s]))`)
	// `sum by(instance) (quantile_over_time(0.99, http_requests[10s]))`
	require.NoError(t, err)



	tests := []struct {
		name   string
		expr   parser.Expr
		labels labels.Labels
	}{
		{
			name:   "quantile_over_time:recording_test",
			expr:   &parser.NumberLiteral{Val: 1},
			labels: labels.Labels{},
		},
	}

	for _, test := range tests {
		test.expr = expr
		rule := NewRecordingRule(test.name, test.expr, test.labels)
		for i := 1; i < 10001; i++ { // recording rule evaluated per second
			evalTime := time.Unix((int64)(i), 0)
			_, _ = rule.Eval(suite.Context(), evalTime, EngineQueryFunc(suite.QueryEngine(), suite.Storage()), nil, 0) // no limit here
			/*
			start_time := time.Now()
			result, _ := rule.Eval(suite.Context(), evalTime, EngineQueryFunc(suite.QueryEngine(), suite.Storage()), nil, 0) // no limit here
			elapsed := time.Since(start_time)
			t.Log("Eval time used is:", elapsed)
			t.Log("Eval result is:\n", result)
			*/
		}
		
	}
}



/*

func TestRecordingRuleHTMLSnippet(t *testing.T) {
	expr, err := parser.ParseExpr(`foo{html="<b>BOLD<b>"}`)
	require.NoError(t, err)
	rule := NewRecordingRule("testrule", expr, labels.FromStrings("html", "<b>BOLD</b>"))

	const want = template.HTML(`record: <a href="/test/prefix/graph?g0.expr=testrule&g0.tab=1">testrule</a>
expr: <a href="/test/prefix/graph?g0.expr=foo%7Bhtml%3D%22%3Cb%3EBOLD%3Cb%3E%22%7D&g0.tab=1">foo{html=&#34;&lt;b&gt;BOLD&lt;b&gt;&#34;}</a>
labels:
  html: '&lt;b&gt;BOLD&lt;/b&gt;'
`)

	got := rule.HTMLSnippet("/test/prefix")
	require.Equal(t, want, got, "incorrect HTML snippet; want:\n\n%s\n\ngot:\n\n%s", want, got)
}

// TestRuleEvalDuplicate tests for duplicate labels in recorded metrics, see #5529.
func TestRuleEvalDuplicate(t *testing.T) {
	storage := teststorage.New(t)
	defer storage.Close()

	opts := promql.EngineOpts{
		Logger:     nil,
		Reg:        nil,
		MaxSamples: 10,
		Timeout:    10 * time.Second,
	}

	engine := promql.NewEngine(opts)
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	now := time.Now()

	expr, _ := parser.ParseExpr(`vector(0) or label_replace(vector(0),"test","x","","")`)
	rule := NewRecordingRule("foo", expr, labels.FromStrings("test", "test"))
	_, err := rule.Eval(ctx, now, EngineQueryFunc(engine, storage), nil, 0)
	require.Error(t, err)
	require.EqualError(t, err, "vector contains metrics with the same labelset after applying rule labels")
}

func TestRecordingRuleLimit(t *testing.T) {
	suite, err := promql.NewTest(t, `
		load 1m
			metric{label="1"} 1
			metric{label="2"} 1
	`)
	require.NoError(t, err)
	defer suite.Close()

	require.NoError(t, suite.Run())

	tests := []struct {
		limit int
		err   string
	}{
		{
			limit: 0,
		},
		{
			limit: -1,
		},
		{
			limit: 2,
		},
		{
			limit: 1,
			err:   "exceeded limit of 1 with 2 series",
		},
	}

	expr, _ := parser.ParseExpr(`metric > 0`)
	rule := NewRecordingRule(
		"foo",
		expr,
		labels.FromStrings("test", "test"),
	)

	evalTime := time.Unix(0, 0)

	for _, test := range tests {
		_, err := rule.Eval(suite.Context(), evalTime, EngineQueryFunc(suite.QueryEngine(), suite.Storage()), nil, test.limit)
		if err != nil {
			require.EqualError(t, err, test.err)
		} else if test.err != "" {
			t.Errorf("Expected error %s, got none", test.err)
		}
	}
}

*/