package main

import (
	"time"
	"fmt"
)

var thisYear string = time.RFC3339[0:8]
func PRINTTIME(){
	fmt.Println(thisYear)
}
