# envloader for GO

[![CI Workflow](https://github.com/mhenselin/envloader/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/mhenselin/envloader/actions/workflows/ci.yml)

This component aims towards easy and lazy loading of struct values via env variables.

as easy as that:

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mhenselin/envloader"
)

// TestType - a sample struct
type TestType struct {
	Name       string `env:"NAME,required"`
	Test1      string `env:"-"`
	Test2      string `env:"TEST"`
	AnonIsHere string
}

func main() {
	// set some env vars for testing
	err := os.Setenv("NAME", "gopher")
	if err != nil {
		log.Fatal(err)
	}

	err = os.Setenv("ANON_IS_HERE", "anonymous")
	if err != nil {
		log.Fatal(err)
	}

	var myTypeOne TestType
	// easy loading does NOT magically "detect" other variables
	err = envloader.LoadEnv(&myTypeOne)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("regular: %#v\n", myTypeOne)

	var myTypeTwo TestType
	// lazy loading instead DOES detect other vars by snake-case-ing 
	// the parameters name
	// so AnonIsHere becomes the search string ANON_IS_HERE and if
	// that env var is present - it is used
	err = envloader.LoadEnvLazy(&myTypeTwo)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("lazy: %#v\n", myTypeTwo)
}
```

output:
```
regular: main.TestType{Name:"gopher", Test1:"", Test2:"", AnonIsHere:""}
lazy: main.TestType{Name:"gopher", Test1:"", Test2:"", AnonIsHere:"anonymous"}
```
