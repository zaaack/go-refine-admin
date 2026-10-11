import React from "react";
import { productMeta, userMeta, roleMeta, deptMeta, dictMeta } from "./resources/meta";
import ResourcePage from "./pages/ResourcePage";
import OrderPage from "./pages/OrderPage";
import Dashboard from "./pages/Dashboard";

// navItems is the single source of truth for both the sider menu and the
// multi-tab content. Each entry maps a route path to a rendered page element.
export interface NavItem {
  path: string;
  title: string;
  icon?: string;
  group?: string;
  element: React.ReactNode;
}

export const navItems: NavItem[] = [
  { path: "/dashboard", title: "仪表盘", icon: "dashboard", element: <Dashboard /> },
  {
    path: productMeta.path,
    title: productMeta.title,
    icon: productMeta.icon,
    group: productMeta.group,
    element: <ResourcePage meta={productMeta} />,
  },
  { path: "/demo/order", title: "订单管理", icon: "shopping", group: "业务示例", element: <OrderPage /> },
  {
    path: userMeta.path,
    title: userMeta.title,
    icon: userMeta.icon,
    group: userMeta.group,
    element: <ResourcePage meta={userMeta} />,
  },
  {
    path: roleMeta.path,
    title: roleMeta.title,
    icon: roleMeta.icon,
    group: roleMeta.group,
    element: <ResourcePage meta={roleMeta} />,
  },
  {
    path: deptMeta.path,
    title: deptMeta.title,
    icon: deptMeta.icon,
    group: roleMeta.group,
    element: <ResourcePage meta={deptMeta} />,
  },
  {
    path: dictMeta.path,
    title: dictMeta.title,
    icon: dictMeta.icon,
    group: dictMeta.group,
    element: <ResourcePage meta={dictMeta} />,
  },
];

export function findNav(path: string): NavItem | undefined {
  return navItems.find((n) => n.path === path);
}
