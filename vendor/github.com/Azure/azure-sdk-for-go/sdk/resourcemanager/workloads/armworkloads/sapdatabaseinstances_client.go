//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armworkloads

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// SAPDatabaseInstancesClient contains the methods for the SAPDatabaseInstances group.
// Don't use this type directly, use NewSAPDatabaseInstancesClient() instead.
type SAPDatabaseInstancesClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewSAPDatabaseInstancesClient creates a new instance of SAPDatabaseInstancesClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewSAPDatabaseInstancesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*SAPDatabaseInstancesClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &SAPDatabaseInstancesClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginCreate - Creates the Database resource corresponding to the Virtual Instance for SAP solutions resource.
// This will be used by service only. PUT by end user will return a Bad Request error.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-04-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - sapVirtualInstanceName - The name of the Virtual Instances for SAP solutions resource
//   - databaseInstanceName - Database resource name string modeled as parameter for auto generation to work correctly.
//   - body - Request body of Database resource of a SAP system.
//   - options - SAPDatabaseInstancesClientBeginCreateOptions contains the optional parameters for the SAPDatabaseInstancesClient.BeginCreate
//     method.
func (client *SAPDatabaseInstancesClient) BeginCreate(ctx context.Context, resourceGroupName string, sapVirtualInstanceName string, databaseInstanceName string, body SAPDatabaseInstance, options *SAPDatabaseInstancesClientBeginCreateOptions) (*runtime.Poller[SAPDatabaseInstancesClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, resourceGroupName, sapVirtualInstanceName, databaseInstanceName, body, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[SAPDatabaseInstancesClientCreateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[SAPDatabaseInstancesClientCreateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Create - Creates the Database resource corresponding to the Virtual Instance for SAP solutions resource.
// This will be used by service only. PUT by end user will return a Bad Request error.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-04-01
func (client *SAPDatabaseInstancesClient) create(ctx context.Context, resourceGroupName string, sapVirtualInstanceName string, databaseInstanceName string, body SAPDatabaseInstance, options *SAPDatabaseInstancesClientBeginCreateOptions) (*http.Response, error) {
	var err error
	const operationName = "SAPDatabaseInstancesClient.BeginCreate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createCreateRequest(ctx, resourceGroupName, sapVirtualInstanceName, databaseInstanceName, body, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createCreateRequest creates the Create request.
func (client *SAPDatabaseInstancesClient) createCreateRequest(ctx context.Context, resourceGroupName string, sapVirtualInstanceName string, databaseInstanceName string, body SAPDatabaseInstance, options *SAPDatabaseInstancesClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Workloads/sapVirtualInstances/{sapVirtualInstanceName}/databaseInstances/{databaseInstanceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if sapVirtualInstanceName == "" {
		return nil, errors.New("parameter sapVirtualInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sapVirtualInstanceName}", url.PathEscape(sapVirtualInstanceName))
	if databaseInstanceName == "" {
		return nil, errors.New("parameter databaseInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseInstanceName}", url.PathEscape(databaseInstanceName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, body); err != nil {
		return nil, err
	}
	return req, nil
}

// BeginDelete - Deletes the Database resource corresponding to a Virtual Instance for SAP solutions resource.
// This will be used by service only. Delete by end user will return a Bad Request error.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-04-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - sapVirtualInstanceName - The name of the Virtual Instances for SAP solutions resource
//   - databaseInstanceName - Database resource name string modeled as parameter for auto generation to work correctly.
//   - options - SAPDatabaseInstancesClientBeginDeleteOptions contains the optional parameters for the SAPDatabaseInstancesClient.BeginDelete
//     method.
func (client *SAPDatabaseInstancesClient) BeginDelete(ctx context.Context, resourceGroupName string, sapVirtualInstanceName string, databaseInstanceName string, options *SAPDatabaseInstancesClientBeginDeleteOptions) (*runtime.Poller[SAPDatabaseInstancesClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, sapVirtualInstanceName, databaseInstanceName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[SAPDatabaseInstancesClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[SAPDatabaseInstancesClientDeleteResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Delete - Deletes the Database resource corresponding to a Virtual Instance for SAP solutions resource.
// This will be used by service only. Delete by end user will return a Bad Request error.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-04-01
func (client *SAPDatabaseInstancesClient) deleteOperation(ctx context.Context, resourceGroupName string, sapVirtualInstanceName string, databaseInstanceName string, options *SAPDatabaseInstancesClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	const operationName = "SAPDatabaseInstancesClient.BeginDelete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, sapVirtualInstanceName, databaseInstanceName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *SAPDatabaseInstancesClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, sapVirtualInstanceName string, databaseInstanceName string, options *SAPDatabaseInstancesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Workloads/sapVirtualInstances/{sapVirtualInstanceName}/databaseInstances/{databaseInstanceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if sapVirtualInstanceName == "" {
		return nil, errors.New("parameter sapVirtualInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sapVirtualInstanceName}", url.PathEscape(sapVirtualInstanceName))
	if databaseInstanceName == "" {
		return nil, errors.New("parameter databaseInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseInstanceName}", url.PathEscape(databaseInstanceName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the SAP Database Instance resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-04-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - sapVirtualInstanceName - The name of the Virtual Instances for SAP solutions resource
//   - databaseInstanceName - Database resource name string modeled as parameter for auto generation to work correctly.
//   - options - SAPDatabaseInstancesClientGetOptions contains the optional parameters for the SAPDatabaseInstancesClient.Get
//     method.
func (client *SAPDatabaseInstancesClient) Get(ctx context.Context, resourceGroupName string, sapVirtualInstanceName string, databaseInstanceName string, options *SAPDatabaseInstancesClientGetOptions) (SAPDatabaseInstancesClientGetResponse, error) {
	var err error
	const operationName = "SAPDatabaseInstancesClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, sapVirtualInstanceName, databaseInstanceName, options)
	if err != nil {
		return SAPDatabaseInstancesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return SAPDatabaseInstancesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return SAPDatabaseInstancesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *SAPDatabaseInstancesClient) getCreateRequest(ctx context.Context, resourceGroupName string, sapVirtualInstanceName string, databaseInstanceName string, options *SAPDatabaseInstancesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Workloads/sapVirtualInstances/{sapVirtualInstanceName}/databaseInstances/{databaseInstanceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if sapVirtualInstanceName == "" {
		return nil, errors.New("parameter sapVirtualInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sapVirtualInstanceName}", url.PathEscape(sapVirtualInstanceName))
	if databaseInstanceName == "" {
		return nil, errors.New("parameter databaseInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseInstanceName}", url.PathEscape(databaseInstanceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *SAPDatabaseInstancesClient) getHandleResponse(resp *http.Response) (SAPDatabaseInstancesClientGetResponse, error) {
	result := SAPDatabaseInstancesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SAPDatabaseInstance); err != nil {
		return SAPDatabaseInstancesClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Lists the Database resources associated with a Virtual Instance for SAP solutions resource.
//
// Generated from API version 2023-04-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - sapVirtualInstanceName - The name of the Virtual Instances for SAP solutions resource
//   - options - SAPDatabaseInstancesClientListOptions contains the optional parameters for the SAPDatabaseInstancesClient.NewListPager
//     method.
func (client *SAPDatabaseInstancesClient) NewListPager(resourceGroupName string, sapVirtualInstanceName string, options *SAPDatabaseInstancesClientListOptions) *runtime.Pager[SAPDatabaseInstancesClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[SAPDatabaseInstancesClientListResponse]{
		More: func(page SAPDatabaseInstancesClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *SAPDatabaseInstancesClientListResponse) (SAPDatabaseInstancesClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "SAPDatabaseInstancesClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, resourceGroupName, sapVirtualInstanceName, options)
			}, nil)
			if err != nil {
				return SAPDatabaseInstancesClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *SAPDatabaseInstancesClient) listCreateRequest(ctx context.Context, resourceGroupName string, sapVirtualInstanceName string, options *SAPDatabaseInstancesClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Workloads/sapVirtualInstances/{sapVirtualInstanceName}/databaseInstances"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if sapVirtualInstanceName == "" {
		return nil, errors.New("parameter sapVirtualInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sapVirtualInstanceName}", url.PathEscape(sapVirtualInstanceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *SAPDatabaseInstancesClient) listHandleResponse(resp *http.Response) (SAPDatabaseInstancesClientListResponse, error) {
	result := SAPDatabaseInstancesClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.SAPDatabaseInstanceList); err != nil {
		return SAPDatabaseInstancesClientListResponse{}, err
	}
	return result, nil
}

// BeginStartInstance - Starts the database instance of the SAP system.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-04-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - sapVirtualInstanceName - The name of the Virtual Instances for SAP solutions resource
//   - databaseInstanceName - Database resource name string modeled as parameter for auto generation to work correctly.
//   - options - SAPDatabaseInstancesClientBeginStartInstanceOptions contains the optional parameters for the SAPDatabaseInstancesClient.BeginStartInstance
//     method.
func (client *SAPDatabaseInstancesClient) BeginStartInstance(ctx context.Context, resourceGroupName string, sapVirtualInstanceName string, databaseInstanceName string, options *SAPDatabaseInstancesClientBeginStartInstanceOptions) (*runtime.Poller[SAPDatabaseInstancesClientStartInstanceResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.startInstance(ctx, resourceGroupName, sapVirtualInstanceName, databaseInstanceName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[SAPDatabaseInstancesClientStartInstanceResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[SAPDatabaseInstancesClientStartInstanceResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// StartInstance - Starts the database instance of the SAP system.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-04-01
func (client *SAPDatabaseInstancesClient) startInstance(ctx context.Context, resourceGroupName string, sapVirtualInstanceName string, databaseInstanceName string, options *SAPDatabaseInstancesClientBeginStartInstanceOptions) (*http.Response, error) {
	var err error
	const operationName = "SAPDatabaseInstancesClient.BeginStartInstance"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.startInstanceCreateRequest(ctx, resourceGroupName, sapVirtualInstanceName, databaseInstanceName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// startInstanceCreateRequest creates the StartInstance request.
func (client *SAPDatabaseInstancesClient) startInstanceCreateRequest(ctx context.Context, resourceGroupName string, sapVirtualInstanceName string, databaseInstanceName string, options *SAPDatabaseInstancesClientBeginStartInstanceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Workloads/sapVirtualInstances/{sapVirtualInstanceName}/databaseInstances/{databaseInstanceName}/start"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if sapVirtualInstanceName == "" {
		return nil, errors.New("parameter sapVirtualInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sapVirtualInstanceName}", url.PathEscape(sapVirtualInstanceName))
	if databaseInstanceName == "" {
		return nil, errors.New("parameter databaseInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseInstanceName}", url.PathEscape(databaseInstanceName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// BeginStopInstance - Stops the database instance of the SAP system.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-04-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - sapVirtualInstanceName - The name of the Virtual Instances for SAP solutions resource
//   - databaseInstanceName - Database resource name string modeled as parameter for auto generation to work correctly.
//   - options - SAPDatabaseInstancesClientBeginStopInstanceOptions contains the optional parameters for the SAPDatabaseInstancesClient.BeginStopInstance
//     method.
func (client *SAPDatabaseInstancesClient) BeginStopInstance(ctx context.Context, resourceGroupName string, sapVirtualInstanceName string, databaseInstanceName string, options *SAPDatabaseInstancesClientBeginStopInstanceOptions) (*runtime.Poller[SAPDatabaseInstancesClientStopInstanceResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.stopInstance(ctx, resourceGroupName, sapVirtualInstanceName, databaseInstanceName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[SAPDatabaseInstancesClientStopInstanceResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[SAPDatabaseInstancesClientStopInstanceResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// StopInstance - Stops the database instance of the SAP system.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-04-01
func (client *SAPDatabaseInstancesClient) stopInstance(ctx context.Context, resourceGroupName string, sapVirtualInstanceName string, databaseInstanceName string, options *SAPDatabaseInstancesClientBeginStopInstanceOptions) (*http.Response, error) {
	var err error
	const operationName = "SAPDatabaseInstancesClient.BeginStopInstance"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.stopInstanceCreateRequest(ctx, resourceGroupName, sapVirtualInstanceName, databaseInstanceName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// stopInstanceCreateRequest creates the StopInstance request.
func (client *SAPDatabaseInstancesClient) stopInstanceCreateRequest(ctx context.Context, resourceGroupName string, sapVirtualInstanceName string, databaseInstanceName string, options *SAPDatabaseInstancesClientBeginStopInstanceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Workloads/sapVirtualInstances/{sapVirtualInstanceName}/databaseInstances/{databaseInstanceName}/stop"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if sapVirtualInstanceName == "" {
		return nil, errors.New("parameter sapVirtualInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sapVirtualInstanceName}", url.PathEscape(sapVirtualInstanceName))
	if databaseInstanceName == "" {
		return nil, errors.New("parameter databaseInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseInstanceName}", url.PathEscape(databaseInstanceName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if options != nil && options.Body != nil {
		if err := runtime.MarshalAsJSON(req, *options.Body); err != nil {
			return nil, err
		}
		return req, nil
	}
	return req, nil
}

// BeginUpdate - Updates the Database resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-04-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - sapVirtualInstanceName - The name of the Virtual Instances for SAP solutions resource
//   - databaseInstanceName - Database resource name string modeled as parameter for auto generation to work correctly.
//   - body - Database resource update request body.
//   - options - SAPDatabaseInstancesClientBeginUpdateOptions contains the optional parameters for the SAPDatabaseInstancesClient.BeginUpdate
//     method.
func (client *SAPDatabaseInstancesClient) BeginUpdate(ctx context.Context, resourceGroupName string, sapVirtualInstanceName string, databaseInstanceName string, body UpdateSAPDatabaseInstanceRequest, options *SAPDatabaseInstancesClientBeginUpdateOptions) (*runtime.Poller[SAPDatabaseInstancesClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, sapVirtualInstanceName, databaseInstanceName, body, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[SAPDatabaseInstancesClientUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[SAPDatabaseInstancesClientUpdateResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// Update - Updates the Database resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-04-01
func (client *SAPDatabaseInstancesClient) update(ctx context.Context, resourceGroupName string, sapVirtualInstanceName string, databaseInstanceName string, body UpdateSAPDatabaseInstanceRequest, options *SAPDatabaseInstancesClientBeginUpdateOptions) (*http.Response, error) {
	var err error
	const operationName = "SAPDatabaseInstancesClient.BeginUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, sapVirtualInstanceName, databaseInstanceName, body, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// updateCreateRequest creates the Update request.
func (client *SAPDatabaseInstancesClient) updateCreateRequest(ctx context.Context, resourceGroupName string, sapVirtualInstanceName string, databaseInstanceName string, body UpdateSAPDatabaseInstanceRequest, options *SAPDatabaseInstancesClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Workloads/sapVirtualInstances/{sapVirtualInstanceName}/databaseInstances/{databaseInstanceName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if sapVirtualInstanceName == "" {
		return nil, errors.New("parameter sapVirtualInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{sapVirtualInstanceName}", url.PathEscape(sapVirtualInstanceName))
	if databaseInstanceName == "" {
		return nil, errors.New("parameter databaseInstanceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{databaseInstanceName}", url.PathEscape(databaseInstanceName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, body); err != nil {
		return nil, err
	}
	return req, nil
}