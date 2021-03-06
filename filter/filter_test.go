// Copyright 2015 Eleme Inc. All rights reserved.

package filter

import (
	"github.com/eleme/banshee/models"
	"github.com/eleme/banshee/util/assert"
	"path/filepath"
	"testing"
)

func TestSimple(t *testing.T) {
	// New and add rules.
	filter := New()
	rule1 := &models.Rule{Pattern: "a.*.c.d"}
	rule2 := &models.Rule{Pattern: "a.b.c.*"}
	filter.addRule(rule1)
	filter.addRule(rule2)
	// Test
	rules1 := filter.MatchedRules(&models.Metric{Name: "nothing"})
	assert.Ok(t, 0 == len(rules1))

	rules2 := filter.MatchedRules(&models.Metric{Name: "a.b.c.e"})
	assert.Ok(t, 1 == len(rules2))
	assert.Ok(t, rules2[0] == rule2)

	rules3 := filter.MatchedRules(&models.Metric{Name: "a.e.c.d"})
	assert.Ok(t, 1 == len(rules3))
	assert.Ok(t, rules3[0] == rule1)

	rules4 := filter.MatchedRules(&models.Metric{Name: "a.b.c.d"})
	assert.Ok(t, 2 == len(rules4))
}

func BenchmarkRules1KNativeBest(b *testing.B) {
	var rules []*models.Rule
	for i := 0; i < 1024; i++ {
		rules = append(rules, &models.Rule{Pattern: "a.*.c." + string(i)})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i := 0; i < 1024; i++ {
			filepath.Match(rules[i].Pattern, "x")
		}
	}
}

func BenchmarkRules1kBest(b *testing.B) {
	filter := New()
	for i := 0; i < 1024; i++ {
		filter.addRule(&models.Rule{Pattern: "a.*.c." + string(i)})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.MatchedRules(&models.Metric{Name: "x.b.c." + string(i&1024)})
	}
}

func BenchmarkRules1kWorst(b *testing.B) {
	filter := New()
	for i := 0; i < 1024; i++ {
		filter.addRule(&models.Rule{Pattern: "a.*.c." + string(i)})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.MatchedRules(&models.Metric{Name: "a.b.c." + string(i&1024)})
	}
}

func BenchmarkRules2kWorst(b *testing.B) {
	filter := New()
	for i := 0; i < 1024*2; i++ {
		filter.addRule(&models.Rule{Pattern: "a.*.c." + string(i)})
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filter.MatchedRules(&models.Metric{Name: "a.b.c." + string(i&65535)})
	}
}
