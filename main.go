package main

import (
	"github.com/GaruGaru/Tao/cmd"
	"time"
	"fmt"
)

func main() {
	cmd.Execute()

	d := time.Duration(2)*time.Second

	fmt.Println(d)
}