import React, { useState } from "react";
import { Button, Layout, Menu, Avatar, Dropdown, Typography, Select } from "antd";
import { MenuFoldOutlined, MenuUnfoldOutlined } from "@ant-design/icons";
import { useLocation, useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useGetIdentity, useLogout } from "@refinedev/core";
import { useTabs } from "../context/TabsContext";
import { navItems } from "../registry";
import { authApi, unwrap } from "../api/client";
import { SUPPORTED_LANGS, setLang } from "../i18n";

const { Sider, Header, Content } = Layout;

// AdminLayout is the application chrome: sider menu, header and the multi-tab
// content area. Clicking a menu item opens (or focuses) its tab. The sider menu
// is dynamic: it only shows the menus granted to the current user's role via
// the role-scoped /menu/tree API (the dashboard home is always available).
// All labels render through react-i18next (titleKey + fallback title).
export default function AdminLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const tabs = useTabs();
  const { t, i18n } = useTranslation();
  const { data: identity } = useGetIdentity();
  const { mutate: logout } = useLogout();
  const [grantedPaths, setGrantedPaths] = useState<string[] | null>(null);
  const [collapsed, setCollapsed] = useState(false);

  // 首页（/）自动打开 dashboard
  React.useEffect(() => {
    if (location.pathname === "/") {
      navigate("/dashboard", { replace: true });
    }
  }, [location.pathname, navigate]);

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
      tabs.openTab({ key: item.path, title: item.title, titleKey: item.titleKey, closable: true });
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location.pathname, visibleNavItems]);

  const menuItems = buildMenu(visibleNavItems, (key, fallback) => t(key, { defaultValue: fallback }));
  const userMenu = {
    items: [{ key: "logout", label: t("layout.logout") }],
    onClick: () => logout(),
  };

  return (
    <Layout style={{ height: "100vh" }}>
      <Sider
        theme="dark"
        width={220}
        collapsible
        collapsed={collapsed}
        onCollapse={setCollapsed}
      >
        <div
          style={{
            color: "#fff",
            padding: 16,
            fontWeight: 600,
            fontSize: 16,
            whiteSpace: "nowrap",
            overflow: "hidden",
          }}
        >
          {collapsed ? t("app.shortTitle") : t("app.title")}
        </div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={[location.pathname]}
          items={menuItems}
          onClick={({ key }) => {
            const it = visibleNavItems.find((n) => n.path === key);
            if (it) {
              tabs.openTab({ key: it.path, title: it.title, titleKey: it.titleKey, closable: true });
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
          <div style={{ display: "flex", alignItems: "center", gap: 12 }}>
            <Button
              type="text"
              icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
              onClick={() => setCollapsed(!collapsed)}
            />
            <Typography.Text strong>{t("app.subtitle")}</Typography.Text>
          </div>
          <div style={{ display: "flex", alignItems: "center", gap: 12 }}>
            <Select
              size="small"
              value={i18n.language}
              aria-label={t("layout.language")}
              options={SUPPORTED_LANGS.map((l) => ({ ...l }))}
              onChange={(v) => setLang(v)}
              style={{ width: 96 }}
            />
            <Dropdown menu={userMenu}>
              <span style={{ cursor: "pointer" }}>
                <Avatar style={{ marginRight: 8 }}>
                  {identity ? String((identity as any).name)[0].toUpperCase() : "U"}
                </Avatar>
                {(identity as any)?.name}
              </span>
            </Dropdown>
          </div>
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

function buildMenu(items: typeof navItems, t: (key: string, fallback: string) => string) {
  const groups = new Map<string, { label: string; children: any[] }>();
  const roots: any[] = [];
  items.forEach((it) => {
    const mi = { key: it.path, label: t(it.titleKey, it.title) };
    if (it.group) {
      const gk = it.groupKey ?? it.group;
      if (!groups.has(gk)) groups.set(gk, { label: t(gk, it.group), children: [] });
      groups.get(gk)!.children.push(mi);
    } else {
      roots.push(mi);
    }
  });
  const out: any[] = [...roots];
  groups.forEach(({ label, children }, g) => out.push({ key: `g-${g}`, label, children }));
  return out;
}
