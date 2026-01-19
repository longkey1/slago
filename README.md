# slago

Slack Log Collector CLI - A command-line tool for collecting Slack messages

## Installation

### Using Homebrew

```bash
brew install longkey1/tap/slago
```

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
export SLACK_MENTION="U12345678"     # Optional
```

Or create a config file at `~/.slago.yaml`:

```yaml
token: xoxp-...
author: your-username
mention:
  - U12345678
  - @team-name
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
slago list -d 2025-01-15 --mention U111 --mention @team

# Parallel execution
slago list -m 2025-01 --parallel 4
```

Output is saved to `logs/YYYY/MM/DD/slack.json`.

#### version

```bash
slago version
```

### Global Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--config` | `-c` | Config file path | `~/.slago.yaml` |
| `--token` | | Slack API token | `$SLACK_API_TOKEN` |

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
| `--from` | | Start date (YYYY-MM-DD) |
| `--to` | | End date (YYYY-MM-DD) |

**Other Flags:**

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--thread` | | Fetch entire threads | `false` |
| `--author` | | Filter by author | `$SLACK_AUTHOR` |
| `--mention` | | Filter by mention (can be specified multiple times) | `$SLACK_MENTION` |
| `--parallel` | `-p` | Number of parallel workers | `1` |

## Required Permissions

The Slack API token requires the following scopes:

- `search:read` - Search messages
- `channels:history` - Read channel history
- `channels:read` - Read channel information
- `groups:history` - Read private channel history (optional)
- `groups:read` - Read private channel information (optional)

## Output Format

Output is in JSON format, grouped by thread.

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

# Install tools
make tools

# Release (dry run)
make release type=patch

# Release (actual)
make release type=patch dryrun=false
```

## License

MIT
