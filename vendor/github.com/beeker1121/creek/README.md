# Creek [![GoDoc](http://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/beeker1121/creek) [![License](http://img.shields.io/badge/license-mit-blue.svg)](https://raw.githubusercontent.com/beeker1121/creek/master/LICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/beeker1121/creek)](https://goreportcard.com/report/github.com/beeker1121/creek) [![Build Status](https://travis-ci.org/beeker1121/creek.svg?branch=master)](https://travis-ci.org/beeker1121/creek)

A simple log rotator for Go on Linux platforms.

## Usage

Creek is meant to be used with the standard log library.

The `New` method accepts a file to log to, as well as the max size, in megabytes, of each log file before rolling.

Sample usage:

```go
package main

import (
	"log"

	"creek"
)

func main() {
	// Create a new logger that stores to a http.log file
	// with a max size of 10 MB before rolling over.
	logger := log.New(creek.New("/var/log/your_app/http.log", 10), "Logged: ", log.Lshortfile|log.LstdFlags)

	// Print to the log.
	logger.Println("Testing the log file")
}
```