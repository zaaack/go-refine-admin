# OrderApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**orderGet**](#orderget) | **GET** /order | List orders (regular mode) with summary|
|[**orderIdDelete**](#orderiddelete) | **DELETE** /order/{id} | Delete order|
|[**orderIdPaidPut**](#orderidpaidput) | **PUT** /order/{id}/paid | Mark order paid|
|[**orderPost**](#orderpost) | **POST** /order | Create order|
|[**orderStatsGet**](#orderstatsget) | **GET** /order/stats | Order stats by status|

# **orderGet**
> ResponseEnvelope orderGet()


### Example

```typescript
import {
    OrderApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OrderApi(configuration);

let minAmount: number; //min amount (optional) (default to undefined)
let status: string; //status, repeatable (optional) (default to undefined)
let channel: string; //channel (optional) (default to undefined)
let start: string; //created >= (RFC3339) (optional) (default to undefined)
let end: string; //created <= (RFC3339) (optional) (default to undefined)
let page: number; //page (optional) (default to undefined)
let pageSize: number; //page size (optional) (default to undefined)

const { status, data } = await apiInstance.orderGet(
    minAmount,
    status,
    channel,
    start,
    end,
    page,
    pageSize
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **minAmount** | [**number**] | min amount | (optional) defaults to undefined|
| **status** | [**string**] | status, repeatable | (optional) defaults to undefined|
| **channel** | [**string**] | channel | (optional) defaults to undefined|
| **start** | [**string**] | created &gt;&#x3D; (RFC3339) | (optional) defaults to undefined|
| **end** | [**string**] | created &lt;&#x3D; (RFC3339) | (optional) defaults to undefined|
| **page** | [**number**] | page | (optional) defaults to undefined|
| **pageSize** | [**number**] | page size | (optional) defaults to undefined|


### Return type

**ResponseEnvelope**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **orderIdDelete**
> ResponseEnvelope orderIdDelete()


### Example

```typescript
import {
    OrderApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OrderApi(configuration);

let id: string; //order id (default to undefined)

const { status, data } = await apiInstance.orderIdDelete(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | order id | defaults to undefined|


### Return type

**ResponseEnvelope**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: */*


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **orderIdPaidPut**
> ResponseEnvelope orderIdPaidPut()


### Example

```typescript
import {
    OrderApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OrderApi(configuration);

let id: string; //order id (default to undefined)

const { status, data } = await apiInstance.orderIdPaidPut(
    id
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **id** | [**string**] | order id | defaults to undefined|


### Return type

**ResponseEnvelope**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **orderPost**
> ResponseEnvelope orderPost(body)


### Example

```typescript
import {
    OrderApi,
    Configuration,
    ModelsOrder
} from './api';

const configuration = new Configuration();
const apiInstance = new OrderApi(configuration);

let body: ModelsOrder; //order

const { status, data } = await apiInstance.orderPost(
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **ModelsOrder**| order | |


### Return type

**ResponseEnvelope**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **orderStatsGet**
> ResponseEnvelope orderStatsGet()


### Example

```typescript
import {
    OrderApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new OrderApi(configuration);

const { status, data } = await apiInstance.orderStatsGet();
```

### Parameters
This endpoint does not have any parameters.


### Return type

**ResponseEnvelope**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

