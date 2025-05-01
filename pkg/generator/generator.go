package generator

import "github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"

func GenerateTypesGo() ([]byte, error) {
	ccfg := codegen.Configuration{
		PackageName: "a",
		Compatibility: codegen.CompatibilityOptions{
			AlwaysPrefixEnumValues: true,
		},
		Generate: codegen.GenerateOptions{
			Models: true,
		},
		OutputOptions: codegen.OutputOptions{
			SkipPrune: true,
			UserTemplates: map[string]string{
				"imports.tmpl": "A",
			},
		},
	}

	codegen.Generate(nil, ccfg)

	return nil, nil
}
