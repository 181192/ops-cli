package config

import (
	"encoding/json"
	"fmt"

	"github.com/181192/ops-cli/pkg/cmd/cmdutils"

	logger "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

var (
	format string
)

func showConfigCmd(cmd *cmdutils.Cmd) {

	cmd.CobraCommand.Use = "show"
	cmd.CobraCommand.Short = "Show current config"
	cmd.CobraCommand.Long = ""
	cmd.CobraCommand.Run = func(_ *cobra.Command, args []string) {
		if format != "yaml" && format != "json" {
			logger.Fatal("unsupported format")
		}

		c := viper.AllSettings()
		bs, err := yaml.Marshal(c)
		if err != nil {
			logger.Fatalf("unable to marshal config to YAML: %v", err)
		}

		if format == "yaml" {
			fmt.Printf("# Config file used %s\n", viper.ConfigFileUsed())
			fmt.Print(string(bs))
			return
		}

		var body interface{}
		if err := yaml.Unmarshal([]byte(bs), &body); err != nil {
			logger.Fatalf("unable to convert YAML to JSON: %v", err)
		}

		body = convert(body)

		b, err := json.MarshalIndent(body, "", " ")
		if err != nil {
			logger.Fatalf("unable to marshal config to JSON: %v", err)
		}

		fmt.Print(string(b))

	}

	cmd.FlagSetGroup.InFlagSet("Show config", func(fs *pflag.FlagSet) {
		fs.StringVar(&format, "format", "yaml", "Output format [yaml|json]")
	})
}

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}
