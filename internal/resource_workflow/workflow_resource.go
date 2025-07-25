package resource_workflow

import (
	"context"
	"fmt"
	"time"

	"github.com/aws-gopher/unstructured-sdk-go"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// WorkflowToModel converts an unstructured.Workflow to a WorkflowModel.
func WorkflowToModel(ctx context.Context, workflow *unstructured.Workflow, diagnostics diag.Diagnostics) *WorkflowModel {
	srcs, d := types.ListValueFrom(ctx, types.StringType, workflow.Sources)
	if d.HasError() {
		diagnostics.Append(d...)
	}

	dsts, d := types.ListValueFrom(ctx, types.StringType, workflow.Destinations)
	if d.HasError() {
		diagnostics.Append(d...)
	}

	// Convert WorkflowNodes
	var workflowNodesList types.List
	if len(workflow.WorkflowNodes) == 0 {
		workflowNodesList = types.ListNull(WorkflowNodesValue{}.Type(ctx))
	} else {
		workflowNodeValues := make([]attr.Value, 0, len(workflow.WorkflowNodes))
		for _, node := range workflow.WorkflowNodes {
			// Convert settings map to ObjectValue
			var settingsObj types.Object
			if len(node.Settings) == 0 {
				settingsObj, d = types.ObjectValue(map[string]attr.Type{}, map[string]attr.Value{})
				if d.HasError() {
					diagnostics.Append(d...)
				}
			} else {
				// Convert map[string]interface{} to map[string]attr.Value
				settingsMap := make(map[string]attr.Value)
				settingsAttrTypes := make(map[string]attr.Type)
				for k, v := range node.Settings {
					switch val := v.(type) {
					case string:
						settingsMap[k] = types.StringValue(val)
						settingsAttrTypes[k] = types.StringType
					case int:
						settingsMap[k] = types.Int64Value(int64(val))
						settingsAttrTypes[k] = types.Int64Type
					case float64:
						settingsMap[k] = types.Float64Value(val)
						settingsAttrTypes[k] = types.Float64Type
					case bool:
						settingsMap[k] = types.BoolValue(val)
						settingsAttrTypes[k] = types.BoolType
					default:
						// For complex types, convert to string representation
						settingsMap[k] = types.StringValue(fmt.Sprintf("%v", val))
						settingsAttrTypes[k] = types.StringType
					}
				}
				settingsObj, d = types.ObjectValue(settingsAttrTypes, settingsMap)
				if d.HasError() {
					diagnostics.Append(d...)
				}
			}

			// Create WorkflowNodesValue using constructor
			workflowNodeValue, d := NewWorkflowNodesValue(
				WorkflowNodesValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"id":       types.StringPointerValue(node.ID),
					"name":     types.StringValue(node.Name),
					"settings": settingsObj,
					"subtype":  types.StringValue(node.Subtype),
					"type":     types.StringValue(node.Type),
				},
			)
			if d.HasError() {
				diagnostics.Append(d...)
				continue // skip this node if conversion fails
			}

			// Convert to ObjectValue for the list
			workflowNodeObj, d := workflowNodeValue.ToObjectValue(ctx)
			if d.HasError() {
				diagnostics.Append(d...)
				continue // skip this node if conversion fails
			}
			workflowNodeValues = append(workflowNodeValues, workflowNodeObj)
		}
		workflowNodesList, d = types.ListValue(WorkflowNodesValue{}.Type(ctx), workflowNodeValues)
		if d.HasError() {
			diagnostics.Append(d...)
		}
	}

	// Convert Schedule from complex nested type to string
	var scheduleString types.String
	if workflow.Schedule == nil {
		scheduleString = types.StringNull()
	} else {
		// Convert the schedule to a string representation
		// This is a simplified conversion - you may need to adjust based on your needs
		if len(workflow.Schedule.CronTabEntries) > 0 {
			scheduleString = types.StringValue(workflow.Schedule.CronTabEntries[0].CronExpression)
		} else {
			scheduleString = types.StringNull()
		}
	}

	// Handle ReprocessAll pointer
	var reprocessAll types.Bool
	if workflow.ReprocessAll == nil {
		reprocessAll = types.BoolValue(false) // Use a default value instead of null
	} else {
		reprocessAll = types.BoolValue(*workflow.ReprocessAll)
	}

	// Handle WorkflowType
	var workflowType types.String
	if workflow.WorkflowType == nil {
		workflowType = types.StringValue("") // Use empty string as default
	} else {
		workflowType = types.StringValue(string(*workflow.WorkflowType))
	}

	// Handle DestinationId and SourceId
	// For now, we'll set them to null since they're not directly available in the API response
	// You may need to adjust this based on your specific requirements
	var destinationId types.String
	var sourceId types.String
	if len(workflow.Destinations) > 0 {
		destinationId = types.StringValue(workflow.Destinations[0])
	} else {
		destinationId = types.StringNull()
	}
	if len(workflow.Sources) > 0 {
		sourceId = types.StringValue(workflow.Sources[0])
	} else {
		sourceId = types.StringNull()
	}

	return &WorkflowModel{
		Id:            types.StringValue(workflow.ID),
		Name:          types.StringValue(workflow.Name),
		Sources:       srcs,
		Destinations:  dsts,
		CreatedAt:     types.StringValue(workflow.CreatedAt.Format(time.RFC3339)),
		ReprocessAll:  reprocessAll,
		Status:        types.StringValue(string(workflow.Status)),
		UpdatedAt:     types.StringValue(workflow.UpdatedAt.Format(time.RFC3339)),
		WorkflowType:  workflowType,
		WorkflowNodes: workflowNodesList,
		Schedule:      scheduleString,
		DestinationId: destinationId,
		SourceId:      sourceId,
	}
}
