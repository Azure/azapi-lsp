//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armlogic

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// WorkflowRunActionsClient contains the methods for the WorkflowRunActions group.
// Don't use this type directly, use NewWorkflowRunActionsClient() instead.
type WorkflowRunActionsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewWorkflowRunActionsClient creates a new instance of WorkflowRunActionsClient with the specified values.
//   - subscriptionID - The subscription id.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewWorkflowRunActionsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*WorkflowRunActionsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &WorkflowRunActionsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Get - Gets a workflow run action.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-05-01
//   - resourceGroupName - The resource group name.
//   - workflowName - The workflow name.
//   - runName - The workflow run name.
//   - actionName - The workflow action name.
//   - options - WorkflowRunActionsClientGetOptions contains the optional parameters for the WorkflowRunActionsClient.Get method.
func (client *WorkflowRunActionsClient) Get(ctx context.Context, resourceGroupName string, workflowName string, runName string, actionName string, options *WorkflowRunActionsClientGetOptions) (WorkflowRunActionsClientGetResponse, error) {
	var err error
	const operationName = "WorkflowRunActionsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, workflowName, runName, actionName, options)
	if err != nil {
		return WorkflowRunActionsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return WorkflowRunActionsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return WorkflowRunActionsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *WorkflowRunActionsClient) getCreateRequest(ctx context.Context, resourceGroupName string, workflowName string, runName string, actionName string, options *WorkflowRunActionsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/workflows/{workflowName}/runs/{runName}/actions/{actionName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workflowName == "" {
		return nil, errors.New("parameter workflowName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workflowName}", url.PathEscape(workflowName))
	if runName == "" {
		return nil, errors.New("parameter runName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{runName}", url.PathEscape(runName))
	if actionName == "" {
		return nil, errors.New("parameter actionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{actionName}", url.PathEscape(actionName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *WorkflowRunActionsClient) getHandleResponse(resp *http.Response) (WorkflowRunActionsClientGetResponse, error) {
	result := WorkflowRunActionsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.WorkflowRunAction); err != nil {
		return WorkflowRunActionsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Gets a list of workflow run actions.
//
// Generated from API version 2019-05-01
//   - resourceGroupName - The resource group name.
//   - workflowName - The workflow name.
//   - runName - The workflow run name.
//   - options - WorkflowRunActionsClientListOptions contains the optional parameters for the WorkflowRunActionsClient.NewListPager
//     method.
func (client *WorkflowRunActionsClient) NewListPager(resourceGroupName string, workflowName string, runName string, options *WorkflowRunActionsClientListOptions) *runtime.Pager[WorkflowRunActionsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[WorkflowRunActionsClientListResponse]{
		More: func(page WorkflowRunActionsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *WorkflowRunActionsClientListResponse) (WorkflowRunActionsClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "WorkflowRunActionsClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, resourceGroupName, workflowName, runName, options)
			}, nil)
			if err != nil {
				return WorkflowRunActionsClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *WorkflowRunActionsClient) listCreateRequest(ctx context.Context, resourceGroupName string, workflowName string, runName string, options *WorkflowRunActionsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/workflows/{workflowName}/runs/{runName}/actions"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workflowName == "" {
		return nil, errors.New("parameter workflowName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workflowName}", url.PathEscape(workflowName))
	if runName == "" {
		return nil, errors.New("parameter runName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{runName}", url.PathEscape(runName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-01")
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *WorkflowRunActionsClient) listHandleResponse(resp *http.Response) (WorkflowRunActionsClientListResponse, error) {
	result := WorkflowRunActionsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.WorkflowRunActionListResult); err != nil {
		return WorkflowRunActionsClientListResponse{}, err
	}
	return result, nil
}

// NewListExpressionTracesPager - Lists a workflow run expression trace.
//
// Generated from API version 2019-05-01
//   - resourceGroupName - The resource group name.
//   - workflowName - The workflow name.
//   - runName - The workflow run name.
//   - actionName - The workflow action name.
//   - options - WorkflowRunActionsClientListExpressionTracesOptions contains the optional parameters for the WorkflowRunActionsClient.NewListExpressionTracesPager
//     method.
func (client *WorkflowRunActionsClient) NewListExpressionTracesPager(resourceGroupName string, workflowName string, runName string, actionName string, options *WorkflowRunActionsClientListExpressionTracesOptions) *runtime.Pager[WorkflowRunActionsClientListExpressionTracesResponse] {
	return runtime.NewPager(runtime.PagingHandler[WorkflowRunActionsClientListExpressionTracesResponse]{
		More: func(page WorkflowRunActionsClientListExpressionTracesResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *WorkflowRunActionsClientListExpressionTracesResponse) (WorkflowRunActionsClientListExpressionTracesResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "WorkflowRunActionsClient.NewListExpressionTracesPager")
			req, err := client.listExpressionTracesCreateRequest(ctx, resourceGroupName, workflowName, runName, actionName, options)
			if err != nil {
				return WorkflowRunActionsClientListExpressionTracesResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return WorkflowRunActionsClientListExpressionTracesResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return WorkflowRunActionsClientListExpressionTracesResponse{}, runtime.NewResponseError(resp)
			}
			return client.listExpressionTracesHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listExpressionTracesCreateRequest creates the ListExpressionTraces request.
func (client *WorkflowRunActionsClient) listExpressionTracesCreateRequest(ctx context.Context, resourceGroupName string, workflowName string, runName string, actionName string, options *WorkflowRunActionsClientListExpressionTracesOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/workflows/{workflowName}/runs/{runName}/actions/{actionName}/listExpressionTraces"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if workflowName == "" {
		return nil, errors.New("parameter workflowName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{workflowName}", url.PathEscape(workflowName))
	if runName == "" {
		return nil, errors.New("parameter runName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{runName}", url.PathEscape(runName))
	if actionName == "" {
		return nil, errors.New("parameter actionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{actionName}", url.PathEscape(actionName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listExpressionTracesHandleResponse handles the ListExpressionTraces response.
func (client *WorkflowRunActionsClient) listExpressionTracesHandleResponse(resp *http.Response) (WorkflowRunActionsClientListExpressionTracesResponse, error) {
	result := WorkflowRunActionsClientListExpressionTracesResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ExpressionTraces); err != nil {
		return WorkflowRunActionsClientListExpressionTracesResponse{}, err
	}
	return result, nil
}