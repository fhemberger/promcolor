package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var version = "development"

// ANSI colors
var (
	Black         = "\033[30m"
	Red           = "\033[31m"
	Green         = "\033[32m"
	Yellow        = "\033[33m"
	Blue          = "\033[34m"
	Magenta       = "\033[35m"
	Cyan          = "\033[36m"
	White         = "\033[37m"
	BrightBlack   = "\033[90m"
	BrightRed     = "\033[91m"
	BrightGreen   = "\033[92m"
	BrightYellow  = "\033[93m"
	BrightBlue    = "\033[94m"
	BrightMagenta = "\033[95m"
	BrightCyan    = "\033[96m"
	BrightWhite   = "\033[97m"
	Reset         = "\033[0m"
)

// Prometheus metric tokens
var (
	CommentColor    = White
	LabelKeyColor   = Blue
	LabelValueColor = Green
	MetricColor     = BrightWhite
	ValueColor      = Magenta
)

// PrometheusLabel is a single key/value combination
type PrometheusLabel struct {
	Key, Value string
}

// PrometheusMetric describes a single metric
type PrometheusMetric struct {
	Name   string
	Labels []PrometheusLabel
	Value  string
}

func splitLabels(token string) []PrometheusLabel {
	var splitLabels []PrometheusLabel

	re := regexp.MustCompile(`([^,]+?)="(.*?)"`)
	labels := re.FindAllStringSubmatch(token, -1)

	for i := 0; i < len(labels); i++ {
		label := labels[i]
		splitLabels = append(splitLabels, PrometheusLabel{Key: label[1], Value: label[2]})
	}
	return splitLabels
}

func colorizeLine(line string) string {

	// Check if line is a comment
	re := regexp.MustCompile(`^#`)
	if re.MatchString(line) {
		return CommentColor + line + Reset
	}

	// Check if line is a metric
	re = regexp.MustCompile(`^([\w_]+)(?:\{(.*)\})?\x20(.+)$`)
	if match := re.FindStringSubmatch(line); len(match) > 1 {
		metric := PrometheusMetric{
			Name:   match[1],
			Labels: splitLabels(match[2]),
			Value:  match[3],
		}

		var colorizedLabels []string
		for i := 0; i < len(metric.Labels); i++ {
			label := metric.Labels[i]
			colorizedLabels = append(colorizedLabels, LabelKeyColor+label.Key+`="`+LabelValueColor+label.Value+`"`+Reset)
		}

		return MetricColor + metric.Name + Reset + "{" + strings.Join(colorizedLabels, ",") + "} " + ValueColor + metric.Value + Reset
	}

	// line is something else, do nothing
	return line
}

func main() {
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 {
		fmt.Println("promcolor", version, "- Colorize piped Prometheus metrics.")
		fmt.Println("\nUsage: curl http://127.0.0.1:9100/metrics | promcolor")
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fmt.Println(colorizeLine(scanner.Text()))
	}
}
