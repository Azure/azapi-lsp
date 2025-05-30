// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armapimanagement

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

// AuthorizationClient contains the methods for the Authorization group.
// Don't use this type directly, use NewAuthorizationClient() instead.
type AuthorizationClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewAuthorizationClient creates a new instance of AuthorizationClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewAuthorizationClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*AuthorizationClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &AuthorizationClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// ConfirmConsentCode - Confirm valid consent code to suppress Authorizations anti-phishing page.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serviceName - The name of the API Management service.
//   - authorizationProviderID - Identifier of the authorization provider.
//   - authorizationID - Identifier of the authorization.
//   - parameters - Create parameters.
//   - options - AuthorizationClientConfirmConsentCodeOptions contains the optional parameters for the AuthorizationClient.ConfirmConsentCode
//     method.
func (client *AuthorizationClient) ConfirmConsentCode(ctx context.Context, resourceGroupName string, serviceName string, authorizationProviderID string, authorizationID string, parameters AuthorizationConfirmConsentCodeRequestContract, options *AuthorizationClientConfirmConsentCodeOptions) (AuthorizationClientConfirmConsentCodeResponse, error) {
	var err error
	const operationName = "AuthorizationClient.ConfirmConsentCode"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.confirmConsentCodeCreateRequest(ctx, resourceGroupName, serviceName, authorizationProviderID, authorizationID, parameters, options)
	if err != nil {
		return AuthorizationClientConfirmConsentCodeResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AuthorizationClientConfirmConsentCodeResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AuthorizationClientConfirmConsentCodeResponse{}, err
	}
	resp, err := client.confirmConsentCodeHandleResponse(httpResp)
	return resp, err
}

