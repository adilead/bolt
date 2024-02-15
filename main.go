package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
    "log"
)

// application -> command -> arguments,
/*
e.g.    search google blablabla
        search youtube blalbala
        browse url https://youtube.com
        browse bookmark youtube
        run evince //here the command become the name of program
        run echo 1 //the arguments are the arguments of the program

        for browser features the browser must be set previously or somehow find the system browser

        use wmctrl for window management
        when using the browser:
        focus a browser windows first
        then execute the run
*/

func main () {
    args := os.Args[1:]
    mode := ""
    if len(args) != 0 {
        mode = args[0]
    }
    applications := []string{"search", "run", "open", "browse"}
    _ = applications
    fmt.Printf("Bolt run -> %s ", mode)
    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')
    if mode == "" {
        inputs := strings.Split(input, " ")
        mode = inputs[0]
        input = strings.Join(inputs[1:], "")
    }
    if mode == "google" {
        input = strings.Replace(input, " ", "+", -1)
        cmd := exec.Command("brave", "--new-tab", "https://google.com/search?q=" + input)
        cmd.Run() //blocking
    } else {
        log.Fatal("No valid mode")
    }
}
