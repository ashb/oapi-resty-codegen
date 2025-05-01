package cmds

import (
	"fmt"
	"os"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"github.com/oapi-codegen/oapi-codegen/v2/pkg/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type configuration struct {
	codegen.Configuration `yaml:",inline"`

	// OutputFile is the filename to output.
	OutputFile string `yaml:"output,omitempty"`
}

var (
	flagConfigFile string
	RootCmd        = &cobra.Command{
		Use:        "codegen [flags] <spec_file>",
		Short:      "Customized OpenAPI genration for Airflow's Go SDK",
		RunE:       Run,
		ArgAliases: []string{"spec-file"},
		Args:       cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	}
)

func Run(cmd *cobra.Command, args []string) error {
	configFile, err := os.ReadFile(flagConfigFile)
	if err != nil {
		return fmt.Errorf("error reading config file '%s': %w", flagConfigFile, err)
	}
	var opts configuration
	if err := yaml.UnmarshalStrict(configFile, &opts); err != nil {
		return fmt.Errorf("error parsing '%s' as YAML: %w", flagConfigFile, err)
	}

	// Now, ensure that the config options are valid.
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	overlayOpts := util.LoadSwaggerWithOverlayOpts{
		Path: opts.OutputOptions.Overlay.Path,
		// default to strict, but can be overridden
		Strict: true,
	}

	if opts.OutputOptions.Overlay.Strict != nil {
		overlayOpts.Strict = *opts.OutputOptions.Overlay.Strict
	}

	opts.AdditionalImports = append(opts.AdditionalImports, codegen.AdditionalImport{Package: "resty.dev/v3"})

	spec, err := util.LoadSwaggerWithOverlay(args[0], overlayOpts)
	if err != nil {
		return fmt.Errorf("error loading swagger spec in %s\n: %w", args[0], err)
	}
	code, err := codegen.Generate(spec, opts.Configuration)
	if err != nil {
		return fmt.Errorf("error generating code: %w", err)
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
