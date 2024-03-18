package apprunner

import (
	"fmt"
	// "os"
	"os/exec"
	"slices"
	"strings"

	"syscall"

	"github.com/adilead/bolt/internal/bolt"
)

type ProgramItem struct {
    Name string
    Description string
    Exec []string
}

func (pi *ProgramItem) ToSlice() []string {
    return []string{pi.Name, pi.Description}
}

func GetPrograms () ([]bolt.Searchable, error) {
    dirs := dataDirs()
    entries, err := Scan(dirs)
    if err != nil {
        return nil, fmt.Errorf("Error when getting programs: %w", err)
    }
    entries_flat := slices.Concat(entries...)
    pis := make([]bolt.Searchable, len(entries_flat))
    for i,e := range entries_flat {
        e.ExpandExec("")
        args := strings.Split(strings.TrimSpace(e.Exec), " ")
        p, err := exec.LookPath(args[0])
        if err == nil {
            args[0] = p
        }
        pis[i] = &ProgramItem{e.Name, "Description", args} 
    }
    return pis, nil
}



func RunProgram(choice bolt.Searchable) error {
    prog := choice.(*ProgramItem) 
    cmd := exec.Command(prog.Exec[0], prog.Exec[1:]...)
    cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
    err := cmd.Start()
    if err != nil {
        panic(err)
    }
    return nil
} 
