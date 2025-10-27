# 📘 MangaHub CLI

**MangaHub CLI** is a command-line interface for managing and tracking your manga reading experience.
It provides access to all core MangaHub features, including:

* Manga discovery
* Reading progress tracking
* Real-time synchronization
* Community chat functionality

---

## 🧩 Table of Contents

* [Features](#features)
* [System Requirements](#system-requirements)
* [Supported Platforms](#supported-platforms)
* [Installation](#installation)
* [Quick Start](#quick-start)
* [Configuration](#configuration)
* [Authentication Commands](#authentication-commands)
* [Manga Management](#manga-management)
* [Library Operations](#library-operations)
* [Progress Tracking](#progress-tracking)
* [Network Features](#network-features)
* [Chat System](#chat-system)
* [Server Management](#server-management)
* [Configuration Management](#configuration-management)
* [Troubleshooting](#troubleshooting)
* [Advanced Features](#advanced-features)
* [Examples & Workflows](#examples--workflows)
* [Support & Updates](#support--updates)

---

## Features

* CLI-based manga management and tracking
* Built-in local SQLite database
* Multi-protocol support: HTTP, TCP, UDP, gRPC, WebSocket
* Real-time synchronization across devices
* Interactive chat and notifications
* Scriptable JSON outputs for automation

---

## System Requirements

* Go **1.19+**
* SQLite **3.x**
* Network connectivity (for sync & notifications)
* Terminal with **UTF-8** support

### Supported Platforms

* Linux (x64, ARM)
* macOS (Intel, Apple Silicon)
* Windows (x64)

---

## Installation

```bash
# Download the latest release
wget https://github.com/yourorg/mangahub/releases/latest/mangahub-cli

# Make executable (Linux/macOS)
chmod +x mangahub-cli

# Move to system path
sudo mv mangahub-cli /usr/local/bin/mangahub

# Verify installation
mangahub version
```

---

## Quick Start

```bash
# 1. Initialize configuration
mangahub init

# 2. Start the server
mangahub server start

# 3. Register and log in
mangahub auth register --username myuser --email myuser@example.com
mangahub auth login --username myuser

# 4. Search and add manga
mangahub manga search "one piece"
mangahub library add --manga-id one-piece --status reading

# 5. Update progress
mangahub progress update --manga-id one-piece --chapter 1095
```

---

## Configuration

Default config file: `~/.mangahub/config.yaml`

```yaml
server:
  host: "localhost"
  http_port: 8080
  tcp_port: 9090
  udp_port: 9091
  grpc_port: 9092
  websocket_port: 9093

database:
  path: "~/.mangahub/data.db"

user:
  username: ""
  token: ""

sync:
  auto_sync: true
  conflict_resolution: "last_write_wins"

notifications:
  enabled: true
  sound: false

logging:
  level: "info"
  path: "~/.mangahub/logs/"
```

---

## Authentication Commands

| Command                         | Description                   |
| ------------------------------- | ----------------------------- |
| `mangahub auth register`        | Create a new account          |
| `mangahub auth login`           | Log in with username or email |
| `mangahub auth logout`          | Log out and clear session     |
| `mangahub auth status`          | Show current login info       |
| `mangahub auth change-password` | Change your password          |

Example:

```bash
mangahub auth register --username johndoe --email john@example.com
mangahub auth login --username johndoe
```

---

## Manga Management

### Search Manga

```bash
mangahub manga search "attack on titan"
mangahub manga search "romance" --genre romance --status completed
```

### View Manga Details

```bash
mangahub manga info one-piece
```

### List All Manga

```bash
mangahub manga list --genre shounen --page 2 --limit 20
```

---

## Library Operations

| Command                   | Description               |
| ------------------------- | ------------------------- |
| `mangahub library add`    | Add manga to your library |
| `mangahub library list`   | View your library         |
| `mangahub library remove` | Remove manga from library |
| `mangahub library update` | Update status or rating   |

Example:

```bash
mangahub library add --manga-id one-piece --status reading
mangahub library update --manga-id one-piece --rating 10
```

---

## Progress Tracking

```bash
mangahub progress update --manga-id one-piece --chapter 1095
mangahub progress history --manga-id one-piece
mangahub progress sync
```

---

## Network Features

### TCP Sync

```bash
mangahub sync connect
mangahub sync status
mangahub sync monitor
```

### UDP Notifications

```bash
mangahub notify subscribe
mangahub notify test
```

### gRPC

```bash
mangahub grpc manga search --query "bleach"
```

---

## Chat System

```bash
mangahub chat join
mangahub chat join --manga-id one-piece
mangahub chat send "Hello everyone!"
```

Interactive commands:

```
/help    - Show help
/users   - List online users
/quit    - Leave chat
/pm      - Send private message
```

---

## Server Management

Start and monitor servers:

```bash
mangahub server start
mangahub server status
mangahub server health
mangahub server logs --follow
mangahub server stop
```

Example services:

* HTTP API (8080)
* TCP Sync (9090)
* UDP Notify (9091)
* gRPC (9092)
* WebSocket Chat (9093)

---

## Configuration Management

```bash
mangahub config show
mangahub config set server.host "192.168.1.100"
mangahub config reset
```

Profiles:

```bash
mangahub profile create --name work
mangahub profile switch --name work
mangahub profile list
```

---

## Advanced Features

* **Batch Operations**

  ```bash
  mangahub library batch-add --file manga-list.txt
  mangahub progress batch-update --file progress.csv
  ```

* **Backup & Restore**

  ```bash
  mangahub backup create --output backup.tar.gz
  mangahub backup restore --input backup.tar.gz
  ```

* **Database Maintenance**

  ```bash
  mangahub db check
  mangahub db optimize
  mangahub db repair
  ```

---

## Troubleshooting

### Authentication Issues

```bash
mangahub auth clear
mangahub auth register --username <username> --email <email>
```

### Connectivity

```bash
mangahub server ping
mangahub sync reconnect
```

### Database

```bash
mangahub db repair
mangahub init --force
```

Enable debug mode:

```bash
mangahub --verbose <command>
mangahub config set logging.level trace
```

---

## Examples & Workflows

Daily routine:

```bash
mangahub server start &
mangahub sync connect
mangahub notify subscribe
mangahub manga search --new-chapters
mangahub progress update --manga-id one-piece --chapter 1095
```

Scripting:

```bash
#!/bin/bash
while IFS=',' read -r id chapter; do
  mangahub progress update --manga-id "$id" --chapter "$chapter"
done < progress.csv
```

JSON Output:

```bash
mangahub library list --output json | jq '.results[].title'
```

---

## Support & Updates

```bash
# Show help
mangahub help
mangahub manga help

# Check version & updates
mangahub version
mangahub update check
mangahub update install
```

To report issues:

1. Run with `--verbose`
2. Check logs: `mangahub logs errors`
3. Include system info: `mangahub system info`

---

### License

This project is distributed under the **MIT License**.
See `LICENSE` for more information.