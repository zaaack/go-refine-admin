# AuthApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**captchaGet**](#captchaget) | **GET** /captcha | 获取验证码|
|[**loginPost**](#loginpost) | **POST** /login | Login|
|[**meGet**](#meget) | **GET** /me | Current user|
|[**menuTreeGet**](#menutreeget) | **GET** /menu/tree | Menu tree (role-scoped)|

# **captchaGet**
> ResponseEnvelope captchaGet()


### Example

```typescript
import {
    AuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

const { status, data } = await apiInstance.captchaGet();
```

### Parameters
This endpoint does not have any parameters.


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

# **loginPost**
> ResponseEnvelope loginPost(body)


### Example

```typescript
import {
    AuthApi,
    Configuration,
    ControllersLoginRequest
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

let body: ControllersLoginRequest; //credentials

const { status, data } = await apiInstance.loginPost(
    body
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **body** | **ControllersLoginRequest**| credentials | |


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

# **meGet**
> ResponseEnvelope meGet()


### Example

```typescript
import {
    AuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

const { status, data } = await apiInstance.meGet();
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

# **menuTreeGet**
> ResponseEnvelope menuTreeGet()


### Example

```typescript
import {
    AuthApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new AuthApi(configuration);

const { status, data } = await apiInstance.menuTreeGet();
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

