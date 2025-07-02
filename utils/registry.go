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

	"golang.org/x/sys/windows/registry"
)

const (
	autostartKey    string       = `Software\Microsoft\Windows\CurrentVersion\Run`
	bootloggerValue string       = "bootlogger"
	hiveKey         registry.Key = registry.CURRENT_USER
)

// GetAutostartStatus retrieves the current autostart status of bootlogger
// by checking the registry key HKCU\Software\Microsoft\Windows\CurrentVersion\Run
// It prints whether bootlogger is set to run at startup and from which path
func GetAutostartStatus() (string, error) {
	key, err := registry.OpenKey(hiveKey, autostartKey, registry.QUERY_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {

			return "bootlogger is not set to run at startup.", nil
		}
		return "", fmt.Errorf("failed to open registry key: %v.", err)
	}
	defer key.Close()

	value, _, err := key.GetStringValue(bootloggerValue)
	if err != nil {
		if err == registry.ErrNotExist {
			return "bootlogger is not set to run at startup.", nil
		}
		return "", fmt.Errorf("failed to read registry value: %v.", err)
	}
	if value != "" {
		return fmt.Sprintf("%s is present in 'HKCU\\%s' and will run from: %s when Windows starts. ", bootloggerValue, autostartKey, value), nil
	}
	return "bootlogger is not set to run at startup.", nil
}

// DisableAutostart removes bootlogger from HKCU\Software\Microsoft\Windows\CurrentVersion\Run
// This prevents bootlogger from running automatically at system startup
// It does not delete the executable, only the registry entry
func DisableAutostart() (string, error) {
	key, err := registry.OpenKey(hiveKey, autostartKey, registry.WRITE)
	if err != nil {
		return "", fmt.Errorf("failed to open registry key: %v", err)
	}
	defer key.Close()

	err = key.DeleteValue(bootloggerValue)
	if err != nil {
		if err == registry.ErrNotExist {
			return "bootlogger is not set to run at startup.", nil
		}
		return "", fmt.Errorf("error deleting registry entry: %v", err)
	}
	return fmt.Sprintf("bootlogger was removed from 'HKCU\\%s' and will no longer start with Windows.", autostartKey), nil
}

// EnableAutostart adds bootlogger to HKCU\Software\Microsoft\Windows\CurrentVersion\Run
// It sets the value to the path of the currently running executable
// This allows bootlogger to run automatically at system startup
func EnableAutostart() (string, error) {
	key, err := registry.OpenKey(hiveKey, autostartKey, registry.CREATE_SUB_KEY|registry.WRITE)
	if err != nil {
		return "", fmt.Errorf("failed to open registry key: %v", err)
	}
	defer key.Close()

	path, err := GetExecutablePath()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %v", err)
	}

	err = key.SetStringValue(bootloggerValue, path)
	if err != nil {
		return "", fmt.Errorf("failed to set registry value: %v", err)
	}
	return fmt.Sprintf("bootlogger was added to 'HKCU\\%s' and will run from: %s when Windows is started.", autostartKey, path), nil
}
