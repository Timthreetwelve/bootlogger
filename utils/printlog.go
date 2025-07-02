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
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

// PrintLog reads the bootlogger log file and prints its contents to the console.
// It checks if the log file exists and is not empty before attempting to read it.
func PrintLog() error {
	// Get the log file path from the configuration
	logFile := viper.GetString("logFile")
	if logFile == "" {
		log.Println("bootlogger: log file path is not configured")
		return fmt.Errorf("log file path is not configured.")
	}

	// Check if the log file exists and is not empty
	stat, err := os.Stat(logFile)
	if os.IsNotExist(err) {
		log.Println("bootlogger: log file exist")
		return fmt.Errorf("log file does not exist: %s", logFile)
	}
	if stat.Size() == 0 {
		log.Println("bootlogger: log file is empty")
		return fmt.Errorf("log file is empty: %s", logFile)
	}

	// Open the log file for reading
	file, err := os.Open(logFile)
	if err != nil {
		log.Println("bootlogger: error opening log file ", err)
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// Using bufio scanner here for memory-efficient line-by-line reading.
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	// Check for errors during scanning
	if err := scanner.Err(); err != nil {
		log.Println("bootlogger: error reading log file ", err)
		return fmt.Errorf("Error reading the log file: %w", err)
	}
	return nil
}
