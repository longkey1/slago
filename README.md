# slago

Slack Log Collector CLI - A command-line tool for collecting Slack messages

## Installation

### Download from Releases

Download the binary from [Releases](https://github.com/longkey1/slago/releases).

### Build from Source

```bash
git clone https://github.com/longkey1/slago.git
cd slago
make build
```

## Usage

### Environment Variables

```bash
export SLACK_API_TOKEN="xoxp-..."  # Required
export SLACK_AUTHOR="your-username"  # Optional
export SLACK_MENTION="U12345678,@john.doe,@team-name"  # Optional: comma-separated
```

### Commands

#### get

Fetch a single message or thread from a Slack URL.

```bash
# Fetch a single message
slago get "https://xxx.slack.com/archives/C123/p456"

# Fetch the entire thread
slago get "https://xxx.slack.com/archives/C123/p456" --thread
```

#### list

Collect messages for a date range and save to files.

```bash
# Collect messages for a specific day
slago list --day 2025-01-15

# Collect messages for an entire month
slago list --month 2025-01

# Collect messages for a custom date range
slago list --from 2025-01-01 --to 2025-01-15

# Combine options
slago list -m 2025-01 --thread --author U12345678
slago list -d 2025-01-15 --mention U111 --mention @john.doe --mention @team

# Parallel execution
slago list -m 2025-01 --parallel 4
```

Output is saved to `logs/YYYY/MM/DD/slack.json`.

#### merge

Merge multiple JSON files and deduplicate threads/messages.

```bash
# Merge all JSON files in a directory
slago merge ./logs

# With explicit --dir flag
slago merge --dir ./logs

# Filter by file pattern
slago merge ./logs --pattern "slack*.json"
slago merge ./logs -p "2025-*.json"

# Recursive search (include subdirectories)
slago merge ./logs --recursive
slago merge ./logs -r -p "*.json"
```

Output is written to stdout.

#### version

```bash
slago version
```

### Global Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--token` | Slack API token | `$SLACK_API_TOKEN` |

### get Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--thread` | Fetch the entire thread | `false` |

### list Flags

**Date Range (mutually exclusive):**

| Flag | Short | Description |
|------|-------|-------------|
| `--day` | `-d` | Single day (YYYY-MM-DD) |
| `--month` | `-m` | Entire month (YYYY-MM) |
| `--from` | | Start date (YYYY-MM-DD, inclusive) |
| `--to` | | End date (YYYY-MM-DD, inclusive) |

**Other Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--thread` | | Fetch entire threads | `false` |
| `--author` | | Filter by author | `$SLACK_AUTHOR` |
| `--mention` | | Filter by mention (User ID or `@username`/`@group-name`, repeatable) | `$SLACK_MENTION` |
| `--parallel` | `-p` | Number of parallel workers | `1` |

### merge Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--dir` | `-d` | Target directory | |
| `--pattern` | `-p` | File name glob pattern | `*.json` |
| `--recursive` | `-r` | Search subdirectories recursively | `false` |

## Required Permissions

The Slack API token requires the following scopes:

- `search:read` - Search messages
- `channels:history` - Read channel history
- `channels:read` - Read channel information
- `groups:history` - Read private channel history (optional)
- `groups:read` - Read private channel information (optional)

## Output Format

Output is in JSON format, grouped by thread.

**Note:** The output structure is not the raw Slack API response. Messages are transformed into a simplified, consistent format:

| slago field | Source |
|-------------|--------|
| `id` | Message timestamp (`ts`) |
| `content` | Message text |
| `author` | User ID |
| `timestamp` | Parsed to ISO 8601 format |
| `mentions` | Extracted from `<@USER\|name>` patterns in text |
| `attached_links` | Extracted from text and attachments |
| `is_thread_parent` | Calculated from `thread_ts` |

```json
[
  {
    "thread_id": "1716192523.567890",
    "thread_permalink": "https://xxx.slack.com/archives/C123/p456",
    "channel": "general",
    "channel_id": "C12345678",
    "messages": [
      {
        "id": "1716192523.567890",
        "type": "slack_message",
        "content": "Hello, World!",
        "author": "U12345678",
        "timestamp": "2025-01-15T10:30:00Z",
        "channel": "general",
        "channel_id": "C12345678",
        "thread_ts": "1716192523.567890",
        "is_thread_parent": true
      }
    ],
    "message_count": 1
  }
]
```

## Development

```bash
# Build
make build

# Test
make test

# Clean
make clean
```

## License

MIT
