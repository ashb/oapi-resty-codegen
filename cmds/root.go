package cmds

import (
	"fmt"
	"os"

	"github.com/ashb/oapi-resty-codegen/pkg/generator"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	flagConfigFile string
	RootCmd        = &cobra.Command{
		Use:        "codegen [flags] <spec_file>",
		Short:      "Customized OpenAPI genration for Airflow's Go SDK",
		RunE:       Run,
		ArgAliases: []string{"spec-file"},
		Args:       cobra.MatchAll(cobra.MaximumNArgs(1), cobra.OnlyValidArgs),
	}
)

func Run(cmd *cobra.Command, args []string) error {
	configFile, err := os.ReadFile(flagConfigFile)
	if err != nil {
		return fmt.Errorf("error reading config file '%s': %w", flagConfigFile, err)
	}
	var opts generator.GenerateArgs
	if err := yaml.UnmarshalStrict(configFile, &opts); err != nil {
		return fmt.Errorf("error parsing '%s' as YAML: %w", flagConfigFile, err)
	}

	if len(args) > 0 {
		opts.Input = args[0]
	}

	code, err := generator.Generate(opts)
	if err != nil {
		return err
	}

	if opts.OutputFile != "" {
		err = os.WriteFile(opts.OutputFile, []byte(code), 0o644)
		if err != nil {
			return fmt.Errorf("error writing generated code to file: %w", err)
		}
	} else {
		fmt.Print(code)
	}
	return nil
}

func init() {
	RootCmd.PersistentFlags().
		StringVar(&flagConfigFile, "config", "", "A YAML config file that controls oapi-codegen behavior.")
	RootCmd.MarkPersistentFlagFilename("config", ".yaml", ".yml")
}
