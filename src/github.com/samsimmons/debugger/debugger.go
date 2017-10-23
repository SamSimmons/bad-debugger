package main

import (
  "flag"
  "log"
  "os"
  "os/exec"
  "syscall"
)

func main() {
  flag.Parse()
  input := flag.Arg(0)
  cmd := exec.Command(input)
  cmd.Args = []string{input}
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  cmd.SysProcAttr = &syscall.SysProcAttr{Ptrace: true}
  err := cmd.Start()
  if err != nil {
    log.Fatal(err)
  }
  err = cmd.Wait()
  log.Printf("State: %v\n", err)
  log.Println("Restarting...")
  err = syscall.PtraceCont(cmd.Process.Pid, 0)
  if err != nil {
    log.Panic(err)
  }
  var ws syscall.WaitStatus
  _, err = syscall.Wait4(cmd.Process.Pid, &ws, syscall.WALL, nil)
  if err != nil {
    log.Fatal(err)
  }
  log.Printf("Exited: %v\n", ws.Exited())
  log.Printf("Exit status: %v\n", ws.ExitStatus())
}
