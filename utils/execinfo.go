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
	"os"
	"path/filepath"
)

// GetExecutablePath retrieves the path of the currently running executable
func GetExecutablePath() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %v", err)
	}
	path, err := filepath.EvalSymlinks(ex)
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %v", err)
	}
	return path, nil
}

// GetExecutableFolder returns the directory of the currently running executable
func GetExecutableFolder() (string, error) {
	exec, err := GetExecutablePath()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exec), nil
}

// GetAppLogFile creates a file used for application logging
func GetAppLogFile() (*os.File, error) {
	ex, err := GetExecutablePath()
	if err != nil {
		return nil, err
	}
	exPath := filepath.Dir(ex)
	fullPath := filepath.Join(exPath, "bootlogger.app.log")
	appLogFile, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error creating temporary file: %w", err)
	}
	defer appLogFile.Close()
	return appLogFile, nil
}
