package autostart

import (
    "bytes"
    "context"
    "fmt"
    "os/exec"
)

/*
This should:
1. contain a map of names to shell scripts and their installation templates
2. install/uninstall the scripts by name
3. be able to install new templates via some import process
*/

/* 
installation steps:
1. script name/template lookup
2. downlaod the script locally
3. execute the script using the appropriate shell
*/

type RunArgs struct {
    Shell bool //is this a shell command?
    Cmd string
    Args []string
    Sudo bool
}


func Run(ctx context.Context, run RunArgs) error {
    var args []string
    if run.Shell {
        args = append(args, "/usr/bin/bash", "-c")
    }
    if run.Sudo {
        args = append(args, "sudo")
    } 
    args = append(args, run.Cmd)
    for _, v := range run.Args {
        args = append(args, v)
    }
    //Run the cmd
    cmd := exec.Command(args[0], args[1:]...)
    var out bytes.Buffer
	cmd.Stdout = &out
    err := cmd.Run()
    if err != nil {
        panic(err)
    }
    fmt.Println("output: ", out.String())
    return nil
}