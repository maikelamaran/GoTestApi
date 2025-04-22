// // package and import statements omitted
// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"os"
// 	"strconv"
// 	"strings"
// )

// func main() {
// 	fmt.Print("Enter a grade: ")
// 	reader := bufio.NewReader(os.Stdin) //la entradastandar que es el keyboard
// 	input, err := reader.ReadString('\n')
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	input = strings.TrimSpace(input)

// 	//quito los espacios en blanco
// 	grade, err := strconv.ParseFloat(input, 64) //input lo que quiero convertir, 64 de float64
// 	var status string
// 	if grade >= 60 {
// 		status = "passing"
// 	} else {
// 		status = "failing"
// 	}
// 	fmt.Println("Your grade is", grade, "and you are", status)
// }
