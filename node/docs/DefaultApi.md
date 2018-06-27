# NetlifysApiDefinition.DefaultApi

All URIs are relative to *https://api.netlify.com/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**addMemberToAccount**](DefaultApi.md#addMemberToAccount) | **POST** /{account_slug}/members | 
[**cancelAccount**](DefaultApi.md#cancelAccount) | **DELETE** /accounts/{account_id} | 
[**configureDNSForSite**](DefaultApi.md#configureDNSForSite) | **PUT** /sites/{site_id}/dns | 
[**createAccount**](DefaultApi.md#createAccount) | **POST** /accounts | 
[**createDeployKey**](DefaultApi.md#createDeployKey) | **POST** /deploy_keys | 
[**createHookBySiteId**](DefaultApi.md#createHookBySiteId) | **POST** /hooks | 
[**createSite**](DefaultApi.md#createSite) | **POST** /sites | 
[**createSiteAsset**](DefaultApi.md#createSiteAsset) | **POST** /sites/{site_id}/assets | 
[**createSiteBuildHook**](DefaultApi.md#createSiteBuildHook) | **POST** /sites/{site_id}/build_hooks | 
[**createSiteDeploy**](DefaultApi.md#createSiteDeploy) | **POST** /sites/{site_id}/deploys | 
[**createSiteSnippet**](DefaultApi.md#createSiteSnippet) | **POST** /sites/{site_id}/snippets | 
[**createTicket**](DefaultApi.md#createTicket) | **POST** /oauth/tickets | 
[**deleteDeployKey**](DefaultApi.md#deleteDeployKey) | **DELETE** /deploy_keys/{key_id} | 
[**deleteHookBySiteId**](DefaultApi.md#deleteHookBySiteId) | **DELETE** /hooks/{hook_id} | 
[**deleteSite**](DefaultApi.md#deleteSite) | **DELETE** /sites/{site_id} | 
[**deleteSiteAsset**](DefaultApi.md#deleteSiteAsset) | **DELETE** /sites/{site_id}/assets/{asset_id} | 
[**deleteSiteBuildHook**](DefaultApi.md#deleteSiteBuildHook) | **DELETE** /sites/{site_id}/build_hooks/{id} | 
[**deleteSiteSnippet**](DefaultApi.md#deleteSiteSnippet) | **DELETE** /sites/{site_id}/snippets/{snippet_id} | 
[**deleteSubmission**](DefaultApi.md#deleteSubmission) | **DELETE** /submissions/{submission_id} | 
[**enableHook**](DefaultApi.md#enableHook) | **POST** /hooks/{hook_id}/enable | 
[**exchangeTicket**](DefaultApi.md#exchangeTicket) | **POST** /oauth/tickets/{ticket_id}/exchange | 
[**getDNSForSite**](DefaultApi.md#getDNSForSite) | **GET** /sites/{site_id}/dns | 
[**getDeploy**](DefaultApi.md#getDeploy) | **GET** /deploys/{deploy_id} | 
[**getDeployKey**](DefaultApi.md#getDeployKey) | **GET** /deploy_keys/{key_id} | 
[**getHook**](DefaultApi.md#getHook) | **GET** /hooks/{hook_id} | 
[**getSite**](DefaultApi.md#getSite) | **GET** /sites/{site_id} | 
[**getSiteAssetInfo**](DefaultApi.md#getSiteAssetInfo) | **GET** /sites/{site_id}/assets/{asset_id} | 
[**getSiteAssetPublicSignature**](DefaultApi.md#getSiteAssetPublicSignature) | **GET** /sites/{site_id}/assets/{asset_id}/public_signature | 
[**getSiteBuild**](DefaultApi.md#getSiteBuild) | **GET** /builds/{build_id} | 
[**getSiteBuildHook**](DefaultApi.md#getSiteBuildHook) | **GET** /sites/{site_id}/build_hooks/{id} | 
[**getSiteDeploy**](DefaultApi.md#getSiteDeploy) | **GET** /sites/{site_id}/deploys/{deploy_id} | 
[**getSiteFileByPathName**](DefaultApi.md#getSiteFileByPathName) | **GET** /sites/{site_id}/files/{file_path} | 
[**getSiteMetadata**](DefaultApi.md#getSiteMetadata) | **GET** /sites/{site_id}/metadata | 
[**getSiteSnippet**](DefaultApi.md#getSiteSnippet) | **GET** /sites/{site_id}/snippets/{snippet_id} | 
[**listAccountAuditEvents**](DefaultApi.md#listAccountAuditEvents) | **GET** /accounts/{account_id}/audit | 
[**listAccountTypesForUser**](DefaultApi.md#listAccountTypesForUser) | **GET** /accounts/types | 
[**listAccountsForUser**](DefaultApi.md#listAccountsForUser) | **GET** /accounts | 
[**listDeployKeys**](DefaultApi.md#listDeployKeys) | **GET** /deploy_keys | 
[**listFormSubmission**](DefaultApi.md#listFormSubmission) | **GET** /submissions/{submission_id} | 
[**listFormSubmissions**](DefaultApi.md#listFormSubmissions) | **GET** /forms/{form_id}/submissions | 
[**listForms**](DefaultApi.md#listForms) | **GET** /forms | 
[**listHookTypes**](DefaultApi.md#listHookTypes) | **GET** /hooks/types | 
[**listHooksBySiteId**](DefaultApi.md#listHooksBySiteId) | **GET** /hooks | 
[**listMembersForAccount**](DefaultApi.md#listMembersForAccount) | **GET** /{account_slug}/members | 
[**listPaymentMethodsForUser**](DefaultApi.md#listPaymentMethodsForUser) | **GET** /billing/payment_methods | 
[**listSiteAssets**](DefaultApi.md#listSiteAssets) | **GET** /sites/{site_id}/assets | 
[**listSiteBuildHooks**](DefaultApi.md#listSiteBuildHooks) | **GET** /sites/{site_id}/build_hooks | 
[**listSiteBuilds**](DefaultApi.md#listSiteBuilds) | **GET** /sites/{site_id}/builds | 
[**listSiteDeployedBranches**](DefaultApi.md#listSiteDeployedBranches) | **GET** /sites/{site_id}/deployed-branches | 
[**listSiteDeploys**](DefaultApi.md#listSiteDeploys) | **GET** /sites/{site_id}/deploys | 
[**listSiteFiles**](DefaultApi.md#listSiteFiles) | **GET** /sites/{site_id}/files | 
[**listSiteForms**](DefaultApi.md#listSiteForms) | **GET** /sites/{site_id}/forms | 
[**listSiteSnippets**](DefaultApi.md#listSiteSnippets) | **GET** /sites/{site_id}/snippets | 
[**listSiteSubmissions**](DefaultApi.md#listSiteSubmissions) | **GET** /sites/{site_id}/submissions | 
[**listSites**](DefaultApi.md#listSites) | **GET** /sites | 
[**listSitesForAccount**](DefaultApi.md#listSitesForAccount) | **GET** /{account_slug}/sites | 
[**lockDeploy**](DefaultApi.md#lockDeploy) | **POST** /deploys/{deploy_id}/lock | 
[**notifyBuildStart**](DefaultApi.md#notifyBuildStart) | **POST** /builds/{build_id}/start | 
[**provisionSiteTLSCertificate**](DefaultApi.md#provisionSiteTLSCertificate) | **POST** /sites/{site_id}/ssl | 
[**restoreSiteDeploy**](DefaultApi.md#restoreSiteDeploy) | **POST** /sites/{site_id}/deploys/{deploy_id}/restore | 
[**showSiteTLSCertificate**](DefaultApi.md#showSiteTLSCertificate) | **GET** /sites/{site_id}/ssl | 
[**showTicket**](DefaultApi.md#showTicket) | **GET** /oauth/tickets/{ticket_id} | 
[**unlockDeploy**](DefaultApi.md#unlockDeploy) | **POST** /deploys/{deploy_id}/unlock | 
[**updateAccount**](DefaultApi.md#updateAccount) | **PUT** /accounts/{account_id} | 
[**updateHook**](DefaultApi.md#updateHook) | **PUT** /hooks/{hook_id} | 
[**updateSite**](DefaultApi.md#updateSite) | **PATCH** /sites/{site_id} | 
[**updateSiteAsset**](DefaultApi.md#updateSiteAsset) | **PUT** /sites/{site_id}/assets/{asset_id} | 
[**updateSiteBuildHook**](DefaultApi.md#updateSiteBuildHook) | **PUT** /sites/{site_id}/build_hooks/{id} | 
[**updateSiteBuildLog**](DefaultApi.md#updateSiteBuildLog) | **POST** /builds/{build_id}/log | 
[**updateSiteDeploy**](DefaultApi.md#updateSiteDeploy) | **PUT** /sites/{site_id}/deploys/{deploy_id} | 
[**updateSiteMetadata**](DefaultApi.md#updateSiteMetadata) | **PUT** /sites/{site_id}/metadata | 
[**updateSiteSnippet**](DefaultApi.md#updateSiteSnippet) | **PUT** /sites/{site_id}/snippets/{snippet_id} | 
[**uploadDeployFile**](DefaultApi.md#uploadDeployFile) | **PUT** /deploys/{deploy_id}/files/{path} | 
[**uploadDeployFunction**](DefaultApi.md#uploadDeployFunction) | **PUT** /deploys/{deploy_id}/functions/{name} | 


<a name="addMemberToAccount"></a>
# **addMemberToAccount**
> [Member] addMemberToAccount(accountSlug, email, opts)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var accountSlug = "accountSlug_example"; // String | 

var email = "email_example"; // String | 

var opts = { 
  'role': "role_example" // String | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.addMemberToAccount(accountSlug, email, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **accountSlug** | **String**|  | 
 **email** | **String**|  | 
 **role** | **String**|  | [optional] 

### Return type

[**[Member]**](Member.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="cancelAccount"></a>
# **cancelAccount**
> cancelAccount(accountId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var accountId = "accountId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.cancelAccount(accountId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **accountId** | **String**|  | 

### Return type

null (empty response body)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="configureDNSForSite"></a>
# **configureDNSForSite**
> [DnsZone] configureDNSForSite(siteId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.configureDNSForSite(siteId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 

### Return type

[**[DnsZone]**](DnsZone.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="createAccount"></a>
# **createAccount**
> AccountMembership createAccount(accountSetup)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var accountSetup = new NetlifysApiDefinition.AccountSetup(); // AccountSetup | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.createAccount(accountSetup, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **accountSetup** | [**AccountSetup**](AccountSetup.md)|  | 

### Return type

[**AccountMembership**](AccountMembership.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="createDeployKey"></a>
# **createDeployKey**
> DeployKey createDeployKey()



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.createDeployKey(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**DeployKey**](DeployKey.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="createHookBySiteId"></a>
# **createHookBySiteId**
> Hook createHookBySiteId(siteId, hook)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var hook = new NetlifysApiDefinition.Hook(); // Hook | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.createHookBySiteId(siteId, hook, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **hook** | [**Hook**](Hook.md)|  | 

### Return type

[**Hook**](Hook.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="createSite"></a>
# **createSite**
> Site createSite(site, opts)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var site = new NetlifysApiDefinition.SiteSetup(); // SiteSetup | 

var opts = { 
  'configureDns': true // Boolean | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.createSite(site, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **site** | [**SiteSetup**](SiteSetup.md)|  | 
 **configureDns** | **Boolean**|  | [optional] 

### Return type

[**Site**](Site.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="createSiteAsset"></a>
# **createSiteAsset**
> AssetSignature createSiteAsset(siteId, name, size, contentType, opts)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var name = "name_example"; // String | 

var size = 789; // Number | 

var contentType = "contentType_example"; // String | 

var opts = { 
  'visibility': "visibility_example" // String | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.createSiteAsset(siteId, name, size, contentType, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **name** | **String**|  | 
 **size** | **Number**|  | 
 **contentType** | **String**|  | 
 **visibility** | **String**|  | [optional] 

### Return type

[**AssetSignature**](AssetSignature.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="createSiteBuildHook"></a>
# **createSiteBuildHook**
> BuildHook createSiteBuildHook(siteId, buildHook)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var buildHook = new NetlifysApiDefinition.BuildHook(); // BuildHook | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.createSiteBuildHook(siteId, buildHook, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **buildHook** | [**BuildHook**](BuildHook.md)|  | 

### Return type

[**BuildHook**](BuildHook.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="createSiteDeploy"></a>
# **createSiteDeploy**
> Deploy createSiteDeploy(siteId, deploy, opts)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var deploy = new NetlifysApiDefinition.DeployFiles(); // DeployFiles | 

var opts = { 
  'title': "title_example" // String | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.createSiteDeploy(siteId, deploy, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **deploy** | [**DeployFiles**](DeployFiles.md)|  | 
 **title** | **String**|  | [optional] 

### Return type

[**Deploy**](Deploy.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="createSiteSnippet"></a>
# **createSiteSnippet**
> Snippet createSiteSnippet(siteId, snippet)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var snippet = new NetlifysApiDefinition.Snippet(); // Snippet | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.createSiteSnippet(siteId, snippet, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **snippet** | [**Snippet**](Snippet.md)|  | 

### Return type

[**Snippet**](Snippet.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="createTicket"></a>
# **createTicket**
> Ticket createTicket(clientId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var clientId = "clientId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.createTicket(clientId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **clientId** | **String**|  | 

### Return type

[**Ticket**](Ticket.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteDeployKey"></a>
# **deleteDeployKey**
> deleteDeployKey(keyId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var keyId = "keyId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteDeployKey(keyId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **keyId** | **String**|  | 

### Return type

null (empty response body)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteHookBySiteId"></a>
# **deleteHookBySiteId**
> deleteHookBySiteId(hookId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var hookId = "hookId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteHookBySiteId(hookId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **hookId** | **String**|  | 

### Return type

null (empty response body)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteSite"></a>
# **deleteSite**
> deleteSite(siteId, )



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteSite(siteId, , callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 

### Return type

null (empty response body)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteSiteAsset"></a>
# **deleteSiteAsset**
> deleteSiteAsset(siteId, assetId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var assetId = "assetId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteSiteAsset(siteId, assetId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **assetId** | **String**|  | 

### Return type

null (empty response body)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteSiteBuildHook"></a>
# **deleteSiteBuildHook**
> deleteSiteBuildHook(siteId, id)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var id = "id_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteSiteBuildHook(siteId, id, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **id** | **String**|  | 

### Return type

null (empty response body)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteSiteSnippet"></a>
# **deleteSiteSnippet**
> deleteSiteSnippet(siteId, snippetId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var snippetId = "snippetId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteSiteSnippet(siteId, snippetId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **snippetId** | **String**|  | 

### Return type

null (empty response body)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="deleteSubmission"></a>
# **deleteSubmission**
> deleteSubmission(submissionId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var submissionId = "submissionId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.deleteSubmission(submissionId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **submissionId** | **String**|  | 

### Return type

null (empty response body)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="enableHook"></a>
# **enableHook**
> Hook enableHook(hookId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var hookId = "hookId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.enableHook(hookId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **hookId** | **String**|  | 

### Return type

[**Hook**](Hook.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="exchangeTicket"></a>
# **exchangeTicket**
> AccessToken exchangeTicket(ticketId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var ticketId = "ticketId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.exchangeTicket(ticketId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ticketId** | **String**|  | 

### Return type

[**AccessToken**](AccessToken.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getDNSForSite"></a>
# **getDNSForSite**
> [DnsZone] getDNSForSite(siteId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getDNSForSite(siteId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 

### Return type

[**[DnsZone]**](DnsZone.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getDeploy"></a>
# **getDeploy**
> Deploy getDeploy(deployId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var deployId = "deployId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getDeploy(deployId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deployId** | **String**|  | 

### Return type

[**Deploy**](Deploy.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getDeployKey"></a>
# **getDeployKey**
> DeployKey getDeployKey(keyId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var keyId = "keyId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getDeployKey(keyId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **keyId** | **String**|  | 

### Return type

[**DeployKey**](DeployKey.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getHook"></a>
# **getHook**
> Hook getHook(hookId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var hookId = "hookId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getHook(hookId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **hookId** | **String**|  | 

### Return type

[**Hook**](Hook.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getSite"></a>
# **getSite**
> Site getSite(siteId, )



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getSite(siteId, , callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 

### Return type

[**Site**](Site.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getSiteAssetInfo"></a>
# **getSiteAssetInfo**
> Asset getSiteAssetInfo(siteId, assetId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var assetId = "assetId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getSiteAssetInfo(siteId, assetId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **assetId** | **String**|  | 

### Return type

[**Asset**](Asset.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getSiteAssetPublicSignature"></a>
# **getSiteAssetPublicSignature**
> AssetPublicSignature getSiteAssetPublicSignature(siteId, assetId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var assetId = "assetId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getSiteAssetPublicSignature(siteId, assetId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **assetId** | **String**|  | 

### Return type

[**AssetPublicSignature**](AssetPublicSignature.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getSiteBuild"></a>
# **getSiteBuild**
> Build getSiteBuild(buildId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var buildId = "buildId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getSiteBuild(buildId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **buildId** | **String**|  | 

### Return type

[**Build**](Build.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getSiteBuildHook"></a>
# **getSiteBuildHook**
> BuildHook getSiteBuildHook(siteId, id)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var id = "id_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getSiteBuildHook(siteId, id, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **id** | **String**|  | 

### Return type

[**BuildHook**](BuildHook.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getSiteDeploy"></a>
# **getSiteDeploy**
> Deploy getSiteDeploy(siteId, deployId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var deployId = "deployId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getSiteDeploy(siteId, deployId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **deployId** | **String**|  | 

### Return type

[**Deploy**](Deploy.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getSiteFileByPathName"></a>
# **getSiteFileByPathName**
> File getSiteFileByPathName(siteId, filePath)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var filePath = "filePath_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getSiteFileByPathName(siteId, filePath, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **filePath** | **String**|  | 

### Return type

**File**

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getSiteMetadata"></a>
# **getSiteMetadata**
> Metadata getSiteMetadata(siteId, )



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getSiteMetadata(siteId, , callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 

### Return type

[**Metadata**](Metadata.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="getSiteSnippet"></a>
# **getSiteSnippet**
> Snippet getSiteSnippet(siteId, snippetId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var snippetId = "snippetId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getSiteSnippet(siteId, snippetId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **snippetId** | **String**|  | 

### Return type

[**Snippet**](Snippet.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listAccountAuditEvents"></a>
# **listAccountAuditEvents**
> [AuditLog] listAccountAuditEvents(accountId, opts)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var accountId = "accountId_example"; // String | 

var opts = { 
  'query': "query_example", // String | 
  'logType': "logType_example" // String | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listAccountAuditEvents(accountId, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **accountId** | **String**|  | 
 **query** | **String**|  | [optional] 
 **logType** | **String**|  | [optional] 

### Return type

[**[AuditLog]**](AuditLog.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listAccountTypesForUser"></a>
# **listAccountTypesForUser**
> [AccountType] listAccountTypesForUser()



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listAccountTypesForUser(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**[AccountType]**](AccountType.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listAccountsForUser"></a>
# **listAccountsForUser**
> [AccountMembership] listAccountsForUser()



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listAccountsForUser(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**[AccountMembership]**](AccountMembership.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listDeployKeys"></a>
# **listDeployKeys**
> [DeployKey] listDeployKeys()



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listDeployKeys(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**[DeployKey]**](DeployKey.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listFormSubmission"></a>
# **listFormSubmission**
> [Submission] listFormSubmission(submissionId, opts)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var submissionId = "submissionId_example"; // String | 

var opts = { 
  'query': "query_example" // String | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listFormSubmission(submissionId, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **submissionId** | **String**|  | 
 **query** | **String**|  | [optional] 

### Return type

[**[Submission]**](Submission.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listFormSubmissions"></a>
# **listFormSubmissions**
> [Submission] listFormSubmissions(formId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var formId = "formId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listFormSubmissions(formId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **formId** | **String**|  | 

### Return type

[**[Submission]**](Submission.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listForms"></a>
# **listForms**
> [Form] listForms(opts)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var opts = { 
  'siteId': "siteId_example" // String | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listForms(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | [optional] 

### Return type

[**[Form]**](Form.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listHookTypes"></a>
# **listHookTypes**
> [HookType] listHookTypes()



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listHookTypes(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**[HookType]**](HookType.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listHooksBySiteId"></a>
# **listHooksBySiteId**
> [Hook] listHooksBySiteId(siteId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listHooksBySiteId(siteId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 

### Return type

[**[Hook]**](Hook.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listMembersForAccount"></a>
# **listMembersForAccount**
> [Member] listMembersForAccount(accountSlug, )



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var accountSlug = "accountSlug_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listMembersForAccount(accountSlug, , callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **accountSlug** | **String**|  | 

### Return type

[**[Member]**](Member.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listPaymentMethodsForUser"></a>
# **listPaymentMethodsForUser**
> [PaymentMethod] listPaymentMethodsForUser()



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listPaymentMethodsForUser(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**[PaymentMethod]**](PaymentMethod.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listSiteAssets"></a>
# **listSiteAssets**
> [Asset] listSiteAssets(siteId, )



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listSiteAssets(siteId, , callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 

### Return type

[**[Asset]**](Asset.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listSiteBuildHooks"></a>
# **listSiteBuildHooks**
> [BuildHook] listSiteBuildHooks(siteId, )



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listSiteBuildHooks(siteId, , callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 

### Return type

[**[BuildHook]**](BuildHook.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listSiteBuilds"></a>
# **listSiteBuilds**
> [Build] listSiteBuilds(siteId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listSiteBuilds(siteId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 

### Return type

[**[Build]**](Build.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listSiteDeployedBranches"></a>
# **listSiteDeployedBranches**
> [DeployedBranch] listSiteDeployedBranches(siteId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listSiteDeployedBranches(siteId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 

### Return type

[**[DeployedBranch]**](DeployedBranch.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listSiteDeploys"></a>
# **listSiteDeploys**
> [Deploy] listSiteDeploys(siteId, )



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listSiteDeploys(siteId, , callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 

### Return type

[**[Deploy]**](Deploy.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listSiteFiles"></a>
# **listSiteFiles**
> [File] listSiteFiles(siteId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listSiteFiles(siteId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 

### Return type

**[File]**

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listSiteForms"></a>
# **listSiteForms**
> [Form] listSiteForms(siteId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listSiteForms(siteId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 

### Return type

[**[Form]**](Form.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listSiteSnippets"></a>
# **listSiteSnippets**
> [Snippet] listSiteSnippets(siteId, )



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listSiteSnippets(siteId, , callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 

### Return type

[**[Snippet]**](Snippet.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listSiteSubmissions"></a>
# **listSiteSubmissions**
> [Submission] listSiteSubmissions(siteId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listSiteSubmissions(siteId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 

### Return type

[**[Submission]**](Submission.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listSites"></a>
# **listSites**
> [Site] listSites(opts)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var opts = { 
  'name': "name_example", // String | 
  'filter': "filter_example" // String | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listSites(opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **String**|  | [optional] 
 **filter** | **String**|  | [optional] 

### Return type

[**[Site]**](Site.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="listSitesForAccount"></a>
# **listSitesForAccount**
> [Site] listSitesForAccount(accountSlug, opts)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var accountSlug = "accountSlug_example"; // String | 

var opts = { 
  'name': "name_example" // String | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.listSitesForAccount(accountSlug, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **accountSlug** | **String**|  | 
 **name** | **String**|  | [optional] 

### Return type

[**[Site]**](Site.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="lockDeploy"></a>
# **lockDeploy**
> Deploy lockDeploy(deployId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var deployId = "deployId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.lockDeploy(deployId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deployId** | **String**|  | 

### Return type

[**Deploy**](Deploy.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="notifyBuildStart"></a>
# **notifyBuildStart**
> notifyBuildStart(buildId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var buildId = "buildId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.notifyBuildStart(buildId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **buildId** | **String**|  | 

### Return type

null (empty response body)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="provisionSiteTLSCertificate"></a>
# **provisionSiteTLSCertificate**
> SniCertificate provisionSiteTLSCertificate(siteId, opts)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var opts = { 
  'certificate': "certificate_example", // String | 
  'key': "key_example", // String | 
  'caCertificates': "caCertificates_example" // String | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.provisionSiteTLSCertificate(siteId, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **certificate** | **String**|  | [optional] 
 **key** | **String**|  | [optional] 
 **caCertificates** | **String**|  | [optional] 

### Return type

[**SniCertificate**](SniCertificate.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="restoreSiteDeploy"></a>
# **restoreSiteDeploy**
> Deploy restoreSiteDeploy(siteId, deployId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var deployId = "deployId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.restoreSiteDeploy(siteId, deployId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **deployId** | **String**|  | 

### Return type

[**Deploy**](Deploy.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="showSiteTLSCertificate"></a>
# **showSiteTLSCertificate**
> SniCertificate showSiteTLSCertificate(siteId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.showSiteTLSCertificate(siteId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 

### Return type

[**SniCertificate**](SniCertificate.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="showTicket"></a>
# **showTicket**
> Ticket showTicket(ticketId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var ticketId = "ticketId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.showTicket(ticketId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ticketId** | **String**|  | 

### Return type

[**Ticket**](Ticket.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="unlockDeploy"></a>
# **unlockDeploy**
> Deploy unlockDeploy(deployId)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var deployId = "deployId_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.unlockDeploy(deployId, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deployId** | **String**|  | 

### Return type

[**Deploy**](Deploy.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="updateAccount"></a>
# **updateAccount**
> AccountMembership updateAccount(accountId, opts)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var accountId = "accountId_example"; // String | 

var opts = { 
  'accountUpdateSetup': new NetlifysApiDefinition.AccountUpdateSetup() // AccountUpdateSetup | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.updateAccount(accountId, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **accountId** | **String**|  | 
 **accountUpdateSetup** | [**AccountUpdateSetup**](AccountUpdateSetup.md)|  | [optional] 

### Return type

[**AccountMembership**](AccountMembership.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="updateHook"></a>
# **updateHook**
> Hook updateHook(hookIdhook)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var hookId = "hookId_example"; // String | 

var hook = new NetlifysApiDefinition.Hook(); // Hook | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.updateHook(hookIdhook, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **hookId** | **String**|  | 
 **hook** | [**Hook**](Hook.md)|  | 

### Return type

[**Hook**](Hook.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="updateSite"></a>
# **updateSite**
> Site updateSite(siteId, site)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var site = new NetlifysApiDefinition.SiteSetup(); // SiteSetup | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.updateSite(siteId, site, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **site** | [**SiteSetup**](SiteSetup.md)|  | 

### Return type

[**Site**](Site.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="updateSiteAsset"></a>
# **updateSiteAsset**
> Asset updateSiteAsset(siteId, assetIdstate)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var assetId = "assetId_example"; // String | 

var state = "state_example"; // String | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.updateSiteAsset(siteId, assetIdstate, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **assetId** | **String**|  | 
 **state** | **String**|  | 

### Return type

[**Asset**](Asset.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="updateSiteBuildHook"></a>
# **updateSiteBuildHook**
> updateSiteBuildHook(siteId, idbuildHook)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var id = "id_example"; // String | 

var buildHook = new NetlifysApiDefinition.BuildHook(); // BuildHook | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.updateSiteBuildHook(siteId, idbuildHook, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **id** | **String**|  | 
 **buildHook** | [**BuildHook**](BuildHook.md)|  | 

### Return type

null (empty response body)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="updateSiteBuildLog"></a>
# **updateSiteBuildLog**
> updateSiteBuildLog(buildId, msg)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var buildId = "buildId_example"; // String | 

var msg = new NetlifysApiDefinition.BuildLogMsg(); // BuildLogMsg | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.updateSiteBuildLog(buildId, msg, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **buildId** | **String**|  | 
 **msg** | [**BuildLogMsg**](BuildLogMsg.md)|  | 

### Return type

null (empty response body)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="updateSiteDeploy"></a>
# **updateSiteDeploy**
> Deploy updateSiteDeploy(siteId, deployId, deploy)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var deployId = "deployId_example"; // String | 

var deploy = new NetlifysApiDefinition.DeployFiles(); // DeployFiles | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.updateSiteDeploy(siteId, deployId, deploy, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **deployId** | **String**|  | 
 **deploy** | [**DeployFiles**](DeployFiles.md)|  | 

### Return type

[**Deploy**](Deploy.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="updateSiteMetadata"></a>
# **updateSiteMetadata**
> updateSiteMetadata(siteId, metadata)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var metadata = new NetlifysApiDefinition.Metadata(); // Metadata | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.updateSiteMetadata(siteId, metadata, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **metadata** | [**Metadata**](Metadata.md)|  | 

### Return type

null (empty response body)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="updateSiteSnippet"></a>
# **updateSiteSnippet**
> updateSiteSnippet(siteId, snippetIdsnippet)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var siteId = "siteId_example"; // String | 

var snippetId = "snippetId_example"; // String | 

var snippet = new NetlifysApiDefinition.Snippet(); // Snippet | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully.');
  }
};
apiInstance.updateSiteSnippet(siteId, snippetIdsnippet, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **siteId** | **String**|  | 
 **snippetId** | **String**|  | 
 **snippet** | [**Snippet**](Snippet.md)|  | 

### Return type

null (empty response body)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="uploadDeployFile"></a>
# **uploadDeployFile**
> File uploadDeployFile(deployId, path, fileBody)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var deployId = "deployId_example"; // String | 

var path = "path_example"; // String | 

var fileBody = B; // Blob | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.uploadDeployFile(deployId, path, fileBody, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deployId** | **String**|  | 
 **path** | **String**|  | 
 **fileBody** | **Blob**|  | 

### Return type

**File**

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/octet-stream
 - **Accept**: application/json

<a name="uploadDeployFunction"></a>
# **uploadDeployFunction**
> ModelFunction uploadDeployFunction(deployId, name, fileBody, opts)



### Example
```javascript
var NetlifysApiDefinition = require('netlifys_api_definition');
var defaultClient = NetlifysApiDefinition.ApiClient.instance;

// Configure OAuth2 access token for authorization: netlifyAuth
var netlifyAuth = defaultClient.authentications['netlifyAuth'];
netlifyAuth.accessToken = 'YOUR ACCESS TOKEN';

var apiInstance = new NetlifysApiDefinition.DefaultApi();

var deployId = "deployId_example"; // String | 

var name = "name_example"; // String | 

var fileBody = B; // Blob | 

var opts = { 
  'runtime': "runtime_example" // String | 
};

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.uploadDeployFunction(deployId, name, fileBody, opts, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **deployId** | **String**|  | 
 **name** | **String**|  | 
 **fileBody** | **Blob**|  | 
 **runtime** | **String**|  | [optional] 

### Return type

[**ModelFunction**](ModelFunction.md)

### Authorization

[netlifyAuth](../README.md#netlifyAuth)

### HTTP request headers

 - **Content-Type**: application/octet-stream
 - **Accept**: application/json

