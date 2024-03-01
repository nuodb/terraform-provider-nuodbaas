// (C) Copyright 2013-2024 Dassault Systemes SE.  All Rights Reserved.
//
// This software is licensed under a BSD 3-Clause License.
// See the LICENSE file provided with this software.

package framework

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AttributeBuilder struct {
	attributes map[string]schema.Attribute
}

func (ab *AttributeBuilder) WithStringAttribute(name, description string, optional bool) *AttributeBuilder {
	ab.attributes[name] = schema.StringAttribute{
		Description:         description,
		MarkdownDescription: description,
		Optional:            optional,
		Computed:            !optional,
	}
	return ab
}

func (ab *AttributeBuilder) WithOptionalStringAttribute(name, description string) *AttributeBuilder {
	return ab.WithStringAttribute(name, description, true)
}

func (ab *AttributeBuilder) WithComputedStringAttribute(name, description string) *AttributeBuilder {
	return ab.WithStringAttribute(name, description, false)
}

func (ab *AttributeBuilder) WithNameAttribute(typeName string) *AttributeBuilder {
	return ab.WithComputedStringAttribute("name", "The name of the "+typeName)
}

func (ab *AttributeBuilder) WithStringListAttribute(name, description string) *AttributeBuilder {
	ab.attributes[name] = schema.ListAttribute{
		Description:         description,
		MarkdownDescription: description,
		ElementType:         types.StringType,
		Optional:            true,
	}
	return ab
}

func (ab *AttributeBuilder) WithNewNestedAttribute(name, description string) *AttributeBuilder {
	childAttributes := make(map[string]schema.Attribute)
	ab.attributes[name] = schema.SingleNestedAttribute{
		Description:         description,
		MarkdownDescription: description,
		Attributes:          childAttributes,
		Optional:            true,
	}
	return &AttributeBuilder{childAttributes}
}

func (ab *AttributeBuilder) WithNewListNestedAttribute(name, description string) *AttributeBuilder {
	childAttributes := make(map[string]schema.Attribute)
	ab.attributes[name] = schema.ListNestedAttribute{
		Description:         description,
		MarkdownDescription: description,
		Computed:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: childAttributes,
		},
	}
	return &AttributeBuilder{childAttributes}
}

type SchemaBuilder struct {
	AttributeBuilder
	description string
}

func (sb *SchemaBuilder) WithDescription(description string) *SchemaBuilder {
	sb.description = description
	return sb
}

const (
	LABEL_FILTER_DESCRIPTION = "List of filters to apply based on labels, which are composed using `AND`. Acceptable filter expressions are:\n" +
		"  * `key` - Only return items that have label with specified key\n" +
		"  * `key=value` - Only return items that have label with specified key set to value\n" +
		"  * `!key` - Only return items that do _not_ have label with specified key\n" +
		"  * `key!=value` - Only return items that do _not_ have label with specified key set to value"
)

// WithOrganizationScopeFilters attaches common attributes for the filter nested
// attribute of an organization-scoped resource in DBaaS.
func (sb *SchemaBuilder) WithOrganizationScopeFilters(typeNamePlural string) *AttributeBuilder {
	return sb.WithNewNestedAttribute("filter", fmt.Sprintf("Filters to apply to %s", typeNamePlural)).
		WithStringListAttribute("labels", LABEL_FILTER_DESCRIPTION).
		WithOptionalStringAttribute("organization", fmt.Sprintf("The organization to filter %s on", typeNamePlural))
}

// WithProjectScopeFilters attaches schema information for the filter nested
// attribute of a project-scoped resource in DBaaS.
func (sb *SchemaBuilder) WithProjectScopeFilters(typeNamePlural string) *AttributeBuilder {
	return sb.WithOrganizationScopeFilters(typeNamePlural).
		WithOptionalStringAttribute("project", fmt.Sprintf("The project to filter %s on. If specified, the organization must also be specified.", typeNamePlural))
}

// WithOrganizationScopeList attaches schema information for the list attribute
// of the supplied type, which is organization-scoped (e.g. projects).
func (sb *SchemaBuilder) WithOrganizationScopeList(typeName, typeNamePlural string) *AttributeBuilder {
	return sb.WithNewListNestedAttribute(typeNamePlural, fmt.Sprintf("The list of %s that satisfy the filter requirements", typeNamePlural)).
		WithComputedStringAttribute("organization", fmt.Sprintf("The organization the %s belongs to", typeName))
}

// WithProjectScopeList attaches schema information for the list attribute of
// the supplied type, which is project-scoped (e.g. databases).
func (sb *SchemaBuilder) WithProjectScopeList(typeName, typeNamePlural string) *AttributeBuilder {
	return sb.WithNewListNestedAttribute(typeNamePlural, fmt.Sprintf("The list of %s that satisfy the filter requirements", typeNamePlural)).
		WithComputedStringAttribute("organization", fmt.Sprintf("The organization the %s belongs to", typeName)).
		WithComputedStringAttribute("project", fmt.Sprintf("The project the %s belongs to", typeName))
}

func (sb *SchemaBuilder) Build() *schema.Schema {
	return &schema.Schema{
		Description:         sb.description,
		MarkdownDescription: sb.description,
		Attributes:          sb.attributes,
	}
}

func NewSchemaBuilder() *SchemaBuilder {
	return &SchemaBuilder{
		AttributeBuilder: AttributeBuilder{make(map[string]schema.Attribute)},
	}
}
