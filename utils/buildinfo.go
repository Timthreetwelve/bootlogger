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
	"strings"

	"golang.org/x/sys/windows/registry"
)

type BuildInfo struct {
	BuildNumber string
	UbrNumber   string
	ProductName string
}

// GetWindowsInfo retrieves the Windows build information from the registry.
// It returns a BuildInfo struct containing the build number, UBR number, and product name.
func GetWindowsInfo() (BuildInfo, error) {
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		log.Println("bootlogger error opening registry key: ", err)
		return BuildInfo{}, fmt.Errorf("error opening registry key: %v", err)
	}
	defer key.Close()

	buildNumber, _, err := key.GetStringValue("CurrentBuildNumber")
	if err != nil {
		log.Println("bootlogger error getting CurrentBuildNumber: ", err)
		return BuildInfo{}, fmt.Errorf("error getting CurrentBuildNumber: %v", err)
	}

	ubrNumber, _, err := key.GetIntegerValue("UBR")
	if err != nil {
		log.Println("bootlogger error getting UBR: ", err)
		return BuildInfo{}, fmt.Errorf("error getting UBR: %v", err)
	}

	productName, _, err := key.GetStringValue("ProductName")
	if err != nil {
		log.Println("bootlogger error getting ProductName: ", err)
		return BuildInfo{}, fmt.Errorf("error getting ProductName: %v", err)
	}
	if buildNumber >= "22000" {
		productName = strings.ReplaceAll(productName, "Windows 10", "Windows 11")
	}

	return BuildInfo{
		BuildNumber: buildNumber,
		UbrNumber:   fmt.Sprintf("%d", ubrNumber),
		ProductName: productName,
	}, nil

}
