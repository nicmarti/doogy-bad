// Timeseries cross product query returns "OK" response

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/DataDog/datadog-api-client-go/v2/api/datadog"
	"github.com/DataDog/datadog-api-client-go/v2/api/datadogV2"
)

// MetricResult represents the result for a single metric
type MetricResult struct {
	Metric            string      `json:"metric"`
	LastSeenTimestamp int64       `json:"last_seen_timestamp"`
	LastSeenDate      string      `json:"last_seen_date"`
	Value             interface{} `json:"value"`
}

// Result wraps the array of metric results
type Result struct {
	Result []MetricResult `json:"result"`
}

// pathToDatadogMetrics converts a web path to a Datadog metric name
// Dynamic parameters enclosed in < > are transformed to _ prefix.
// Example: /bm/branding/<uuid:id> -> "/bm/branding/_uuid:id"
func pathToDatadogMetrics(path string) string {
	// Replace < and > with _
	metricPath := strings.ReplaceAll(path, "<", "_")
	metricPath = strings.ReplaceAll(metricPath, ">", "")

	return metricPath
}

// parseDate parses a date string and returns the timestamp in milliseconds.
// If parsing fails or the string is empty, it returns the defaultValue.
func parseDate(dateStr string, defaultValue int64) int64 {
	if dateStr == "" {
		return defaultValue
	}
	parsed, err := strconv.ParseInt(dateStr, 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing date: %v\n", err)
		fmt.Fprintf(os.Stderr, "Using default value: %d\n", defaultValue)
		return defaultValue
	}
	return parsed
}

// findLastNonNullValues processes the Datadog API response and returns the last non-null timestamp for each metric
func findLastNonNullValues(resp datadogV2.TimeseriesFormulaQueryResponse) Result {
	var results []MetricResult

	// Check if we have the expected data structure
	if resp.Data == nil || resp.Data.Attributes == nil {
		return Result{Result: results}
	}

	attrs := resp.Data.Attributes
	
	// Get series, times, and values
	series := attrs.Series
	times := attrs.Times
	values := attrs.Values

	// Process each series
	for i, s := range series {
		if i >= len(values) {
			break
		}

		// Get the metric name from group_tags
		var metricName string
		if len(s.GroupTags) > 0 {
			metricName = s.GroupTags[0]
		}

		// Find the last non-null value for this series
		seriesValues := values[i]
		var lastTimestamp int64
		var lastValue interface{}

		// Iterate from the end to find the last non-null value
		for j := len(seriesValues) - 1; j >= 0; j-- {
			if seriesValues[j] != nil {
				if j < len(times) {
					lastTimestamp = times[j]
					lastValue = seriesValues[j]
				}
				break
			}
		}

		// Only add if we found a non-null value
		if lastValue != nil {
			// Convert timestamp to human-readable date in GMT
			lastSeenTime := time.Unix(lastTimestamp/1000, 0).UTC()
			lastSeenDate := lastSeenTime.Format("2006-01-02 15:04:05 MST")
			
			results = append(results, MetricResult{
				Metric:            metricName,
				LastSeenTimestamp: lastTimestamp,
				LastSeenDate:      lastSeenDate,
				Value:             lastValue,
			})
		}
	}

	return Result{Result: results}
}

