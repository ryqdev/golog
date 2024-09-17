# golog

## Get golog
```shell
go get -u github.com/ryqdev/golog
```

## Usage
```go
package main

import "github.com/ryqdev/golog"

func main() {
	golog.Info("Hello world")
	golog.ShowDetail(true)
	golog.Info("More details")

	// write log to local file
	golog.SetLogFile(true)
	golog.Info("write to log directory")
}