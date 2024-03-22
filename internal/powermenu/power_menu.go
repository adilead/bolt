package powermenu

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"

	"github.com/adilead/bolt/internal/bolt"
)


type PowerMenuItem struct {
    name string 
    command string
}


func (pi *PowerMenuItem) ToSlice() []string {
    return []string{pi.name, pi.command}
}

func GetPrograms () ([]bolt.Searchable, error) {
    pms := []bolt.Searchable{
        &PowerMenuItem{"Lock", "i3lock -i ~/Images/2560x1440-dark-archlinux.png -C"},
        // PowerMenuItem{"Suspend", "i3lock && systemctl suspend"},
        &PowerMenuItem{"Power Off", "shutdown -h now --poweroff"},
        &PowerMenuItem{"Reboot", "shutdown -h now --reboot"},
    }
    return pms, nil
}



func RunProgram(choice bolt.Searchable) error {
    pi := choice.(*PowerMenuItem)
    commands := strings.Split(pi.command, " ")
    cmdPath, err := exec.LookPath(commands[0])
    if err != nil {
        return fmt.Errorf("powermenu fails due to %w", err)
    }
    cmd := exec.Command(cmdPath, commands[1:]...)
    cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
    err = cmd.Start()
    if err != nil {
        return fmt.Errorf("powermenu fails due to %w", err)
    }
    cmd.Wait()
    return nil
} 
