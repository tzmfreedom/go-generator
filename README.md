# Go Code Generator From PHP

## Install

```bash
$ go install github.com/tzmfreedom/go-generator
```

## Usage

Input PHP Code to stdin.
```bash
$ echo '<?php echo "hello world";' | go-generator
```

Output result from stdout.
```go
package hoge

import "fmt"

func main() {
fmt.Println("hello world")
}
```

