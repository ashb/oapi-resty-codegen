package generator

import (
	"fmt"
	"testing"

	"github.com/oapi-codegen/oapi-codegen/v2/pkg/codegen"
	"github.com/stretchr/testify/assert"
)

func TestConvertOperationWithTag(t *testing.T) {
	cases := []struct {
		tag      string
		opid     string
		expected string
	}{
		{"Asset Events", "GetAssetEventByAssetNameUri", "GetByAssetNameUri"},
		{"Variables", "GetVariable", "Get"},
		{"Pet", "FindPetsByStatus", "FindByStatus"},
		{"Pet", "FindPetsByStatus", "FindByStatus"},
		{"Abilities", "GetAbility", "Get"},
		{"Abilities", "AbilityGet", "Get"},
	}

	for _, test := range cases {
		t.Run(fmt.Sprintf("(%s,%s)", test.tag, test.opid), func(t *testing.T) {
			op := codegen.OperationDefinition{OperationId: test.opid}
			assert.Equal(t, convertOperationWithTag(test.tag, op), test.expected)
		})
	}
}
