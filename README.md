Process monitor in GO

Program to monitor started processes in Linux (mainly in RedHat).

It log every started process in Linux to log file in /var/log/processes/* - directory is automatically created.
Log file are "rotated" every day.





Install gopsutil:

1. Initialize a Go Module: In your project directory, run the following command to initialize a new Go module:

go mod init your_module_name

2. Install gopsutil: Now that your project is a module, you can install gopsutil using:

go get github.com/shirou/gopsutil/process


3. Build Your Project: You can now compile your Go program as before:

go build -o process_monitor
