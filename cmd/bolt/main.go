package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/adilead/bolt/internal/apprunner"
	"github.com/adilead/bolt/internal/bolt"
	"github.com/ktr0731/go-fuzzyfinder"
)

type generatorFunc func() ([]bolt.Searchable, error)
type optionExecutorFunc func(bolt.Searchable) error 

type BoltApp struct {
    args []string
    programs map[string]bool
    program string
    command string
    commandArgs []string

    optionGenerator map[string]generatorFunc
    optionExecutor map[string]optionExecutorFunc 
}



func NewBoltApp(args []string) (BoltApp, error) {
    return BoltApp {
        args: args, 
        programs: map[string]bool{"search":true, "run":true, "open":true, "browse":true},
        program: "",
        command: "",
        commandArgs: make([]string, 0),
        optionGenerator: map[string]generatorFunc{
            "run": apprunner.GetPrograms,
        },
        optionExecutor: map[string]optionExecutorFunc {
            "run": apprunner.RunProgram,
        },
    }, nil
}

func (b *BoltApp) Input() {
    program := ""
    if len(b.args) == 1 {
        program = b.args[0]
    } else if len(b.args) != 0 {
        log.Fatal("Call")
    }
    fmt.Printf("Bolt run -> %s ", program)
    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')
    if program == "" {
        inputs := strings.Split(input, " ")
        program = inputs[0]
        input = strings.Join(inputs[1:], "")
    }
}

func (b *BoltApp) GetOptions(cmd string) ([]bolt.Searchable, error){
    gen, ok := b.optionGenerator[cmd]
    if !ok {
        return nil, errors.New(fmt.Sprintf("%s is no valid command", cmd))
    }
    return gen()
}

func (b *BoltApp) Run(cmd string, choice bolt.Searchable) error {
    exec, ok := b.optionExecutor[cmd]
    if !ok {
        return errors.New(fmt.Sprintf("%s is no valid command", cmd))
    }
    return exec(choice)
}

func GetSearchUrl(cmd string) (string, error) {
    if cmd == "google" {
        return "https://google.com/search?q=", nil
    } else if cmd == "youtube" {
        return "https://youtube.com/search?q=", nil
    } else {
        return "", errors.New("No valid search cmd")    
    }
}

func fuzzyFind[I bolt.Searchable](options []I, promptString string) (I, error) {
    idx, err := fuzzyfinder.Find(options, 
        func (i int) string {return options[i].ToSlice()[0]},
        fuzzyfinder.WithPromptString(promptString))
    if err != nil {
        var result I
        return result, fmt.Errorf("fuzzy finder crashed %w", err)
    }
    return options[idx], nil
}

func main () {
    args := os.Args[1:]
    cmd := "run"
    if len(args) == 1 {
        cmd = args[0]
    } else if len(args) > 1 {
        log.Fatal("Usage: bolt [cmd]")
        return
    }     
    boltApp, err := NewBoltApp(args)

    options, err := boltApp.GetOptions(cmd)
    if err != nil {
        log.Fatal(err)
        return
    }

    choice, err := fuzzyFind(options, fmt.Sprintf("Bolt %s >",cmd))
    if err != nil {
        if errors.Unwrap(err) != fuzzyfinder.ErrAbort {
            log.Fatal(err)
        }
        return
    }
    fmt.Printf("You have chosen %s\n", choice.ToSlice()[0])

    //Run choice
    if err := boltApp.Run(cmd, choice); err != nil {
        log.Fatal(err)
    }
}
