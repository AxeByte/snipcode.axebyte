package main

import (
  "os"

  log "github.com/sirupsen/logrus"
  "github.com/urfave/cli/v2"

  "github.com/AxeByte/snipcode.axebyte/internal/collector"
  "github.com/AxeByte/snipcode.axebyte/internal/config"
  "github.com/AxeByte/snipcode.axebyte/internal/writer"
)

func main() {
  app := &cli.App{
    Name:  "snipcode",
    Usage: "Collect & format code for LLM consumption",
    Commands: []*cli.Command{
      {Name: "init", Usage: "Generate .grepattern.yaml", Action: config.InitConfig},
      {Name: "compile", Usage: "Build compilation file", Flags: []cli.Flag{
         &cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: "Override output file"},
         &cli.BoolFlag{Name: "with-tree", Usage: "Append file tree listing"},
       }, Action: runCompile},
    },
  }
  if err := app.Run(os.Args); err != nil {
    log.Fatal(err)
  }
}

func runCompile(c *cli.Context) error {
  log.Info("Loading config from .grepattern.yaml")
  cfg, err := config.LoadConfig(".grepattern.yaml")
  if err != nil {
    return err
  }

  out := cfg.DefaultName
  if s := c.String("output"); s != "" {
    out = s
  }

  log.Info("Collecting files...")
  files, err := collector.Collect(".", cfg.IgnorePatterns)
  if err != nil {
    return err
  }
  log.Infof("Found %d files", len(files))

  if err := writer.Write(out, files, c.Bool("with-tree")); err != nil {
    return err
  }

  info, err := os.Stat(out)
  if err != nil {
    return err
  }
  log.WithFields(log.Fields{"output": out, "files": len(files), "size": info.Size()}).Info("Compilation complete")
  return nil
}
