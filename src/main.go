/*
 * This application can be used to experiment and test various serial port options
 */

package main

import (
 "fmt"
 "chuankou"
)

func main() {
    num1,num2:=chuankou.Read()
    fmt.Println(num1,num2)
}