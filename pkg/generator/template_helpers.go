package generator

import (
	"maps"
	"strings"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func operationsByTag(defs []codegen.OperationDefinition) map[string][]codegen.OperationDefinition {
	byTag := map[string][]codegen.OperationDefinition{}
	for _, def := range defs {
		tag := ""
		if len(def.Spec.Tags) >= 1 {
			tag = def.Spec.Tags[0]
		}
		byTag[tag] = append(byTag[tag], def)

	}
	return byTag
}

func tagToClass(tag string) string {
	return strings.ReplaceAll(cases.Title(language.English).String(tag), " ", "")
}

func convertOperationWithTag(tag string, op codegen.OperationDefinition) string {
	// Attempt to convert the following cases:
	//
	// - ("Asset Events", "GetAssetEventByAssetNameUri") -> "GetByAssetNameUri"
	// - ("Variables", "GetVariable") -> "Get"

	// TODO: Add config to supoprt
	// - ("Task Instances", "TiHeartbeat") -> "Heartbeat"

	class := tagToClass(tag)

	before, after, found := strings.Cut(op.OperationId, class)

	if !found && strings.HasSuffix(class, "s") {
		before, after, found = strings.Cut(op.OperationId, strings.TrimSuffix(class, "s"))
	}

	if !found {
		return op.OperationId
	}

	return before + after
}

func init() {
	maps.Copy(codegen.TemplateFunctions, map[string]any{
		"operationsByTag":         operationsByTag,
		"tagToClass":              tagToClass,
		"convertOperationWithTag": convertOperationWithTag,
	})
}
