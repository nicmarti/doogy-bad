# Datadog Metrics Query Tool

This Go application queries Datadog metrics API to retrieve Django request hit counts filtered by resource names.

## Purpose

The tool fetches time-series data from Datadog for Django application metrics, specifically tracking HTTP request hits with status code 2xx (successful requests) filtered by resource name patterns.

## Prerequisites

### Required Go Packages

Install the Datadog API client for Go:

```bash
go get github.com/DataDog/datadog-api-client-go/v2
```

### Environment Variables

The following environment variables must be set for Datadog API authentication:

- `DD_API_KEY`: Your Datadog API key
- `DD_APP_KEY`: Your Datadog application key  
- `DD_SITE`: Your Datadog site (e.g., `datadoghq.com`)

## Parameters

### Required Parameters

- **RESOURCE_FILTER** (mandatory): Filter pattern for resource names
  - Can be provided as command-line flag: `-RESOURCE_FILTER='pattern'`
  - Or as environment variable: `RESOURCE_FILTER='pattern'`
  - Example: `*admin*` to filter all admin-related endpoints

### Optional Parameters

- **START_DATE**: Start date timestamp in milliseconds
  - Default: One year ago at midnight (calculated dynamically)
  - Can be provided as command-line flag: `-START_DATE=timestamp`
  - Or as environment variable: `START_DATE=timestamp`

- **END_DATE**: End date timestamp in milliseconds
  - Default: Today at midnight (calculated dynamically)
  - Can be provided as command-line flag: `-END_DATE=timestamp`
  - Or as environment variable: `END_DATE=timestamp`

- **SERVICE**: Service name prefix to filter
  - Default: `badoom`
  - Can be provided as command-line flag: `-SERVICE=service_name`
  - Or as environment variable: `SERVICE=service_name`
  - Note: Automatically appends `*` to match any service starting with this prefix

- **--score**: Enable obsolescence scoring analysis
  - Command-line flag only: `--score`
  - When enabled, calculates a score (0-100) for each metric indicating how obsolete it is
  - Higher scores indicate endpoints that are likely unused (old + low usage)
  - Useful for identifying endpoints that can be safely deprecated

**Note**: If neither START_DATE nor END_DATE is specified, the tool uses a default range of the last 12 months (both calculated dynamically from today's date).

## Usage Examples

### Basic Usage with Required Parameter

```bash
# Using environment variables
DD_API_KEY=your_api_key \
DD_APP_KEY=your_app_key \
DD_SITE=datadoghq.com \
RESOURCE_FILTER='/admin/auth/user' \
go run main.go
```

### Using Command-Line Flags

```bash
# Set API credentials as environment variables
export DD_API_KEY=your_api_key
export DD_APP_KEY=your_app_key
export DD_SITE=datadoghq.com

# Run with command-line flag
go run main.go -RESOURCE_FILTER='/admin/auth/user'
```

### Custom Date Range

```bash
DD_API_KEY=your_api_key \
DD_APP_KEY=your_app_key \
DD_SITE=datadoghq.com \
RESOURCE_FILTER='/admin/auth/user' \
START_DATE=1720000000000 \
END_DATE=1730000000000 \
go run main.go
```

### Mixed Parameters (Environment Variables and Flags)

```bash
# API credentials in environment
export DD_API_KEY=your_api_key
export DD_APP_KEY=your_app_key
export DD_SITE=datadoghq.com

# Resource filter as environment variable, dates as flags
RESOURCE_FILTER='/api/v1/products' go run main.go -START_DATE=1720000000000 -END_DATE=1730000000000
```

### Custom Service Name

```bash
# Query a different service (e.g., "myapp" instead of "badoom")
DD_API_KEY=your_api_key \
DD_APP_KEY=your_app_key \
DD_SITE=datadoghq.com \
RESOURCE_FILTER='/api/v1/users' \
SERVICE='myapp' \
go run main.go
```

### Obsolescence Scoring

```bash
# Basic scoring analysis
DD_API_KEY=your_api_key \
DD_APP_KEY=your_app_key \
DD_SITE=datadoghq.com \
RESOURCE_FILTER='*admin*' \
go run main.go --score

# Scoring with built binary
DD_API_KEY=your_api_key \
DD_APP_KEY=your_app_key \
DD_SITE=datadoghq.com \
RESOURCE_FILTER='*admin/address_rules/merchantaddressrule*' \
./doggy_bad --score
```

## Query Details

The application uses a binary search algorithm to find the most recent activity:
- **Metric**: `trace.django.request.hits`
- **Service**: Configurable via SERVICE parameter (default: `badoom*`)
- **HTTP Status Codes**: `2*` (all 2xx success codes)
- **Resource Name**: Converted from URL path to Datadog metric format
- **Aggregation**: Maximum values grouped by resource_name
- **Interval**: Daily (84600000 ms)
- **Search Method**: Binary search to efficiently find the most recent non-null value

## Output

The tool outputs:
1. The parameters being used (dates and filter)
2. JSON-formatted response from Datadog containing:
   - Time series data for matching resources
   - Resource names in group tags
   - Hit counts and metrics metadata

### Standard Output Format

```json
{
  "result": [
    {
      "metric": "resource_name:get_admin/auth/user/",
      "last_seen_timestamp": 1739361600000,
      "last_seen_date": "2025-02-12 12:00:00 UTC",
      "value": 5
    }
  ]
}
```

### Scoring Output Format (with --score flag)

```json
{
  "result": [
    {
      "metric": "resource_name:get_admin/auth/user/",
      "last_seen_timestamp": 1739361600000,
      "last_seen_date": "2025-02-12 12:00:00 UTC", 
      "value": 1,
      "score": 100
    }
  ]
}
```

### Scoring Algorithm

When `--score` flag is used, each metric receives an obsolescence score (0-100):

- **Age Factor** (70% weight): Based on time since last activity
  - 0 points = recent activity
  - 100 points = 5+ months old
  
- **Value Factor** (30% weight): Based on usage intensity
  - 100 points = value ≤ 1 (minimal usage)
  - Decreases logarithmically for higher values
  
- **Final Score**: `(age_score × 0.7) + (value_score × 0.3)`

**Score Ranges:**
- **90-100**: Highly obsolete (strong removal candidate)
- **70-89**: Likely obsolete (investigate usage patterns)
- **40-69**: Moderate risk (monitor trends)
- **0-39**: Active/recent (focus optimization here)

## Error Handling

- If RESOURCE_FILTER is not provided, the program exits with an error message
- Invalid date formats will use default values and print a warning
- API errors are printed to stderr

## Building and Running

### Build the Application

```bash
go build -o doggy_bad main.go
```

### Run the Built Binary

```bash
DD_API_KEY=your_api_key \
DD_APP_KEY=your_app_key \
DD_SITE=datadoghq.com \
RESOURCE_FILTER='/admin/auth/user' \
./doggy_bad
```

## Common Resource Filter Patterns

- `*admin*` - All admin endpoints
- `*api*` - All API endpoints
- `*customer*` - Customer-related endpoints
- `get_*` - All GET request endpoints
- `post_*` - All POST request endpoints
- `*order*` - Order-related endpoints

## Notes

- The tool includes a helper function `pathToDatadogMetrics()` that converts URL paths to Datadog metric names, though it's not currently used in the main query
- Time intervals are set to 1 day (84600000 milliseconds)
- The query returns a maximum of 10 results ordered by count in descending order