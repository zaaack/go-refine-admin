import React, { useMemo, useState } from "react";
import { Layout, Menu, Avatar, Dropdown, Typography } from "antd";
import { useLocation, useNavigate } from "react-router-dom";
import { useGetIdentity, useLogout } from "@refinedev/core";
import { useTabs } from "../context/TabsContext";
import { navItems } from "../registry";
import { authApi, unwrap } from "../api/client";

const { Sider, Header, Content } = Layout;

// AdminLayout is the application chrome: sider menu, header and the multi-tab
// content area. Clicking a menu item opens (or focuses) its tab. The sider menu
// is dynamic: it only shows the menus granted to the current user's role via
// the role-scoped /menu/tree API (the dashboard home is always available).
export default function AdminLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const tabs = useTabs();
  const { data: identity } = useGetIdentity();
  const { mutate: logout } = useLogout();
  const [grantedPaths, setGrantedPaths] = useState<string[] | null>(null);

  // Load the menus granted to the current role; null means "load failed, show all".
  React.useEffect(() => {
    unwrap<any[]>(authApi.menuTreeGet())
      .then((data) => {
        const paths: string[] = [];
        const walk = (nodes: any[]) =>
          nodes.forEach((n: any) => {
            if (n.path) paths.push(n.path);
            if (n.children?.length) walk(n.children);
          });
        walk(data ?? []);
        setGrantedPaths(paths);
      })
      .catch(() => setGrantedPaths(null));
  }, []);

  const visibleNavItems = React.useMemo(() => {
    if (grantedPaths === null) return navItems;
    return navItems.filter((n) => n.path === "/dashboard" || grantedPaths.includes(n.path));
  }, [grantedPaths]);

  // Make sure the route we landed on is represented as an open tab.
  React.useEffect(() => {
    const item = visibleNavItems.find((n) => n.path === location.pathname);
    if (item && !tabs.opened.some((t) => t.key === item.path)) {
      tabs.openTab({ key: item.path, title: item.title, closable: true });
    }
  }, [location.pathname, tabs, visibleNavItems]);

  const menuItems = buildMenu(visibleNavItems);
  const userMenu = {
    items: [{ key: "logout", label: "退出登录" }],
    onClick: () => logout(),
  };

  return (
    <Layout style={{ height: "100vh" }}>
      <Sider theme="dark" width={220}>
        <div style={{ color: "#fff", padding: 16, fontWeight: 600, fontSize: 16 }}>
          Go Refine Admin
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={({ key }) => {
            const it = visibleNavItems.find((n) => n.path === key);
            if (it) {
              tabs.openTab({ key: it.path, title: it.title, closable: true });
              navigate(it.path);
            }
          }}
        />
      </Sider>
      <Layout>
        <Header
          style={{
            background: "#fff",
            display: "flex",
            justifyContent: "space-between",
            alignItems: "center",
            padding: "0 16px",
          }}
        >
          <Typography.Text strong>企业级后台脚手架 · Actions / 常规双轨</Typography.Text>
          <Dropdown menu={userMenu}>
            <span style={{ cursor: "pointer" }}>
              <Avatar style={{ marginRight: 8 }}>
                {identity ? String((identity as any).name)[0].toUpperCase() : "U"}
              </Avatar>
              {(identity as any)?.name}
            </span>
          </Dropdown>
        </Header>
        <Content style={{ overflow: "hidden" }}>
          <TabsViewMount />
        </Content>
      </Layout>
    </Layout>
  );
}

// TabsView is imported lazily to keep this file focused.
import TabsView from "./TabsView";
function TabsViewMount() {
  return <TabsView />;
}

function buildMenu(items: typeof navItems) {
  const groups = new Map<string, any[]>();
  const roots: any[] = [];
  items.forEach((it) => {
    const mi = { key: it.path, label: it.title };
    if (it.group) {
      if (!groups.has(it.group)) groups.set(it.group, []);
      groups.get(it.group)!.push(mi);
    } else {
      roots.push(mi);
    }
  });
  const out: any[] = [...roots];
  groups.forEach((children, g) => out.push({ key: `g-${g}`, label: g, children }));
  return out;
}
