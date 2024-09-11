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
}

```
Output:
```shell
[INFO] Hello world 
[INFO] 2024-09-11 22:36:32.729342 +0800 CST m=+0.000234335 main.go:8 More details 
```