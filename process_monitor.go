package main

import (
    "compress/gzip"
    "flag"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/shirou/gopsutil/process"
)

const version = "2025-04-07a"

func main() {
    // Define command-line flags
    helpFlag := flag.Bool("h", false, "Display help")
    versionFlag := flag.Bool("v", false, "Display version")
    flag.BoolVar(helpFlag, "help", false, "Display help")
    flag.BoolVar(versionFlag, "version", false, "Display version")

    // Parse command-line flags
    flag.Parse()

    // Handle help and version flags
    if *helpFlag {
        fmt.Println("process monitoring program")
        fmt.Println("-v version")
        return
    }

    if *versionFlag {
        fmt.Println("process monitoring program")
        fmt.Println("Daniel Staniek, TietoEvry - version", version)
        return
    }

    // Create the log directory if it doesn't exist
    logDir := "/var/log/processes"
    if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
        fmt.Printf("Error creating log directory: %v\n", err)
        return
    }

    // Compress old logs
    compressOldLogs(logDir)

    // Get the current date for the log file name
    dateStr := time.Now().Format("2006-01-02_15-04-05")
    logFile := filepath.Join(logDir, fmt.Sprintf("processes-%s.log", dateStr))

    // Initial log of all processes
    processes := getProcesses()
    if err := writeToFile(logFile, processes); err != nil {
        fmt.Printf("Error writing to log file: %v\n", err)
        return
    }

    // Infinite loop to monitor processes every second
    loggedProcesses := make(map[string]struct{})
    for _, proc := range processes {
        loggedProcesses[proc] = struct{}{}
    }

    for {
        currentProcesses := getProcesses()
        newProcesses := []string{}

        for _, proc := range currentProcesses {
            if _, exists := loggedProcesses[proc]; !exists {
                newProcesses = append(newProcesses, proc)
                loggedProcesses[proc] = struct{}{}
            }
        }

        if len(newProcesses) > 0 {
            if err := writeToFile(logFile, newProcesses); err != nil {
                fmt.Printf("Error writing to log file: %v\n", err)
                return
            }
        }

        time.Sleep(1 * time.Second)
    }
}

func compressOldLogs(logDir string) {
    files, err := os.ReadDir(logDir)
    if err != nil {
        fmt.Printf("Error reading log directory: %v\n", err)
        return
    }

    for _, file := range files {
        if !file.IsDir() && strings.HasSuffix(file.Name(), ".log") {
            logFilePath := filepath.Join(logDir, file.Name())
            gzipFilePath := logFilePath + ".gz"

            if _, err := os.Stat(gzipFilePath); os.IsNotExist(err) {
                err := compressFile(logFilePath, gzipFilePath)
                if err != nil {
                    fmt.Printf("Error compressing file %s: %v\n", logFilePath, err)
                } else {
                    os.Remove(logFilePath) // Remove the original log file after compression
                }
            }
        }
    }
}

func compressFile(src, dst string) error {
    srcFile, err := os.Open(src)
    if err != nil {
        return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer dstFile.Close()

    gzipWriter := gzip.NewWriter(dstFile)
    defer gzipWriter.Close()

    _, err = io.Copy(gzipWriter, srcFile)
    return err
}

func getProcesses() []string {
    var processes []string
    procs, err := process.Processes()
    if err != nil {
        fmt.Printf("Error retrieving processes: %v\n", err)
        return processes
    }

    for _, proc := range procs {
        username, _ := proc.Username()
        pid := proc.Pid
        ppid, _ := proc.Ppid()
        createTime, _ := proc.CreateTime()
        cmdline, _ := proc.Cmdline()

        createTimeStr := time.Unix(createTime/1000, 0).Format("2006-01-02 15:04:05")
        if cmdline == "" {
            cmdline = getKernelThreadName(pid)
        }

        processes = append(processes, fmt.Sprintf("%s %d %d %s %s", username, pid, ppid, createTimeStr, cmdline))
    }

    return processes
}

func getKernelThreadName(pid int32) string {
    statusFile := fmt.Sprintf("/proc/%d/status", pid)
    data, err := os.ReadFile(statusFile)
    if err != nil {
        return "[unknown kernel thread]"
    }

    lines := strings.Split(string(data), "\n")
    for _, line := range lines {
        if strings.HasPrefix(line, "Name:") {
            return strings.TrimSpace(strings.Split(line, ":")[1])
        }
    }

    return "[unknown kernel thread]"
}

func writeToFile(filename string, lines []string) error {
    f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return err
    }
    defer f.Close()

    for _, line := range lines {
        if _, err := f.WriteString(line + "\n"); err != nil {
            return err
        }
    }

    return nil
}
