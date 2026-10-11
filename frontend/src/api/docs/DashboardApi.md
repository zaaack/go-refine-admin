# DashboardApi

All URIs are relative to */api/v1*

|Method | HTTP request | Description|
|------------- | ------------- | -------------|
|[**dashboardStatsGet**](#dashboardstatsget) | **GET** /dashboard/stats | Dashboard stats|

# **dashboardStatsGet**
> ResponseEnvelope dashboardStatsGet()


### Example

```typescript
import {
    DashboardApi,
    Configuration
} from './api';

const configuration = new Configuration();
const apiInstance = new DashboardApi(configuration);

const { status, data } = await apiInstance.dashboardStatsGet();
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

