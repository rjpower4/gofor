package main

import (
	"fmt"
	"github.com/rjpower4/gofor/pkg/gofor"
)

func main() {
	s, err := gofor.New("samples/gofor.toml")
	if err != nil {
		fmt.Println(err)
	}

	s.Fetch([]string{"naif0012", "de430", "de440"}...)

	_ = s
}
