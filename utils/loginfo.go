/*
Copyright © 2025 Tim Kennedy

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

// formatLogLine formats the log line with the computer name, date, time, OS caption, and build number.
// It uses the preferred timestamp format (12-hour or 24-hour) based on the configuration.
func formatLogLine() string {
	var logLine string

	computerName := GetFormattedComputerName()
	dateTimeStr := GetBootTimeFormatted()
	osInfo, _ := GetWindowsInfo()

	if viper.GetBool("no-text") {
		logLine = fmt.Sprintf("%s %s", computerName, dateTimeStr)
	} else {

		logLine = fmt.Sprintf("%s was rebooted on %s", computerName, dateTimeStr)
	}

	if !viper.GetBool("no-buildinfo") {
		logLine += fmt.Sprintf(" [%s build %s.%s]", osInfo.ProductName, osInfo.BuildNumber, osInfo.UbrNumber)
	}
	return logLine
}

// WriteLog writes the log line to the specified log file.
// It opens the file for appending, writes the log line, and flushes the file.
func WriteLog() error {
	// Get the log file path from the configuration
	logFile := viper.GetString("logFile")

	// Get the formatted log line
	logLine := formatLogLine()

	// If dry-run mode is enabled, print the log line instead of writing to the file
	if viper.GetBool("dryrun") {
		log.Println("bootlogger (dry run) log record:", logLine)
		fmt.Println("Dry run: the following record would be written to", logFile)
		fmt.Println(logLine)
		return nil
	}

	// Open the log file for appending
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("bootlogger error: %v", err)
		return fmt.Errorf("error opening log file: %v", err)
	}
	defer file.Close()

	// Write the log line to the file
	if _, err := file.WriteString(logLine + "\r\n"); err != nil {
		return fmt.Errorf("error writing to log file: %v", err)
	}
	log.Println("bootlogger writing log record:", logLine)

	// Print the log line to the console
	if !viper.GetBool("quiet") {
		fmt.Println(logLine)
	}

	// Flush the file to ensure all data is written
	if err := file.Sync(); err != nil {
		log.Printf("bootlogger error flushing log file: %v", err)
		return fmt.Errorf("error flushing log file: %v", err)
	}
	return nil
}
