package parser

import (
	"dependabot/internal/task/types"
	"io/ioutil"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/stretchr/testify/assert"
)

func TestTerraformParserEmptyFile(t *testing.T) {
	parser := CreateTerraformParser()
	dependenciesFromFile, _ := parser.ParseRequirementFile("")
	assert.Equal(t, 0, len(dependenciesFromFile))
}

func TestTerraformParserCustomFile(t *testing.T) {
	v1, _ := types.MakeVersion("v5.0.14")
	v2, _ := types.MakeVersion("v1.2.4")
	expectedDependencies := []types.Dependency{
		{
			SourceRaw:     "git::ssh://git@source.golabs.io/gopay_infra/terraform/aws-basic-instance.git?ref=v5.0.14",
			SourceBaseUrl: "source.golabs.io/gopay_infra/terraform/aws-basic-instance",
			Version:       *v1,
			Type:          "terraform",
		},
		{
			SourceRaw:     "git::ssh://git@source.golabs.io/gopay_infra/terraform/gcloud-postgresql.git?ref=v1.2.4",
			SourceBaseUrl: "source.golabs.io/gopay_infra/terraform/gcloud-postgresql",
			Version:       *v2,
			Type:          "terraform",
		},
	}

	parser := CreateTerraformParser()
	testFileRawContent, _ := ioutil.ReadFile("./test/terraform_test.tf")
	dependenciesFromFile, _ := parser.ParseRequirementFile(string(testFileRawContent))
	assert.Equal(t, expectedDependencies, dependenciesFromFile)
}

func TestParseTerraformSource(t *testing.T) {
	parser := CreateTerraformParser()

	source := "git::ssh://git@source.golabs.io/gopay_infra/terraform/gcloud-postgresql.git?ref=v9.10.15"
	parsedDependencyFromSource, err := parser.convertSourceToDependency(source)
	assert.Nil(t, err)
	v1, _ := types.MakeVersion("v9.10.15")
	expectedDependencies := types.Dependency{
		SourceRaw:     source,
		SourceBaseUrl: "source.golabs.io/gopay_infra/terraform/gcloud-postgresql",
		Type:          "terraform",
		Version:       *v1,
	}
	assert.Equal(t, expectedDependencies, parsedDependencyFromSource)

	source = "git::ssh://git@gitlab.com/test/terraform/gcloud-postgresql.git?ref=v9.10.15"
	_, err = parser.convertSourceToDependency(source)
	assert.NotNil(t, err)

	source = "git::ssh://git@source.golabs.io/gopay_infra/terraform/gcloud-postgresql.git"
	parsedDependencyFromSource, err = parser.convertSourceToDependency(source)
	assert.Nil(t, err)
	v2, _ := types.MakeVersion("latest")
	expectedDependencies = types.Dependency{
		SourceRaw:     source,
		SourceBaseUrl: "source.golabs.io/gopay_infra/terraform/gcloud-postgresql",
		Type:          "terraform",
		Version:       *v2,
	}
	assert.Equal(t, expectedDependencies, parsedDependencyFromSource)
}

func TestTerraformGetBaseUrlFromRawSource(t *testing.T) {
	parser := CreateTerraformParser()

	rawSource := "git::ssh://git@source.golabs.io/gopay_infra/terraform/gcloud-postgresql.git?ref=v1.2.4"
	expectedBaseUrl := "source.golabs.io/gopay_infra/terraform/gcloud-postgresql.git?ref=v1.2.4"
	parsedBaseUrl, _ := parser.GetBaseUrlFromRawSource(rawSource)
	assert.Equal(t, expectedBaseUrl, parsedBaseUrl)

	rawSource = faker.URL()
	_, err := parser.GetBaseUrlFromRawSource(rawSource)
	assert.NotNil(t, err)
}
