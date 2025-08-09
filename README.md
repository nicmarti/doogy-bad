# Doggy Bad - Datadog Metrics Activity Analyzer

[![Go](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A powerful Go command-line tool that analyzes Django application metrics in Datadog to identify endpoint usage patterns and detect potentially unused endpoints.

Author : Nicolas MARTIGNOLE (Back Market)

# Generated with Claude code

See the prompts that were used to generate this project [claude_prompt.md](./claude_prompt.md)
See also the general [CLAUDE.md](./CLAUDE.md) that helps Claude Code to understand what does this project do.

## 🎯 What it does

**Doggy Bad** helps you answer the critical question: *"When was this endpoint last accessed?"* using Datadog Trace metrics.

The tool uses an intelligent **binary search algorithm** to efficiently find the most recent activity for any Django endpoint by querying Datadog metrics. This is particularly useful for:

- **API cleanup**: Identify unused endpoints that can be safely deprecated
- **Performance monitoring**: Track endpoint usage patterns over time  
- **Security auditing**: Monitor access to sensitive admin endpoints
- **Resource optimization**: Focus optimization efforts on actively used endpoints

## 🔍 How it works

1. **Convert URL paths** to Datadog metric names (e.g., `/admin/auth/user` → `admin/auth/user`)
2. **Binary search** through your specified date range to find the most recent activity
3. **Report results**:
   - ✅ **SUCCESS**: No activity found (potentially unused endpoint)
   - ⚠️ **WARNING**: Found recent activity with timestamp and hit count

## 🚀 Quick Start

### Prerequisites

- Go 1.19 or later
- Datadog API access with valid credentials
- Django application metrics in Datadog

### Installation

```bash
git clone <repository-url>
cd doggy_bad
go build -o doggy_bad main.go
```

### Basic Usage

```bash
# Check admin endpoint activity in the last 12 months
DD_API_KEY=your_api_key \
DD_APP_KEY=your_app_key \
DD_SITE=datadoghq.com \
RESOURCE_FILTER='/admin/auth/user' \
./doggy_bad
```

### Example Output

```
Using START_DATE: 1719525600000 (2024-06-28 00:00:00)
Using END_DATE: 1751061600000 (2025-06-28 00:00:00)
Using RESOURCE_FILTER path: *admin/address_rules/merchantaddressrule*
Converted to Datadog metric: *admin/address_rules/merchantaddressrule*
Using SERVICE: badoom*

Querying for metric: *admin/address_rules/merchantaddressrule*
Response from `MetricsApi.QueryTimeseriesData` for *admin/address_rules/merchantaddressrule*:
{
  "data": {
    "attributes": {
      "series": [
        {
          "group_tags": [
            "resource_name:get_admin/address_rules/merchantaddressrule/_path:object_id_/change/"
          ],
          "query_index": 0,
          "unit": [
            {
              "family": "cache",
              "id": 39,
              "name": "hit",
              "plural": "hits",
              "scale_factor": 1,
              "short_name": ""
            },
            {}
          ]
        },
        {
          "group_tags": [
            "resource_name:get_admin/address_rules/merchantaddressrule/"
          ],
          "query_index": 0,
          "unit": [
            {
              "family": "cache",
              "id": 39,
              "name": "hit",
              "plural": "hits",
              "scale_factor": 1,
              "short_name": ""
            },
            {}
          ]
        }
      ],
      "times": [
        1739361600000
      ],
      "values": [
        [
          1
        ],
        [
          1
        ]
      ]
    },
    "id": "0",
    "type": "timeseries_response"
  }
}

Last seen values:
{
  "result": [
    {
      "metric": "resource_name:get_admin/address_rules/merchantaddressrule/_path:object_id_/change/",
      "last_seen_timestamp": 1739361600000,
      "last_seen_date": "2025-02-12 12:00:00 UTC",
      "value": 1
    },
    {
      "metric": "resource_name:get_admin/address_rules/merchantaddressrule/",
      "last_seen_timestamp": 1739361600000,
      "last_seen_date": "2025-02-12 12:00:00 UTC",
      "value": 1
    }
  ]
}
```

## 📋 Parameters

### Required
- **RESOURCE_FILTER**: URL path to analyze (e.g., `/admin/auth/user`)

### Optional
- **START_DATE**: Start timestamp in milliseconds (default: one year ago at midnight)
- **END_DATE**: End timestamp in milliseconds (default: today at midnight)
- **SERVICE**: Service name prefix (default: `badoom`)
- **--score**: Enable obsolescence scoring to identify unused endpoints

### Environment Variables
- **DD_API_KEY**: Your Datadog API key
- **DD_APP_KEY**: Your Datadog application key  
- **DD_SITE**: Your Datadog site (e.g., `datadoghq.com`)

## 💡 Use Cases

### Find Unused Admin Endpoints
```bash
RESOURCE_FILTER='/admin/auth/user' ./doggy_bad
```

### Check API Endpoint Activity
```bash
RESOURCE_FILTER='/api/v1/products' SERVICE='myapp' ./doggy_bad
```

### Custom Date Range Analysis
```bash
RESOURCE_FILTER='/customer/orders' \
START_DATE=1704067200000 \
END_DATE=1735689600000 \
./doggy_bad
```

### Obsolescence Scoring Analysis
```bash
# Get obsolescence scores for admin endpoints
RESOURCE_FILTER='*admin*' ./doggy_bad --score

# Identify highly obsolete endpoints with detailed scoring
DD_API_KEY=your_api_key \
DD_APP_KEY=your_app_key \
DD_SITE=datadoghq.com \
RESOURCE_FILTER='*admin/address_rules*' \
./doggy_bad --score
```

## 📊 Obsolescence Scoring

The `--score` flag enables advanced analysis to identify potentially obsolete endpoints using a sophisticated scoring algorithm.

### How Scoring Works

The obsolescence score (0-100) combines two factors:

1. **Age Factor** (70% weight): How long since last activity
   - Score increases linearly with age
   - 100 points = 5+ months old
   - 0 points = recent activity

2. **Value Factor** (30% weight): Usage intensity  
   - Higher scores for lower usage values
   - Uses logarithmic scale: value=1 → 100 points, value=10 → 80 points, value=100 → 60 points

**Formula**: `final_score = (age_score × 0.7) + (value_score × 0.3)`

### Score Interpretation

| Score Range | Status | Description |
|-------------|---------|-------------|
| **90-100** | 🔴 **Highly Obsolete** | 5+ months old with minimal usage (≤5 hits) |
| **70-89** | 🟡 **Likely Obsolete** | Either very old OR very low usage |
| **40-69** | 🟠 **Moderate Risk** | Moderately old with some usage |
| **0-39** | 🟢 **Active/Recent** | Recent activity or high usage |

### Example Output with Scoring

```bash
RESOURCE_FILTER='*admin*' ./doggy_bad --score
```

```json
{
  "result": [
    {
      "metric": "resource_name:get_admin/address_rules/merchantaddressrule/",
      "last_seen_timestamp": 1739361600000,
      "last_seen_date": "2025-02-12 12:00:00 UTC",
      "value": 1,
      "score": 100
    }
  ]
}
```

This endpoint scores **100** (highly obsolete) because:
- ⏰ **Age**: Last accessed months ago (high age score)
- 📊 **Usage**: Only 1 hit recorded (high value score)
- 🎯 **Recommendation**: Strong candidate for removal

### Best Practices

- **Score ≥ 90**: Review for removal - likely unused legacy endpoints
- **Score 70-89**: Investigate usage patterns - may be seasonal or edge-case endpoints  
- **Score 40-69**: Monitor for trends - consider consolidation opportunities
- **Score < 40**: Active endpoints - focus optimization efforts here

## 🔧 Advanced Configuration

The tool supports flexible configuration through command-line flags or environment variables:

```bash
# Using command-line flags
./doggy_bad -RESOURCE_FILTER='/admin/users' -SERVICE='backend' -START_DATE=1704067200000

# Using command-line flags with scoring
./doggy_bad -RESOURCE_FILTER='*admin*' --score

# Using environment variables  
export DD_API_KEY=your_key
export RESOURCE_FILTER='/api/health'
export SERVICE='frontend'
./doggy_bad

# Mixed approach with scoring
export DD_API_KEY=your_key
./doggy_bad -RESOURCE_FILTER='*deprecated*' --score
```

## 📖 Documentation

For detailed documentation, build instructions, and advanced usage examples, see:

📚 **[CLAUDE.md](./CLAUDE.md)** - Complete documentation including:
- Detailed installation steps
- All configuration options
- Query details and technical specifications
- Troubleshooting guide
- Common usage patterns

## 🏗️ Technical Details

- **Language**: Go 1.19+
- **Algorithm**: Binary search for efficient date range queries
- **API**: Datadog Metrics API v2
- **Output**: Human-readable console output + structured JSON
- **Time Range**: Configurable with smart defaults (last 12 months)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔒 Security

⚠️ **Important**: Never commit API keys (DD_APP_KEY or DD_API_KEY) to version control. Use environment variables or secure secret management systems.

---

**Built with ❤️ using [Claude Code](https://claude.ai/code)**