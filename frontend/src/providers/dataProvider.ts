import type { DataProvider } from "@refinedev/core";
import { actionsApi, apiAxios, orderApi, unwrap } from "../api/client";

// buildListParams translates refine's pagination/filters/sorters into the
// backend's filter convention: field__op=value query params, sort=-field,field.
// Page/sort ride the generated method's typed params; dynamic field filters go
// through options.params (they are intentionally not enumerated in swagger).
function buildListParams(params: {
  pagination?: { current?: number; pageSize?: number };
  filters?: any[];
  sorters?: any[];
}): { page: number; pageSize: number; sort?: string; filterParams: Record<string, string> } {
  const page = params.pagination?.current ?? 1;
  const pageSize = params.pagination?.pageSize ?? 10;
  const filterParams: Record<string, string> = {};
  (params.filters ?? []).forEach((f) => {
    if (f.operator === "in" && Array.isArray(f.value)) {
      filterParams[`${f.field}__in`] = f.value.join(",");
    } else if (f.operator === "between" && Array.isArray(f.value)) {
      filterParams[`${f.field}__between`] = f.value.join(",");
    } else if (f.operator === "contains") {
      filterParams[`${f.field}__like`] = String(f.value);
    } else {
      const opMap: Record<string, string> = {
        eq: "eq",
        ne: "neq",
        gt: "gt",
        gte: "gte",
        lt: "lt",
        lte: "lte",
      };
      const op = opMap[f.operator as string] ?? "eq";
      filterParams[`${f.field}__${op}`] = String(f.value);
    }
  });
  const parts = (params.sorters ?? []).map(
    (s) => `${s.order === "desc" ? "-" : ""}${s.field}`
  );
  return { page, pageSize, sort: parts.length ? parts.join(",") : undefined, filterParams };
}

// dataProvider speaks the backend envelope via the swagger-generated clients:
// Actions-mode resources go through ActionsApi, the hand-written order resource
// through OrderApi. Response shape ({ list, total } / record) is unchanged.
export const dataProvider: DataProvider = {
  getList: async ({ resource, pagination, filters, sorters }) => {
    const { page, pageSize, sort, filterParams } = buildListParams({ pagination, filters, sorters });
    if (resource === "order") {
      const d = await unwrap<any>(orderApi.orderGet({ page, pageSize }));
      return { data: d.list ?? [], total: d.total ?? 0 };
    }
    const d = await unwrap<any>(
      actionsApi.resourceGet({ resource, page, pageSize, sort }, { params: filterParams })
    );
    return { data: d.list ?? [], total: d.total ?? 0 };
  },
  getOne: async ({ resource, id }) => {
    const data = await unwrap<any>(
      actionsApi.resourceIdGet({ resource, id: String(id) })
    );
    return { data };
  },
  getMany: async ({ resource, ids }) => {
    const d = await unwrap<any>(
      actionsApi.resourceGet({ resource }, { params: { id__in: ids.join(",") } })
    );
    return { data: d.list ?? [] };
  },
  create: async ({ resource, variables }) => {
    if (resource === "order") {
      const data = await unwrap<any>(orderApi.orderPost({ body: variables as any }));
      return { data };
    }
    const data = await unwrap<any>(
      actionsApi.resourcePost({ resource, body: (variables ?? {}) as object })
    );
    return { data };
  },
  update: async ({ resource, id, variables }) => {
    const data = await unwrap<any>(
      actionsApi.resourceIdPut({ resource, id: String(id), body: (variables ?? {}) as object })
    );
    return { data };
  },
  deleteOne: async ({ resource, id }) => {
    if (resource === "order") {
      const data = await unwrap<any>(orderApi.orderIdDelete({ id: String(id) }));
      return { data: data ?? {} };
    }
    const data = await unwrap<any>(
      actionsApi.resourceIdDelete({ resource, id: String(id) })
    );
    return { data: data ?? {} };
  },
  getApiUrl: () => "/api/v1",
  custom: async ({ url, method, payload, headers }) => {
    const data = await unwrap<any>(
      apiAxios.request({ url, method: method as any, data: payload, headers: headers as any })
    );
    return { data };
  },
};
