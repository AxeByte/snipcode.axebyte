package main

import (
	"os"
	"path/filepath"

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
			{Name: "init-admin", Usage: "Generate ~/.config/snipcode/.grepattern.yaml", Action: config.InitAdmin},
			{
				Name:  "compile",
				Usage: "Build compilation file",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "output", Aliases: []string{"o"}, Usage: "Override output file"},
					&cli.BoolFlag{Name: "with-tree", Usage: "Append file tree listing"},
				},
				Action: runCompile,
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func runCompile(c *cli.Context) error {
	// determine config paths
	var cfgMerged *config.Config
	// global
	if cfgDir, err := os.UserConfigDir(); err == nil {
		adminPath := filepath.Join(cfgDir, "snipcode", ".grepattern.yaml")
		if _, err := os.Stat(adminPath); err == nil {
			log.Infof("Loading global config from %s", adminPath)
			globalCfg, err := config.LoadConfig(adminPath)
			if err != nil {
				return err
			}
			cfgMerged = globalCfg
		}
	}
	// local
	localPath := ".grepattern.yaml"
	if _, err := os.Stat(localPath); err == nil {
		log.Infof("Loading local config from %s", localPath)
		localCfg, err := config.LoadConfig(localPath)
		if err != nil {
			return err
		}
		if cfgMerged == nil {
			cfgMerged = localCfg
		} else {
			// merge: local overrides default and adds ignore patterns
			cfgMerged.DefaultName = localCfg.DefaultName
			cfgMerged.IgnorePatterns = append(cfgMerged.IgnorePatterns, localCfg.IgnorePatterns...)
		}
	}
	// fallback if no config found
	if cfgMerged == nil {
		log.Info("No config found; using defaults")
		cfgMerged = &config.Config{
			DefaultName:    "comp_code.txt",
			IgnorePatterns: []string{},
		}
	}

	out := cfgMerged.DefaultName
	if s := c.String("output"); s != "" {
		out = s
	}

	log.Infof("Collecting files (skipping any named %q)...", out)
	files, err := collector.Collect(".", cfgMerged.IgnorePatterns)
	if err != nil {
		return err
	}

	// filter out the output file itself
	filtered := make([]string, 0, len(files))
	for _, f := range files {
		if filepath.Clean(f) == filepath.Clean(out) {
			log.Infof("Skipping output file from collection: %s", f)
			continue
		}
		filtered = append(filtered, f)
	}
	files = filtered

	log.Infof("Found %d files", len(files))

	if err := writer.Write(out, files, c.Bool("with-tree")); err != nil {
		return err
	}

	info, err := os.Stat(out)
	if err != nil {
		return err
	}
	log.WithFields(log.Fields{
		"output": out,
		"files":  len(files),
		"size":   info.Size(),
	}).Info("Compilation complete")
	return nil
}
