# Text-Collab: Real-time Collaborative Text Editor

A terminal-based collaborative text editor built in Go that allows multiple users to edit the same document simultaneously in real-time.

## Features

- Real-time collaboration
- Terminal-based user interface
- CRDT-based conflict resolution
- Multiple user support with unique colors
- Document synchronization
- Secure WebSocket communication
- Automatic user presence detection
- Debug logging system

## Prerequisites

- Go 1.15 or higher
- Terminal with Unicode support
- Network connectivity for collaboration

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/text-collab.git
cd text-collab
```

2. Install dependencies:
```bash
go mod download
```

## Usage

### Starting the Server

```bash
# Basic server start (default port 8080)
go run editor/server/main.go

# Custom port
go run editor/server/main.go -addr :3000
```

### Starting a Client

```bash
# Basic client connection
go run editor/client/main.go

# With custom options
go run editor/client/main.go -server localhost:8080 -login -scroll
```

### Command Line Flags

#### Server Flags
- `-addr`: Server address (default: ":8080")

#### Client Flags
- `-server`: Server address (default: "localhost:8080")
- `-secure`: Enable secure WebSocket (wss://)
- `-login`: Enable manual username entry
- `-file`: Load content from file
- `-debug`: Enable debug logging
- `-scroll`: Enable cursor scrolling (default: true)

## Network Setup

### Local Network Collaboration

1. Find server machine's IP address:
```bash
# Linux/Mac
ifconfig

# Windows
ipconfig
```

2. Start server on host machine:
```bash
go run editor/server/main.go
```

3. Connect from other machines:
```bash
go run editor/client/main.go -server 192.168.x.x:8080
```

## Editor Controls

### Basic Navigation
- Arrow keys: Move cursor
- Home/End: Start/end of line
- Ctrl+B/F: Move left/right
- Ctrl+P/N: Move up/down

### Document Operations
- Ctrl+S: Save document
- Ctrl+L: Load document
- Esc/Ctrl+C: Exit editor

## File Management

- Default save location: `pairpad-content.txt`
- Log files location: `~/.pairpad/`
  - Regular logs: `pairpad.log`
  - Debug logs: `pairpad-debug.log`

## Collaboration Features

- Real-time character updates
- User presence indication
- Join/leave notifications
- Concurrent editing support
- Automatic conflict resolution



## Security Considerations

- Use `-secure` flag for encrypted connections
- Default WebSocket connection is unencrypted
- Files are saved with 644 permissions
- Logs are stored in user's home directory

## Architecture

The project uses:
- CRDT for conflict resolution
- WebSocket for real-time communication
- Goroutines for concurrent operations
- Mutex for thread safety
- Channel-based message passing

## Flow Chart

![Flow Architecture](flow_chart.png)