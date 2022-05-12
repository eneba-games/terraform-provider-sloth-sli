package sli

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDataSourceSLI(t *testing.T) {
	successOutput := "something generated"

	tests := []struct {
		name            string
		sliConfigIn     string
		wantOutput      string
		expectError     *regexp.Regexp
		generatorMockFn func(t testing.TB, sliConfigInput string) *mockSliGenerator
	}{
		{
			name:        "generator fails",
			sliConfigIn: "some sli input config",
			wantOutput:  "",
			expectError: regexp.MustCompile("failed"),
			generatorMockFn: func(t testing.TB, sliConfigInput string) *mockSliGenerator {
				mockGenerator := newMockSliGenerator(t)

				mockGenerator.On("Generate", mock.Anything).Return("", errors.New("failed"))

				return mockGenerator
			},
		},
		{
			name:        "success",
			sliConfigIn: "some sli input config",
			wantOutput:  successOutput,
			expectError: nil,
			generatorMockFn: func(t testing.TB, sliConfigInput string) *mockSliGenerator {
				mockGenerator := newMockSliGenerator(t)

				mockGenerator.On(
					"Generate",
					mock.MatchedBy(func(fileName string) bool {
						b, err := ioutil.ReadFile(fileName)
						if err != nil || string(b) != sliConfigInput {
							return false
						}

						return true
					}),
				).Return(successOutput, nil)

				return mockGenerator
			},
		},
	}

	for _, tc := range tests {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mockGenerator := tc.generatorMockFn(t, tc.sliConfigIn)
			configureFunc := testConfigureProviderFunc(mockGenerator)
			configSrc := testConfig(tc.sliConfigIn)

			resource.UnitTest(t, resource.TestCase{
				ProviderFactories: map[string]func() (*schema.Provider, error){
					"sli": func() (*schema.Provider, error) {
						return Provider(configureFunc), nil
					},
				},
				Steps: []resource.TestStep{
					{
						Config:      configSrc,
						ExpectError: tc.expectError,
						Check: func(s *terraform.State) error {
							got := s.RootModule().Outputs["rendered"]

							assert.Equal(t, tc.wantOutput, got.Value)

							return nil
						},
					},
				},
			})
		})
	}
}

func testConfigureProviderFunc(mockGenerator *mockSliGenerator) schema.ConfigureContextFunc {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		var diags diag.Diagnostics

		return mockGenerator, diags
	}
}

func testConfig(sliConfigContent string) string {
	return fmt.Sprintf(`
		provider "sli" {
			sloth_path = "foo-bar-baz"
		}
		data "sli" "cfg" {
			sli_config = "%s"
		}
		output "rendered" {
			value = "${data.sli.cfg.rendered}"
		}`, sliConfigContent)
}
