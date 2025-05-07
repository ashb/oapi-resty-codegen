package openapi31downgrade

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type Suite struct {
	suite.Suite
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) TestDowngrade() {
	// Find the paths of all input files in the data directory.
	const suffix = "_in.yaml"
	paths, err := filepath.Glob(filepath.Join("testdata", "*"+suffix))
	if err != nil {
		s.T().Fatal(err)
	}
	specToMap := func(spec *openapi3.T) map[string]any {
		s.T().Helper()
		var out map[string]any
		bytes, err := yaml.Marshal(spec)
		s.Require().NoError(err)
		s.Require().NoError(yaml.Unmarshal(bytes, &out))
		return out
	}
	for _, path := range paths {
		path := path
		base := filepath.Base(path)
		testname := base[:len(base)-len(suffix)]
		s.Run(testname, func() {
			var convertedYAML, expectedYAML map[string]any

			sourceBytes, err := os.ReadFile(path)
			s.Require().NoError(err)

			loader := openapi3.NewLoader()

			inputSpec, err := loader.LoadFromData(sourceBytes)
			s.Require().NoError(err)

			expectedBytes, err := os.ReadFile(path[:len(path)-len(suffix)] + "_out.yaml")
			s.Require().NoError(err)
			err = yaml.Unmarshal(expectedBytes, &expectedYAML)
			s.Require().NoError(err)

			output, err := DowngradeTo3_0(inputSpec)
			s.Require().NoError(err)

			convertedYAML = specToMap(output)

			// We use gomega here as it gives better error indications
			if diff := cmp.Diff(convertedYAML, expectedYAML); diff != "" {
				s.Failf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
			// g.Expect(actual).Should(gomega.MatchYAML(expected))
			// s.EqualValues(expectedSpec, output)
		})
	}
}
