package generator

import (
	"fmt"
	"iter"
	"maps"
	"strings"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/gertd/go-pluralize"
)

func operationsByTag(defs []codegen.OperationDefinition) map[string][]codegen.OperationDefinition {
	byTag := map[string][]codegen.OperationDefinition{}
	for _, def := range defs {
		tag := ""
		if len(def.Spec.Tags) >= 1 {
			tag = def.Spec.Tags[0]
		} else if len(def.Spec.Tags) == 0 {
			continue
		}
		byTag[tag] = append(byTag[tag], def)

	}
	return byTag
}

func untaggedOperations(defs []codegen.OperationDefinition) iter.Seq[codegen.OperationDefinition] {
	return func(yield func(codegen.OperationDefinition) bool) {
		for _, def := range defs {
			if len(def.Spec.Tags) > 0 {
				continue
			}
			if !yield(def) {
				return
			}
		}
	}
}

func tagToClass(tag string) string {
	return strings.ReplaceAll(cases.Title(language.English).String(tag), " ", "")
}

func convertOperationWithTag(tag string, op codegen.OperationDefinition) string {
	// Attempt to convert the following cases:
	//
	// - ("Asset Events", "GetAssetEventByAssetNameUri") -> "GetByAssetNameUri"
	// - ("Variables", "GetVariable") -> "Get"
	// - ("Pet", "FindPetsByStatus") -> "FindByStatus"
	// etc

	// TODO: Add config to supoprt
	// - ("Task Instances", "TiHeartbeat") -> "Heartbeat"

	class := tagToClass(tag)

	pluralize := pluralize.NewClient()
	// We always try the longer form
	class = pluralize.Plural(class)

	before, after, found := strings.Cut(op.OperationId, class)

	if !found && pluralize.IsPlural(class) {
		class = pluralize.Singular(class)
		before, after, found = strings.Cut(op.OperationId, class)
	}

	if !found {
		return op.OperationId
	}

	return before + after
}

func jsonTypeOrFirst[T any](items []T) T {
	var first T
	for i, item := range items {
		if i == 0 {
			first = item
		}
	}

	return first
}

func paramLocationToSetter(p codegen.ParameterDefinition) string {
	switch p.In {
	case "path":
		return "SetPathParam"
	case "header":
		return "SetHeader"
	case "cookie":
		return "SetQueryParam"
	default:
		panic(fmt.Sprintf("Unhandled parameter location %q", p.In))
	}
}

func init() {
	maps.Copy(codegen.TemplateFunctions, map[string]any{
		"operationsByTag":         operationsByTag,
		"untaggedOperations":      untaggedOperations,
		"tagToClass":              tagToClass,
		"convertOperationWithTag": convertOperationWithTag,
		"bestBody":                jsonTypeOrFirst[codegen.RequestBodyDefinition],
		"paramLocationToSetter":   paramLocationToSetter,
	})
}
