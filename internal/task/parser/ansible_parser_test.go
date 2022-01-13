package parser

import (
	"dependabot/internal/task/types"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnsibleParserEmptyFile(t *testing.T) {
	parser := CreateAnsibleParser()
	dependenciesFromFile, _ := parser.ParseRequirementFile("")
	assert.Equal(t, 0, len(dependenciesFromFile))
}

func TestAnsibleParserCustomFile(t *testing.T) {
	v1, _ := types.MakeVersion("v1.2.0")
	v2, _ := types.MakeVersion("v1.0.6")
	v3, _ := types.MakeVersion("v2.0.0")
	v4, _ := types.MakeVersion("latest")
	expectedDependencies := []types.Dependency{
		{
			SourceRaw:     "https://source.golabs.io/gopay_infra/ansible/playbooks/basic-instance-playbook.git",
			SourceBaseUrl: "source.golabs.io/gopay_infra/ansible/playbooks/basic-instance-playbook",
			Version:       *v1,
			Type:          "ansible",
		},
		{
			SourceRaw:     "https://source.golabs.io/gopay_infra/ansible/playbooks/rabbitmq-playbook.git",
			SourceBaseUrl: "source.golabs.io/gopay_infra/ansible/playbooks/rabbitmq-playbook",
			Version:       *v2,
			Type:          "ansible",
		},
		{
			SourceRaw:     "git@source.golabs.io:farras.f.interns/rabbitmq-role.git",
			SourceBaseUrl: "source.golabs.io/farras.f.interns/rabbitmq-role",
			Version:       *v3,
			Type:          "ansible",
		},
		{
			SourceRaw:     "git@source.golabs.io:farras.f.interns/rabbitmq-exporter.git",
			SourceBaseUrl: "source.golabs.io/farras.f.interns/rabbitmq-exporter",
			Version:       *v4,
			Type:          "ansible",
		},
	}

	parser := CreateAnsibleParser()
	testFileRawContent, _ := ioutil.ReadFile("./test/ansible_test.yaml")
	dependenciesFromFile, _ := parser.ParseRequirementFile(string(testFileRawContent))
	assert.Equal(t, expectedDependencies, dependenciesFromFile)
}

func TestAnsibleGetBaseUrlFromRawSource(t *testing.T) {
	parser := CreateAnsibleParser()

	rawSource := "https://source.golabs.io/gopay_infra/ansible/playbooks/scp-playbook.git"
	expectedBaseUrl := "source.golabs.io/gopay_infra/ansible/playbooks/scp-playbook"
	parsedBaseUrl, _ := parser.GetBaseUrlFromRawSource(rawSource)
	assert.Equal(t, expectedBaseUrl, parsedBaseUrl)
}
