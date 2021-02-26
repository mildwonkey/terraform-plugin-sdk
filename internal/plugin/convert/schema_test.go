package convert

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/hashicorp/go-cty/cty"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/internal/configs/configschema"
)

var (
	equateEmpty   = cmpopts.EquateEmpty()
	typeComparer  = cmp.Comparer(cty.Type.Equals)
	valueComparer = cmp.Comparer(cty.Value.RawEquals)
)

// Test that we can convert configschema to protobuf types and back again.
func TestConvertSchemaBlocks(t *testing.T) {
	tests := map[string]struct {
		Block *tfprotov6.SchemaBlock
		Want  *configschema.Block
	}{
		"attributes": {
			&tfprotov6.SchemaBlock{
				Attributes: []*tfprotov6.SchemaAttribute{
					{
						Name: "computed",
						Type: tftypes.List{
							ElementType: tftypes.Bool,
						},
						Computed: true,
					},
					{
						Name:     "optional",
						Type:     tftypes.String,
						Optional: true,
					},
					{
						Name: "optional_computed",
						Type: tftypes.Map{
							AttributeType: tftypes.Bool,
						},
						Optional: true,
						Computed: true,
					},
					{
						Name:     "required",
						Type:     tftypes.Number,
						Required: true,
					},
				},
			},
			&configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"computed": {
						Type:     cty.List(cty.Bool),
						Computed: true,
					},
					"optional": {
						Type:     cty.String,
						Optional: true,
					},
					"optional_computed": {
						Type:     cty.Map(cty.Bool),
						Optional: true,
						Computed: true,
					},
					"required": {
						Type:     cty.Number,
						Required: true,
					},
				},
			},
		},
		"blocks": {
			&tfprotov6.SchemaBlock{
				BlockTypes: []*tfprotov6.SchemaNestedBlock{
					{
						TypeName: "list",
						Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
						Block:    &tfprotov6.SchemaBlock{},
					},
					{
						TypeName: "map",
						Nesting:  tfprotov6.SchemaNestedBlockNestingModeMap,
						Block:    &tfprotov6.SchemaBlock{},
					},
					{
						TypeName: "set",
						Nesting:  tfprotov6.SchemaNestedBlockNestingModeSet,
						Block:    &tfprotov6.SchemaBlock{},
					},
					{
						TypeName: "single",
						Nesting:  tfprotov6.SchemaNestedBlockNestingModeSingle,
						Block: &tfprotov6.SchemaBlock{
							Attributes: []*tfprotov6.SchemaAttribute{
								{
									Name:     "foo",
									Type:     tftypes.DynamicPseudoType,
									Required: true,
								},
							},
						},
					},
				},
			},
			&configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"list": {
						Nesting: configschema.NestingList,
					},
					"map": {
						Nesting: configschema.NestingMap,
					},
					"set": {
						Nesting: configschema.NestingSet,
					},
					"single": {
						Nesting: configschema.NestingSingle,
						Block: configschema.Block{
							Attributes: map[string]*configschema.Attribute{
								"foo": {
									Type:     cty.DynamicPseudoType,
									Required: true,
								},
							},
						},
					},
				},
			},
		},
		"deep block nesting": {
			&tfprotov6.SchemaBlock{
				BlockTypes: []*tfprotov6.SchemaNestedBlock{
					{
						TypeName: "single",
						Nesting:  tfprotov6.SchemaNestedBlockNestingModeSingle,
						Block: &tfprotov6.SchemaBlock{
							BlockTypes: []*tfprotov6.SchemaNestedBlock{
								{
									TypeName: "list",
									Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
									Block: &tfprotov6.SchemaBlock{
										BlockTypes: []*tfprotov6.SchemaNestedBlock{
											{
												TypeName: "set",
												Nesting:  tfprotov6.SchemaNestedBlockNestingModeSet,
												Block:    &tfprotov6.SchemaBlock{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			&configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"single": {
						Nesting: configschema.NestingSingle,
						Block: configschema.Block{
							BlockTypes: map[string]*configschema.NestedBlock{
								"list": {
									Nesting: configschema.NestingList,
									Block: configschema.Block{
										BlockTypes: map[string]*configschema.NestedBlock{
											"set": {
												Nesting: configschema.NestingSet,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			converted := ProtoToConfigSchema(tc.Block)
			if !cmp.Equal(converted, tc.Want, typeComparer, valueComparer, equateEmpty) {
				t.Fatal(cmp.Diff(converted, tc.Want, typeComparer, valueComparer, equateEmpty))
			}
		})
	}
}

// Test that we can convert configschema to protobuf types and back again.
func TestConvertProtoSchemaBlocks(t *testing.T) {
	tests := map[string]struct {
		Want  *tfprotov6.SchemaBlock
		Block *configschema.Block
	}{
		"attributes": {
			&tfprotov6.SchemaBlock{
				Attributes: []*tfprotov6.SchemaAttribute{
					{
						Name: "computed",
						Type: tftypes.List{
							ElementType: tftypes.Bool,
						},
						Computed: true,
					},
					{
						Name:     "optional",
						Type:     tftypes.String,
						Optional: true,
					},
					{
						Name: "optional_computed",
						Type: tftypes.Map{
							AttributeType: tftypes.Bool,
						},
						Optional: true,
						Computed: true,
					},
					{
						Name:     "required",
						Type:     tftypes.Number,
						Required: true,
					},
				},
			},
			&configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"computed": {
						Type:     cty.List(cty.Bool),
						Computed: true,
					},
					"optional": {
						Type:     cty.String,
						Optional: true,
					},
					"optional_computed": {
						Type:     cty.Map(cty.Bool),
						Optional: true,
						Computed: true,
					},
					"required": {
						Type:     cty.Number,
						Required: true,
					},
				},
			},
		},
		"blocks": {
			&tfprotov6.SchemaBlock{
				BlockTypes: []*tfprotov6.SchemaNestedBlock{
					{
						TypeName: "list",
						Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
						Block:    &tfprotov6.SchemaBlock{},
					},
					{
						TypeName: "map",
						Nesting:  tfprotov6.SchemaNestedBlockNestingModeMap,
						Block:    &tfprotov6.SchemaBlock{},
					},
					{
						TypeName: "set",
						Nesting:  tfprotov6.SchemaNestedBlockNestingModeSet,
						Block:    &tfprotov6.SchemaBlock{},
					},
					{
						TypeName: "single",
						Nesting:  tfprotov6.SchemaNestedBlockNestingModeSingle,
						Block: &tfprotov6.SchemaBlock{
							Attributes: []*tfprotov6.SchemaAttribute{
								{
									Name:     "foo",
									Type:     tftypes.DynamicPseudoType,
									Required: true,
								},
							},
						},
					},
				},
			},
			&configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"list": {
						Nesting: configschema.NestingList,
					},
					"map": {
						Nesting: configschema.NestingMap,
					},
					"set": {
						Nesting: configschema.NestingSet,
					},
					"single": {
						Nesting: configschema.NestingSingle,
						Block: configschema.Block{
							Attributes: map[string]*configschema.Attribute{
								"foo": {
									Type:     cty.DynamicPseudoType,
									Required: true,
								},
							},
						},
					},
				},
			},
		},
		"deep block nesting": {
			&tfprotov6.SchemaBlock{
				BlockTypes: []*tfprotov6.SchemaNestedBlock{
					{
						TypeName: "single",
						Nesting:  tfprotov6.SchemaNestedBlockNestingModeSingle,
						Block: &tfprotov6.SchemaBlock{
							BlockTypes: []*tfprotov6.SchemaNestedBlock{
								{
									TypeName: "list",
									Nesting:  tfprotov6.SchemaNestedBlockNestingModeList,
									Block: &tfprotov6.SchemaBlock{
										BlockTypes: []*tfprotov6.SchemaNestedBlock{
											{
												TypeName: "set",
												Nesting:  tfprotov6.SchemaNestedBlockNestingModeSet,
												Block:    &tfprotov6.SchemaBlock{},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			&configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"single": {
						Nesting: configschema.NestingSingle,
						Block: configschema.Block{
							BlockTypes: map[string]*configschema.NestedBlock{
								"list": {
									Nesting: configschema.NestingList,
									Block: configschema.Block{
										BlockTypes: map[string]*configschema.NestedBlock{
											"set": {
												Nesting: configschema.NestingSet,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			converted := ConfigSchemaToProto(tc.Block)
			if !cmp.Equal(converted, tc.Want, typeComparer, equateEmpty) {
				t.Fatal(cmp.Diff(converted, tc.Want, typeComparer, equateEmpty))
			}
		})
	}
}
