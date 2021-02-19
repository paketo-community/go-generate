package main

import "fmt"

//go:generate mv /workspace/main.txt /workspace/main_moved.txt
//go:generate sleep 5

func main() {
	if true {
		fmt.Println("hello world!")
	}
}
