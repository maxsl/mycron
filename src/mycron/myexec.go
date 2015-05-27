package mycron

import (
	//"fmt"
	"os/exec"
	"time"
	"os"
	"strings"
	"errors"
    "bytes"
    "runtime"
)

func ExecWithTimeout(d time.Duration, line string)(string, error) {
    var cmd * exec.Cmd
    if runtime.GOOS == "windows"{
        cmd = exec.Command("cmd", "/C", line)
    }else {
        shell := os.Getenv("SHELL")
        cmd = exec.Command(shell, "-c", line)
    }
    var out bytes.Buffer
    cmd.Stdout = &out
	if err := cmd.Start(); err != nil {
		return "" , err
	}
	if d <= 0 {
		cmd.Wait()
		return strings.TrimSpace(string(out.String())),err
	}

	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()
	//fmt.Println(cmd.Process.Pid)
	select {
	case <-time.After(d):
		cmd.Process.Kill()
		return "",errors.New("time out")
	case  err :=<-done:
		if err !=nil{
			return "",err
		}
		return strings.TrimSpace(string(out.String())),err
	}
}

func ShellRun(line string) (string, error) {
	shell := os.Getenv("SHELL")
	b, err := exec.Command(shell, "-c", line).Output()
	if err != nil {
		return "", errors.New(err.Error() + ":" + string(b))
	}
	return strings.TrimSpace(string(b)), nil
}