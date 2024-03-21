package windowctrl

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"github.com/adilead/bolt/internal/bolt"
)

type WindowItem struct {
    id string
    desktop string
    os string
    windowName string
}

func (wi *WindowItem) ToSlice() []string {
    return []string{wi.windowName, wi.desktop}
}

func GetPrograms () ([]bolt.Searchable, error) {
    wmctrlPath, err := exec.LookPath("wmctrl")
    if err != nil {
        return nil, fmt.Errorf("window_ctrl fails due to %w", err)
    }
    cmd := exec.Command(wmctrlPath, "-l")
    out, err := cmd.Output()
    if err != nil {
        return nil, fmt.Errorf("window_ctrl fails due to %w", err)
    }
    out_lines := strings.Split(string(out), "\n")
    mapFunc := func(x string) bolt.Searchable {
        s := strings.Fields(x)
        return &WindowItem{s[0], s[1], s[2], strings.Join(s[3:], " ")}
    }
    items := bolt.Map(out_lines[:len(out_lines)-1], mapFunc)
    return items, nil
}



func RunProgram(choice bolt.Searchable) error {
    wi := choice.(*WindowItem)
    wmctrlPath, err := exec.LookPath("wmctrl")
    if err != nil {
        return fmt.Errorf("window_ctrl fails due to %w", err)
    }
    cmd := exec.Command(wmctrlPath, "-a", wi.windowName)
    cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
    err = cmd.Start()
    if err != nil {
        return fmt.Errorf("window_ctrl fails due to %w", err)
    }
    cmd.Wait()
    return nil
} 
