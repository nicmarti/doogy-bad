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
# Check admin endpoint activity in the last 18 months
DD_API_KEY=your_api_key \
DD_APP_KEY=your_app_key \
DD_SITE=datadoghq.com \
RESOURCE_FILTER='/admin/auth/user' \
./doggy_bad
```

### Example Output

```
=== Starting Binary Search for metric: admin/auth/user ===
Date range: 28/12/2023 to 28/06/2025

[Depth 0] Searching range: 2023-12-28 to 2025-06-28
[Depth 1] Searching upper half...
Found value: 157 at 15 March 2024 14:30:25 UTC

=== Binary Search Results ===
⚠️  WARNING: Found recent activity for metric: admin/auth/user
   Last seen: 15 March 2024 14:30:25 UTC
   Value: 157
   Timestamp: 1710509425000
```

## 📋 Parameters

### Required
- **RESOURCE_FILTER**: URL path to analyze (e.g., `/admin/auth/user`)

### Optional
- **START_DATE**: Start timestamp in milliseconds (default: 18 months ago)
- **END_DATE**: End timestamp in milliseconds (default: today)
- **SERVICE**: Service name prefix (default: `badoom`)

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

## 🔧 Advanced Configuration

The tool supports flexible configuration through command-line flags or environment variables:

```bash
# Using command-line flags
./doggy_bad -RESOURCE_FILTER='/admin/users' -SERVICE='backend' -START_DATE=1704067200000

# Using environment variables  
export DD_API_KEY=your_key
export RESOURCE_FILTER='/api/health'
export SERVICE='frontend'
./doggy_bad
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
- **Time Range**: Configurable with smart defaults (18 months)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔒 Security

⚠️ **Important**: Never commit API keys to version control. Use environment variables or secure secret management systems.

---

**Built with ❤️ using [Claude Code](https://claude.ai/code)**