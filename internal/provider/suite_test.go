package provider

import (
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/stretchr/testify/assert"
)

type TfConfigBuilder struct {
	resources   map[string]any
	dataSources map[string]any
}

func NewTfConfigBuilder() *TfConfigBuilder {
	return &TfConfigBuilder{
		resources:   make(map[string]any),
		dataSources: make(map[string]any),
	}

}

func (b *TfConfigBuilder) WithDatabaseResource(name string, database *DatabaseResourceModel) *TfConfigBuilder {
	b.resources["nuodbaas_database "+name] = database
	return b
}

func (b *TfConfigBuilder) WithProjectResource(name string, project *ProjectResourceModel) *TfConfigBuilder {
	b.resources["nuodbaas_project "+name] = project
	return b
}

func (b *TfConfigBuilder) WithDatabaseDataSource(name string, database *DatabaseNameModel) *TfConfigBuilder {
	b.dataSources["nuodbaas_database "+name] = database
	return b
}

func (b *TfConfigBuilder) WithProjectDataSource(name string, project *ProjectNameModel) *TfConfigBuilder {
	b.dataSources["nuodbaas_project "+name] = project
	return b
}

func (b *TfConfigBuilder) Build() string {
	f := hclwrite.NewEmptyFile()
	for key, value := range b.resources {
		gohcl.EncodeIntoBody(value, f.Body().AppendNewBlock("resource", strings.Split(key, " ")).Body())
		f.Body().AppendNewline()
	}
	for key, value := range b.dataSources {
		gohcl.EncodeIntoBody(value, f.Body().AppendNewBlock("data", strings.Split(key, " ")).Body())
		f.Body().AppendNewline()
	}
	return string(f.Bytes())
}

func TestTfConfigBuilder(t *testing.T) {
	tfConfig := NewTfConfigBuilder().WithProjectResource("proj", &ProjectResourceModel{
		Organization: "org",
		Name:         "proj",
	}).WithDatabaseResource("db", &DatabaseResourceModel{
		Organization: "org",
		Project:      "proj",
		Name:         "db",
	}).WithProjectDataSource("proj", &ProjectNameModel{
		Organization: "org",
		Name:         "proj",
	}).WithDatabaseDataSource("db", &DatabaseNameModel{
		Organization: "org",
		Project:      "proj",
		Name:         "db",
	}).Build()

	// Check that project resource appears in config
	assert.Contains(t, tfConfig, `resource "nuodbaas_project" "proj" {
  organization = "org"
  name         = "proj"
  sla          = ""
  tier         = ""
}`)
	// Check that database resource appears in config
	assert.Contains(t, tfConfig, `resource "nuodbaas_database" "db" {
  organization = "org"
  project      = "proj"
  name         = "db"
`)
	// Check that project data source appears in config
	assert.Contains(t, tfConfig, `data "nuodbaas_project" "proj" {
  organization = "org"
  name         = "proj"
}`)
	// Check that database data source appears in config
	assert.Contains(t, tfConfig, `data "nuodbaas_database" "db" {
  organization = "org"
  project      = "proj"
  name         = "db"
`)
}