func main() {
	// Define command line flags
	var startDateStr string
	var endDateStr string
	var resourceFilter string
	var serviceName string
	flag.StringVar(&startDateStr, "START_DATE", "", "Start date timestamp in milliseconds")
	flag.StringVar(&endDateStr, "END_DATE", "", "End date timestamp in milliseconds")
	flag.StringVar(&resourceFilter, "RESOURCE_FILTER", "", "Resource name filter (e.g., *admin*)")
	flag.StringVar(&serviceName, "SERVICE", "", "Service name prefix (default: badoom)")
	flag.Parse()

	// Default dates
	defaultStartDate := int64(1719525600000)
	defaultEndDate := int64(1751061600000) // 2025-03-01 00:00:00

	// Override if START_DATE env var is set
	if envStartDate := os.Getenv("START_DATE"); envStartDate != "" && startDateStr == "" {
		startDateStr = envStartDate
	}

	// Override if END_DATE env var is set
	if envEndDate := os.Getenv("END_DATE"); envEndDate != "" && endDateStr == "" {
		endDateStr = envEndDate
	}

	// Override if RESOURCE_FILTER env var is set
	if envResourceFilter := os.Getenv("RESOURCE_FILTER"); envResourceFilter != "" && resourceFilter == "" {
		resourceFilter = envResourceFilter
	}

	// Override if SERVICE env var is set
	if envService := os.Getenv("SERVICE"); envService != "" && serviceName == "" {
		serviceName = envService
	}

	// Default service name
	if serviceName == "" {
		serviceName = "badoom"
	}

	// RESOURCE_FILTER is mandatory
	if resourceFilter == "" {
		fmt.Fprintf(os.Stderr, "Error: RESOURCE_FILTER is required\n")
		fmt.Fprintf(os.Stderr, "Usage: %s -RESOURCE_FILTER=<path> or set RESOURCE_FILTER environment variable\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Example: %s -RESOURCE_FILTER='/admin/auth/user'\n", os.Args[0])
		os.Exit(1)
	}

	// Parse the dates using the generic function
	startDate := parseDate(startDateStr, defaultStartDate)
	endDate := parseDate(endDateStr, defaultEndDate)

	// Convert the resource filter path to Datadog metric names
	metricName := pathToDatadogMetrics(resourceFilter)

	// Convert timestamps to human-readable format
	startDateTime := time.Unix(startDate/1000, 0)
	endDateTime := time.Unix(endDate/1000, 0)
	fmt.Printf("Using START_DATE: %d (%s)\n", startDate, startDateTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Using END_DATE: %d (%s)\n", endDate, endDateTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Using RESOURCE_FILTER path: %s\n", resourceFilter)
	fmt.Printf("Converted to Datadog metric: %v\n", metricName)
	fmt.Printf("Using SERVICE: %s*\n", serviceName)

	// Initialize Datadog API client
	ctx := datadog.NewDefaultContext(context.Background())
	configuration := datadog.NewConfiguration()
	apiClient := datadog.NewAPIClient(configuration)
	api := datadogV2.NewMetricsApi(apiClient)

	fmt.Printf("\nQuerying for metric: %s\n", metricName)

	body := datadogV2.TimeseriesFormulaQueryRequest{
		Data: datadogV2.TimeseriesFormulaRequest{
			Attributes: datadogV2.TimeseriesFormulaRequestAttributes{
				Formulas: []datadogV2.QueryFormula{
					{
						Formula: "a",
						Limit: &datadogV2.FormulaLimit{
							Count: datadog.PtrInt32(10),
							Order: datadogV2.QUERYSORTORDER_DESC.Ptr(),
						},
					},
				},
				From:     startDate,
				Interval: datadog.PtrInt64(84600000), // 1 day
				Queries: []datadogV2.TimeseriesQuery{
					datadogV2.TimeseriesQuery{
						MetricsTimeseriesQuery: &datadogV2.MetricsTimeseriesQuery{
							DataSource: datadogV2.METRICSDATASOURCE_METRICS,
							Query:      fmt.Sprintf("max:trace.django.request.hits{service:%s*, http.status_code:2*, resource_name:%s} by {resource_name}.as_count()", serviceName, metricName),
							Name:       datadog.PtrString("a"),
						}},
				},
				To: endDate,
			},
			Type: datadogV2.TIMESERIESFORMULAREQUESTTYPE_TIMESERIES_REQUEST,
		},
	}

	resp, r, err := api.QueryTimeseriesData(ctx, body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `MetricsApi.QueryTimeseriesData` for %s: %v\n", metricName, err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}

	responseContent, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Fprintf(os.Stdout, "Response from `MetricsApi.QueryTimeseriesData` for %s:\n%s\n", metricName, responseContent)

	// Process the response to find last non-null values
	result := findLastNonNullValues(resp)
	resultJSON, _ := json.MarshalIndent(result, "", "  ")
	fmt.Printf("\nLast seen values:\n%s\n", resultJSON)
}
