// Resource metadata describes a CRUD resource to the auto-generated pages. One
// description drives the list table, the search form and the create/edit form —
// this is the front-end counterpart of the backend's Actions mode (zero-code).
//
//   - name : refine resource id AND backend path segment (e.g. "demo/product")
//   - path : route path used by the multi-tab router (e.g. "/demo/product")
//   - fields: columns / form items

export type FieldType =
  | "text"
  | "number"
  | "textarea"
  | "password"
  | "switch"
  | "datetime"
  | "select"
  | "menutree"; // menu-permission tree (role -> menu binding)

export interface FieldMeta {
  name: string; // backend field name (used in filters/sorters)
  label: string; // fallback text; UI renders resources.${meta}.fields.${name}
  type: FieldType;
  searchable?: boolean; // show in the search form
  sortable?: boolean; // allow column sorting
  required?: boolean; // required on create/edit
  options?: OptionMeta[]; // for select
  hideInForm?: boolean;
  hideInTable?: boolean;
  width?: number;
}

export interface ResourceMeta {
  name: string;
  title: string; // fallback text; UI renders t(`resources.${i18nKey}.title`)
  i18nKey: string; // e.g. "product" -> title/fields resolve under resources.*
  path: string;
  icon?: string;
  group?: string; // menu group
  // Backend write-permission prefix (F-node perms). When set, the auto page
  // hides 新建/编辑/删除 for users lacking <prefix>:create/:edit/:delete.
  permPrefix?: string;
  fields: FieldMeta[];
}

// Field labels resolve as t(`resources.${meta.i18nKey}.fields.${name}`) with
// label as the fallback, so FieldMeta needs no key of its own. Select/switch
// options are shared constants, so they carry explicit labelKeys.
export interface OptionMeta {
  label: string; // fallback text
  labelKey: string; // e.g. "common.status.enabled"
  value: any;
}

const statusOptions: OptionMeta[] = [
  { label: "启用", labelKey: "common.status.enabled", value: 1 },
  { label: "禁用", labelKey: "common.status.disabled", value: 0 },
];

const productStatusOptions: OptionMeta[] = [
  { label: "在售", labelKey: "common.productStatus.onSale", value: 1 },
  { label: "下架", labelKey: "common.productStatus.offSale", value: 0 },
];

const yesNoOptions: OptionMeta[] = [
  { label: "是", labelKey: "common.yes", value: "Y" },
  { label: "否", labelKey: "common.no", value: "N" },
];

// --- Actions-mode resources (backed by zero-code CRUD on the backend) ---
export const productMeta: ResourceMeta = {
  name: "demo/product", i18nKey: "product", permPrefix: "demo:product",

  title: "商品管理",
  path: "/demo/product",
  icon: "appstore",
  group: "业务示例",
  fields: [
    { name: "id", label: "ID", type: "number", hideInForm: true, sortable: true, width: 70 },
    { name: "name", label: "名称", type: "text", searchable: true, sortable: true },
    { name: "sku", label: "SKU", type: "text", searchable: true },
    { name: "category", label: "分类", type: "text", searchable: true },
    { name: "price", label: "价格", type: "number", sortable: true },
    { name: "stock", label: "库存", type: "number", sortable: true },
    { name: "status", label: "状态", type: "switch", options: productStatusOptions, width: 90 },
    { name: "remark", label: "备注", type: "textarea", hideInTable: true },
  ],
};

export const userMeta: ResourceMeta = {
  name: "sys/user", i18nKey: "user", permPrefix: "sys:user",

  title: "用户管理",
  path: "/system/user",
  icon: "user",
  group: "系统管理",
  fields: [
    { name: "id", label: "ID", type: "number", hideInForm: true, width: 70 },
    { name: "username", label: "用户名", type: "text", searchable: true, required: true },
    { name: "password", label: "密码", type: "password", hideInTable: true },
    { name: "nickName", label: "昵称", type: "text" },
    { name: "phone", label: "手机号", type: "text", searchable: true },
    { name: "email", label: "邮箱", type: "text" },
    { name: "status", label: "状态", type: "switch", options: statusOptions, width: 90 },
    { name: "remark", label: "备注", type: "textarea", hideInTable: true },
  ],
};

export const roleMeta: ResourceMeta = {
  name: "sys/role", i18nKey: "role", permPrefix: "sys:role",

  title: "角色管理",
  path: "/system/role",
  icon: "safety",
  group: "系统管理",
  fields: [
    { name: "id", label: "ID", type: "number", hideInForm: true, width: 70 },
    { name: "name", label: "角色名", type: "text", searchable: true, required: true },
    { name: "roleKey", label: "角色标识", type: "text", required: true },
    { name: "status", label: "状态", type: "switch", options: statusOptions, width: 90 },
    { name: "sort", label: "排序", type: "number" },
    { name: "menuIds", label: "菜单权限", type: "menutree", hideInTable: true },
    { name: "remark", label: "备注", type: "textarea", hideInTable: true },
  ],
};

export const deptMeta: ResourceMeta = {
  name: "sys/dept", i18nKey: "dept", permPrefix: "sys:dept",

  title: "部门管理",
  path: "/system/dept",
  icon: "apartment",
  group: "系统管理",
  fields: [
    { name: "id", label: "ID", type: "number", hideInForm: true, width: 70 },
    { name: "name", label: "部门名", type: "text", searchable: true, required: true },
    { name: "leader", label: "负责人", type: "text" },
    { name: "phone", label: "电话", type: "text" },
    { name: "status", label: "状态", type: "switch", options: statusOptions, width: 90 },
  ],
};

export const dictMeta: ResourceMeta = {
  name: "sys/dict", i18nKey: "dict", permPrefix: "sys:dict",

  title: "字典管理",
  path: "/system/dict",
  icon: "book",
  group: "系统管理",
  fields: [
    { name: "id", label: "ID", type: "number", hideInForm: true, width: 70 },
    { name: "name", label: "字典名", type: "text", searchable: true, required: true },
    { name: "type", label: "类型", type: "text", searchable: true, required: true },
    { name: "status", label: "状态", type: "switch", options: statusOptions, width: 90 },
    { name: "remark", label: "备注", type: "textarea", hideInTable: true },
  ],
};

// System parameter: Actions-mode resource (backend zero-code CRUD on sys/config).
// configValue stays a plain input (values mix text/numbers/flags).
export const configMeta: ResourceMeta = {
  name: "sys/config", i18nKey: "config", permPrefix: "sys:config",

  title: "参数管理",
  path: "/system/config",
  icon: "config",
  group: "系统管理",
  fields: [
    { name: "id", label: "ID", type: "number", hideInForm: true, width: 70 },
    { name: "configName", label: "参数名", type: "text", searchable: true, required: true },
    { name: "configKey", label: "参数键", type: "text", searchable: true, required: true },
    { name: "configValue", label: "参数值", type: "text" },
    { name: "configType", label: "内置", type: "select", options: yesNoOptions, width: 90 },
    { name: "isFrontend", label: "前台可读", type: "switch", width: 100 },
    { name: "remark", label: "备注", type: "textarea", hideInTable: true },
  ],
};
