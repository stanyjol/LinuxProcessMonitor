# Linux Process Monitor

A Go program to monitor started processes in Linux (primarily tested on RedHat).

## Overview

This tool logs every started process in Linux to log files in `/var/log/processes/*` - the directory is automatically created if it doesn't exist. Log files are rotated daily to manage disk space.

## Requirements

- Go 1.23.0 or later
- Linux operating system
- Root privileges (recommended for full process monitoring)

## Installation

### 1. Initialize a Go Module

In your project directory, run the following command to initialize a new Go module:

```bash
go mod init your_module_name
```

### 2. Install Dependencies

Install the required [gopsutil](https://github.com/shirou/gopsutil) package:

```bash
go get github.com/shirou/gopsutil/process
```

### 3. Build the Project

Compile your Go program:

```bash
go build -o process_monitor
```

## Usage

Run the process monitor:

```bash
./process_monitor
```

### Command Line Options

- `-h`, `--help`: Display help information
- `-v`, `--version`: Display version information

## Log Files

- **Location**: `/var/log/processes/`
- **Format**: Daily rotation with timestamps
- **Content**: Process start events with detailed information