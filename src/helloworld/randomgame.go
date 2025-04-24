// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"os"
// 	"strconv"
// 	"strings"
// 	"time"
// )

// func main() {
// 	seconds := time.Now().Unix()
// 	rand.Seed(seconds)
// 	target := rand.Intn(100) + 1
// 	fmt.Println("he elegido un número random entre 1 y 100.")
// 	fmt.Println("lo adivinas?")
// 	fmt.Println(target)
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print("Make a guess:")
// 	input, err := reader.ReadString('\n')
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	input = strings.TrimSpace(input)  //quitar la nueva linea
// 	guess, err := strconv.Atoi(input) //convert the input to integer
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if guess < target {
// 		fmt.Println("Bajito!!!")
// 	} else if guess > target {
// 		fmt.Println("Alto!!!")
// 	}

// }
