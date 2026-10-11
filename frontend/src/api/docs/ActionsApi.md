# ActionsApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**resourceGet**](#resourceget) | **GET** /{resource} | List resources (Actions mode)|
|[**resourceIdDelete**](#resourceiddelete) | **DELETE** /{resource}/{id} | Delete resource (Actions mode)|
|[**resourceIdGet**](#resourceidget) | **GET** /{resource}/{id} | Get one resource (Actions mode)|
|[**resourceIdPut**](#resourceidput) | **PUT** /{resource}/{id} | Update resource (Actions mode)|
|[**resourcePost**](#resourcepost) | **POST** /{resource} | Create resource (Actions mode)|

# **resourceGet**
> ResponseEnvelope resourceGet()


### Example

```typescript
import {
    ActionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ActionsApi(configuration);

let resource: string; //resource name, e.g. demo/product (default to undefined)
let page: number; //page (optional) (default to undefined)
let pageSize: number; //page size (optional) (default to undefined)
let sort: string; //sort, e.g. -createdAt,name (optional) (default to undefined)

const { status, data } = await apiInstance.resourceGet(
    resource,
    page,
    pageSize,
    sort
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **resource** | [**string**] | resource name, e.g. demo/product | defaults to undefined|
| **page** | [**number**] | page | (optional) defaults to undefined|
| **pageSize** | [**number**] | page size | (optional) defaults to undefined|
| **sort** | [**string**] | sort, e.g. -createdAt,name | (optional) defaults to undefined|


### Return type

**ResponseEnvelope**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **resourceIdDelete**
> ResponseEnvelope resourceIdDelete()


### Example

```typescript
import {
    ActionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ActionsApi(configuration);

let resource: string; //resource name, e.g. demo/product (default to undefined)
let id: string; //primary key(s), comma-separated (default to undefined)

const { status, data } = await apiInstance.resourceIdDelete(
    resource,
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **resource** | [**string**] | resource name, e.g. demo/product | defaults to undefined|
| **id** | [**string**] | primary key(s), comma-separated | defaults to undefined|


### Return type

**ResponseEnvelope**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **resourceIdGet**
> ResponseEnvelope resourceIdGet()


### Example

```typescript
import {
    ActionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ActionsApi(configuration);

let resource: string; //resource name, e.g. demo/product (default to undefined)
let id: string; //primary key (default to undefined)

const { status, data } = await apiInstance.resourceIdGet(
    resource,
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **resource** | [**string**] | resource name, e.g. demo/product | defaults to undefined|
| **id** | [**string**] | primary key | defaults to undefined|


### Return type

**ResponseEnvelope**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **resourceIdPut**
> ResponseEnvelope resourceIdPut(body)


### Example

```typescript
import {
    ActionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ActionsApi(configuration);

let resource: string; //resource name, e.g. demo/product (default to undefined)
let id: string; //primary key (default to undefined)
let body: object; //record

const { status, data } = await apiInstance.resourceIdPut(
    resource,
    id,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **object**| record | |
| **resource** | [**string**] | resource name, e.g. demo/product | defaults to undefined|
| **id** | [**string**] | primary key | defaults to undefined|


### Return type

**ResponseEnvelope**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **resourcePost**
> ResponseEnvelope resourcePost(body)


### Example

```typescript
import {
    ActionsApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new ActionsApi(configuration);

let resource: string; //resource name, e.g. demo/product (default to undefined)
let body: object; //record

const { status, data } = await apiInstance.resourcePost(
    resource,
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **object**| record | |
| **resource** | [**string**] | resource name, e.g. demo/product | defaults to undefined|


### Return type

**ResponseEnvelope**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

