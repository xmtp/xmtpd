package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const MARKDOWN_OUTPUT = "doc/metrics_catalog.md"

type Metric struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Description string `yaml:"description,omitempty"`
	File        string `yaml:"file,omitempty"`
}

type MetricType struct {
	enabled bool
	name    string
}

var metricTypes = map[string]MetricType{
	"NewCounterVec":   {true, "Counter"},
	"NewGaugeVec":     {true, "Gauge"},
	"NewHistogramVec": {true, "Histogram"},
	"NewSummaryVec":   {true, "Summary"},
	"NewCounter":      {true, "Counter"},
	"NewGauge":        {true, "Gauge"},
	"NewHistogram":    {true, "Histogram"},
	"NewSummary":      {true, "Summary"},
}

func main() {
	var metrics []Metric
	root := "pkg/metrics"

	fmt.Printf("Parsing %s to generate %s\n", root, MARKDOWN_OUTPUT)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error walking path %s: %v", path, err)
			return err
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		fileMetrics, err := parseFile(path)
		if err != nil {
			log.Printf("Error parsing %s: %v", path, err)
			return err
		}
		metrics = append(metrics, fileMetrics...)
		return nil
	})
	if err != nil {
		log.Fatalf("Error walking through files: %v", err)
	}

	dumpToMarkdown(metrics)
}

func dumpToMarkdown(metrics []Metric) {
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].Name < metrics[j].Name
	})

	var sb strings.Builder
	sb.WriteString("| Name | Type | Description | File |\n")
	sb.WriteString("|------|------|-------------|------|\n")

	for _, m := range metrics {
		desc := m.Description
		if desc == "" {
			desc = "-"
		}
		sb.WriteString(fmt.Sprintf("| `%s` | `%s` | %s | `%s` |\n",
			m.Name, m.Type, desc, m.File))
	}

	if err := os.WriteFile(MARKDOWN_OUTPUT, []byte(sb.String()), 0o644); err != nil {
		log.Fatalf("Error writing Markdown file: %v", err)
	}

	fmt.Printf("✅ %s generated with %d metrics\n", MARKDOWN_OUTPUT, len(metrics))
}

func parseFile(path string) ([]Metric, error) {
	var results []Metric
	fs := token.NewFileSet()
	node, err := parser.ParseFile(fs, path, nil, parser.AllErrors)
	if err != nil {
		return nil, err
	}

	ast.Inspect(node, func(n ast.Node) bool {
		callExpr, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok || !metricTypes[sel.Sel.Name].enabled {
			return true
		}

		// Parse metric name and help string
		if len(callExpr.Args) > 0 {
			firstArg, ok := callExpr.Args[0].(*ast.CompositeLit)
			if ok {
				metric := Metric{Type: metricTypes[sel.Sel.Name].name, File: path}

				for _, elt := range firstArg.Elts {
					if kv, ok := elt.(*ast.KeyValueExpr); ok {
						key := fmt.Sprint(kv.Key)
						switch key {
						case "Name":
							if val, ok := kv.Value.(*ast.BasicLit); ok {
								metric.Name = strings.Trim(val.Value, `"`)
							}
						case "Help":
							if val, ok := kv.Value.(*ast.BasicLit); ok {
								metric.Description = strings.Trim(val.Value, `"`)
							}
						}
					}
				}

				if metric.Name != "" {
					results = append(results, metric)
				}
			}
		}
		return true
	})
	return results, nil
}
