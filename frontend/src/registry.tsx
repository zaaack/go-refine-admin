import React from "react";
import { productMeta, userMeta, roleMeta, deptMeta, dictMeta } from "./resources/meta";
import ResourcePage from "./pages/ResourcePage";
import OrderPage from "./pages/OrderPage";
import Dashboard from "./pages/Dashboard";

// navItems is the single source of truth for both the sider menu and the
// multi-tab content. Each entry maps a route path to a rendered page element.
// Titles are i18n keys rendered with useTranslation (title is the fallback):
//   t(titleKey, { defaultValue: title })
export interface NavItem {
  path: string;
  title: string;
  titleKey: string;
  icon?: string;
  group?: string;
  groupKey?: string;
  element: React.ReactNode;
}

export const navItems: NavItem[] = [
  { path: "/dashboard", title: "仪表盘", titleKey: "nav.dashboard", icon: "dashboard", element: <Dashboard /> },
  {
    path: productMeta.path,
    title: productMeta.title,
    titleKey: `resources.${productMeta.i18nKey}.title`,
    icon: productMeta.icon,
    group: productMeta.group,
    groupKey: "nav.group.business",
    element: <ResourcePage meta={productMeta} />,
  },
  { path: "/demo/order", title: "订单管理", titleKey: "nav.order", icon: "shopping", group: "业务示例", groupKey: "nav.group.business", element: <OrderPage /> },
  {
    path: userMeta.path,
    title: userMeta.title,
    titleKey: `resources.${userMeta.i18nKey}.title`,
    icon: userMeta.icon,
    group: userMeta.group,
    groupKey: "nav.group.system",
    element: <ResourcePage meta={userMeta} />,
  },
  {
    path: roleMeta.path,
    title: roleMeta.title,
    titleKey: `resources.${roleMeta.i18nKey}.title`,
    icon: roleMeta.icon,
    group: roleMeta.group,
    groupKey: "nav.group.system",
    element: <ResourcePage meta={roleMeta} />,
  },
  {
    path: deptMeta.path,
    title: deptMeta.title,
    titleKey: `resources.${deptMeta.i18nKey}.title`,
    icon: deptMeta.icon,
    group: roleMeta.group,
    groupKey: "nav.group.system",
    element: <ResourcePage meta={deptMeta} />,
  },
  {
    path: dictMeta.path,
    title: dictMeta.title,
    titleKey: `resources.${dictMeta.i18nKey}.title`,
    icon: dictMeta.icon,
    group: dictMeta.group,
    groupKey: "nav.group.system",
    element: <ResourcePage meta={dictMeta} />,
  },
];

export function findNav(path: string): NavItem | undefined {
  return navItems.find((n) => n.path === path);
}