// confirmConsentCodeCreateRequest creates the ConfirmConsentCode request.
func (client *AuthorizationClient) confirmConsentCodeCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, authorizationProviderID string, authorizationID string, parameters AuthorizationConfirmConsentCodeRequestContract, _ *AuthorizationClientConfirmConsentCodeOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/authorizationProviders/{authorizationProviderId}/authorizations/{authorizationId}/confirmConsentCode"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if authorizationProviderID == "" {
		return nil, errors.New("parameter authorizationProviderID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{authorizationProviderId}", url.PathEscape(authorizationProviderID))
	if authorizationID == "" {
		return nil, errors.New("parameter authorizationID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{authorizationId}", url.PathEscape(authorizationID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// confirmConsentCodeHandleResponse handles the ConfirmConsentCode response.
func (client *AuthorizationClient) confirmConsentCodeHandleResponse(resp *http.Response) (AuthorizationClientConfirmConsentCodeResponse, error) {
	result := AuthorizationClientConfirmConsentCodeResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	return result, nil
}

// CreateOrUpdate - Creates or updates authorization.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serviceName - The name of the API Management service.
//   - authorizationProviderID - Identifier of the authorization provider.
//   - authorizationID - Identifier of the authorization.
//   - parameters - Create parameters.
//   - options - AuthorizationClientCreateOrUpdateOptions contains the optional parameters for the AuthorizationClient.CreateOrUpdate
//     method.
func (client *AuthorizationClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, serviceName string, authorizationProviderID string, authorizationID string, parameters AuthorizationContract, options *AuthorizationClientCreateOrUpdateOptions) (AuthorizationClientCreateOrUpdateResponse, error) {
	var err error
	const operationName = "AuthorizationClient.CreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, serviceName, authorizationProviderID, authorizationID, parameters, options)
	if err != nil {
		return AuthorizationClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AuthorizationClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return AuthorizationClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *AuthorizationClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, authorizationProviderID string, authorizationID string, parameters AuthorizationContract, options *AuthorizationClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/authorizationProviders/{authorizationProviderId}/authorizations/{authorizationId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if authorizationProviderID == "" {
		return nil, errors.New("parameter authorizationProviderID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{authorizationProviderId}", url.PathEscape(authorizationProviderID))
	if authorizationID == "" {
		return nil, errors.New("parameter authorizationID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{authorizationId}", url.PathEscape(authorizationID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if options != nil && options.IfMatch != nil {
		req.Raw().Header["If-Match"] = []string{*options.IfMatch}
	}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *AuthorizationClient) createOrUpdateHandleResponse(resp *http.Response) (AuthorizationClientCreateOrUpdateResponse, error) {
	result := AuthorizationClientCreateOrUpdateResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.AuthorizationContract); err != nil {
		return AuthorizationClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes specific Authorization from the Authorization provider.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serviceName - The name of the API Management service.
//   - authorizationProviderID - Identifier of the authorization provider.
//   - authorizationID - Identifier of the authorization.
//   - ifMatch - ETag of the Entity. ETag should match the current entity state from the header response of the GET request or
//     it should be * for unconditional update.
//   - options - AuthorizationClientDeleteOptions contains the optional parameters for the AuthorizationClient.Delete method.
func (client *AuthorizationClient) Delete(ctx context.Context, resourceGroupName string, serviceName string, authorizationProviderID string, authorizationID string, ifMatch string, options *AuthorizationClientDeleteOptions) (AuthorizationClientDeleteResponse, error) {
	var err error
	const operationName = "AuthorizationClient.Delete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, serviceName, authorizationProviderID, authorizationID, ifMatch, options)
	if err != nil {
		return AuthorizationClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AuthorizationClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return AuthorizationClientDeleteResponse{}, err
	}
	return AuthorizationClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *AuthorizationClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, authorizationProviderID string, authorizationID string, ifMatch string, _ *AuthorizationClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/authorizationProviders/{authorizationProviderId}/authorizations/{authorizationId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if authorizationProviderID == "" {
		return nil, errors.New("parameter authorizationProviderID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{authorizationProviderId}", url.PathEscape(authorizationProviderID))
	if authorizationID == "" {
		return nil, errors.New("parameter authorizationID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{authorizationId}", url.PathEscape(authorizationID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	req.Raw().Header["If-Match"] = []string{ifMatch}
	return req, nil
}

// Get - Gets the details of the authorization specified by its identifier.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serviceName - The name of the API Management service.
//   - authorizationProviderID - Identifier of the authorization provider.
//   - authorizationID - Identifier of the authorization.
//   - options - AuthorizationClientGetOptions contains the optional parameters for the AuthorizationClient.Get method.
func (client *AuthorizationClient) Get(ctx context.Context, resourceGroupName string, serviceName string, authorizationProviderID string, authorizationID string, options *AuthorizationClientGetOptions) (AuthorizationClientGetResponse, error) {
	var err error
	const operationName = "AuthorizationClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, serviceName, authorizationProviderID, authorizationID, options)
	if err != nil {
		return AuthorizationClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AuthorizationClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AuthorizationClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *AuthorizationClient) getCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, authorizationProviderID string, authorizationID string, _ *AuthorizationClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/authorizationProviders/{authorizationProviderId}/authorizations/{authorizationId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if authorizationProviderID == "" {
		return nil, errors.New("parameter authorizationProviderID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{authorizationProviderId}", url.PathEscape(authorizationProviderID))
	if authorizationID == "" {
		return nil, errors.New("parameter authorizationID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{authorizationId}", url.PathEscape(authorizationID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *AuthorizationClient) getHandleResponse(resp *http.Response) (AuthorizationClientGetResponse, error) {
	result := AuthorizationClientGetResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.AuthorizationContract); err != nil {
		return AuthorizationClientGetResponse{}, err
	}
	return result, nil
}

// NewListByAuthorizationProviderPager - Lists a collection of authorization providers defined within a authorization provider.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serviceName - The name of the API Management service.
//   - authorizationProviderID - Identifier of the authorization provider.
//   - options - AuthorizationClientListByAuthorizationProviderOptions contains the optional parameters for the AuthorizationClient.NewListByAuthorizationProviderPager
//     method.
func (client *AuthorizationClient) NewListByAuthorizationProviderPager(resourceGroupName string, serviceName string, authorizationProviderID string, options *AuthorizationClientListByAuthorizationProviderOptions) *runtime.Pager[AuthorizationClientListByAuthorizationProviderResponse] {
	return runtime.NewPager(runtime.PagingHandler[AuthorizationClientListByAuthorizationProviderResponse]{
		More: func(page AuthorizationClientListByAuthorizationProviderResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AuthorizationClientListByAuthorizationProviderResponse) (AuthorizationClientListByAuthorizationProviderResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "AuthorizationClient.NewListByAuthorizationProviderPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByAuthorizationProviderCreateRequest(ctx, resourceGroupName, serviceName, authorizationProviderID, options)
			}, nil)
			if err != nil {
				return AuthorizationClientListByAuthorizationProviderResponse{}, err
			}
			return client.listByAuthorizationProviderHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByAuthorizationProviderCreateRequest creates the ListByAuthorizationProvider request.
func (client *AuthorizationClient) listByAuthorizationProviderCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, authorizationProviderID string, options *AuthorizationClientListByAuthorizationProviderOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/authorizationProviders/{authorizationProviderId}/authorizations"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if authorizationProviderID == "" {
		return nil, errors.New("parameter authorizationProviderID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{authorizationProviderId}", url.PathEscape(authorizationProviderID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Skip != nil {
		reqQP.Set("$skip", strconv.FormatInt(int64(*options.Skip), 10))
	}
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByAuthorizationProviderHandleResponse handles the ListByAuthorizationProvider response.
func (client *AuthorizationClient) listByAuthorizationProviderHandleResponse(resp *http.Response) (AuthorizationClientListByAuthorizationProviderResponse, error) {
	result := AuthorizationClientListByAuthorizationProviderResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AuthorizationCollection); err != nil {
		return AuthorizationClientListByAuthorizationProviderResponse{}, err
	}
	return result, nil
}
