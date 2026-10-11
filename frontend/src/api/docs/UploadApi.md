# UploadApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**uploadPost**](#uploadpost) | **POST** /upload | Upload a file|

# **uploadPost**
> ResponseEnvelope uploadPost()


### Example

```typescript
import {
    UploadApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new UploadApi(configuration);

let file: File; //file (default to undefined)

const { status, data } = await apiInstance.uploadPost(
    file
);
```

### Parameters

|Name | Type | Description  | Notes|
|------------- | ------------- | ------------- | -------------|
| **file** | [**File**] | file | defaults to undefined|


### Return type

**ResponseEnvelope**

### Authorization

[Bearer](../README.md#Bearer)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
|**200** | OK |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

