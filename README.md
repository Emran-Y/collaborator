# Text-Collab: Real-time Terminal Collaborative Text Editor

A terminal-based collaborative text editor built in Go that allows multiple users to edit the same document simultaneously in real-time.

## Team Members

| Name                  | Student ID       |
|-----------------------|-----------------|
| Emran Yonas          | UGR/1649/14      |
| Fuad Mohammed Obsu   | UGR/6052/14      |
| Zemenu Mekuriya      | UGR/5017/14      |
| Salman Ali           | UGR/7808/14      |
| Feysel Hussien       | UGR/5898/14      |



## Features

- Real-time collaboration
- Terminal-based user interface
- CRDT-based conflict resolution
- Multiple user support with unique colors
- Document synchronization
- Secure WebSocket communication
- Automatic user presence detection
- Debug logging system



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

## Collaboration Features

- Real-time character updates
- User presence indication
- Join/leave notifications
- Concurrent editing support
- Automatic conflict resolution
## Architecture

The project uses:
- CRDT for conflict resolution
- WebSocket for real-time communication
- Goroutines for concurrent operations
- Mutex for thread safety
- Channel-based message passing

## Flow Chart

![Flow Architecture](/image/flowchart.png)
