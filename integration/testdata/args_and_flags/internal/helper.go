package internal

import "fmt"

//go:generate mv /workspace/internal.txt /workspace/internal_moved.txt
//go:generate sleep 5

func helper() {
	if true {
		fmt.Println("hello world!")
	}
}
