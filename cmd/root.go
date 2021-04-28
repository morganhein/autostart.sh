package cmd

import (
  "context"
  "fmt"
  "github.com/morganhein/autostart-sh/autostart"
  "github.com/spf13/cobra"
  "os"
  "time"
)

var rootCmd = &cobra.Command{
  Use:   "hugo",
  Short: "Hugo is a very fast static site generator",
  Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
  Run: func(cmd *cobra.Command, args []string) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
    defer cancel()
    err := autostart.Run(ctx, autostart.RunArgs{
     Cmd:  "echo",
     Args: []string{"sup playa!"},
     Sudo: false,
    })
    if err != nil {
     panic(err)
    }
    //err := autostart.Run(ctx, autostart.RunArgs{
    //  Cmd:  "apt",
    //  Args: []string{"search", "vim"},
    //  Sudo: true,
    //})
    //if err != nil {
    //  panic(err)
    //}
  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}


