import React, { useState } from "react";
import {
  Breadcrumb,
  Button,
  Layout,
  Menu,
  Avatar,
  Dropdown,
  Typography,
  Select,
  Modal,
  Input,
  List,
  Form,
  Table,
  Tag,
  message,
  Tooltip,
} from "antd";
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  SearchOutlined,
  DashboardOutlined,
  AppstoreOutlined,
  ShoppingOutlined,
  UserOutlined,
  SafetyOutlined,
  ApartmentOutlined,
  BookOutlined,
  SettingOutlined,
  ShopOutlined,
  FolderOutlined,
  MenuOutlined,
  HomeOutlined,
  FileTextOutlined,
  LoginOutlined,
  TeamOutlined,
  CloudServerOutlined,
  ToolOutlined,
} from "@ant-design/icons";
import { useLocation, useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useGetIdentity, useLogout } from "@refinedev/core";
import { useTabs } from "../context/TabsContext";
import { navItems, type NavItem } from "../registry";
import { authApi, unwrap, apiAxios } from "../api/client";
import { backendErrorMessage } from "../i18n";
import { SUPPORTED_LANGS, setLang } from "../i18n";

const { Sider, Header, Content } = Layout;

// icon string (from registry/meta) -> antd icon node. Covers every icon used
// in registry.tsx / resources/meta.ts plus the menu group icons below.
const ICONS: Record<string, React.ReactNode> = {
  dashboard: <DashboardOutlined />,
  appstore: <AppstoreOutlined />,
  shopping: <ShoppingOutlined />,
  user: <UserOutlined />,
  safety: <SafetyOutlined />,
  apartment: <ApartmentOutlined />,
  book: <BookOutlined />,
  menu: <MenuOutlined />,
  log: <FileTextOutlined />,
  login: <LoginOutlined />,
  online: <TeamOutlined />,
  server: <CloudServerOutlined />,
  config: <ToolOutlined />,
};

function iconOf(name?: string): React.ReactNode {
  if (!name) return <FolderOutlined />;
  return ICONS[name] ?? <FolderOutlined />;
}

const GROUP_ICONS: Record<string, React.ReactNode> = {
  "nav.group.business": <ShopOutlined />,
  "nav.group.system": <SettingOutlined />,
};

function groupIcon(key: string): React.ReactNode {
  return GROUP_ICONS[key] ?? <FolderOutlined />;
}

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
  const [searchOpen, setSearchOpen] = useState(false);
  const [searchText, setSearchText] = useState("");
  const [profileOpen, setProfileOpen] = useState(false);
  const [pwdOpen, setPwdOpen] = useState(false);
  const [loginLogOpen, setLoginLogOpen] = useState(false);
  const [myLogs, setMyLogs] = useState<any[]>([]);
  const [myLogsTotal, setMyLogsTotal] = useState(0);
  const [myLogsPage, setMyLogsPage] = useState(1);
  const [myLogsLoading, setMyLogsLoading] = useState(false);
  const [profile, setProfile] = useState<any>(null);
  const [profileForm] = Form.useForm();
  const [pwdForm] = Form.useForm();

  const tr = (it: Pick<NavItem, "title" | "titleKey">) => t(it.titleKey, { defaultValue: it.title });

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

  const goPath = (path: string) => {
    const it = visibleNavItems.find((n) => n.path === path);
    if (it) {
      tabs.openTab({ key: it.path, title: it.title, titleKey: it.titleKey, closable: true });
      navigate(it.path);
    } else {
      navigate(path);
    }
  };

  // Ctrl/Cmd+K opens the quick-search modal.
  React.useEffect(() => {
    const onKey = (e: KeyboardEvent) => {
      if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === "k") {
        e.preventDefault();
        setSearchText("");
        setSearchOpen(true);
      }
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, []);

  // Breadcrumb: 首页(/dashboard) / 分组 / 当前页面, each node with a path jumps.
  const crumbs = React.useMemo(() => {
    const cur = visibleNavItems.find((n) => n.path === location.pathname);
    const home = visibleNavItems.find((n) => n.path === "/dashboard");
    const items: { title: string; path?: string }[] = [];
    if (!cur) {
      return [{ title: t("nav.dashboard", { defaultValue: "仪表盘" }), path: "/dashboard" }];
    }
    if (cur.path === "/dashboard") {
      return [{ title: tr(cur) }];
    }
    if (home) items.push({ title: tr(home), path: home.path });
    else items.push({ title: t("nav.dashboard", { defaultValue: "仪表盘" }), path: "/dashboard" });
    if (cur.group) items.push({ title: t(cur.groupKey ?? cur.group, { defaultValue: cur.group }) });
    items.push({ title: tr(cur), path: cur.path });
    return items;
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location.pathname, visibleNavItems, i18n.language]);

  const searchResults = React.useMemo(() => {
    const kw = searchText.trim().toLowerCase();
    if (!kw) return visibleNavItems;
    return visibleNavItems.filter((n) => {
      const label = t(n.titleKey, { defaultValue: n.title }).toLowerCase();
      const group = n.group ? t(n.groupKey ?? n.group, { defaultValue: n.group }).toLowerCase() : "";
      return (
        label.includes(kw) || n.title.toLowerCase().includes(kw) || n.path.toLowerCase().includes(kw) || group.includes(kw)
      );
    });
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [searchText, visibleNavItems, i18n.language]);

  const openProfile = async () => {
    setProfileOpen(true);
    try {
      const res = await apiAxios.get("/api/v1/me/profile");
      const data = (res.data as any)?.data ?? res.data;
      setProfile(data);
      profileForm.setFieldsValue({
        nickName: data?.nickName ?? "",
        phone: data?.phone ?? "",
        email: data?.email ?? "",
        remark: data?.remark ?? "",
      });
    } catch (e: any) {
      message.error(backendErrorMessage(e));
    }
  };

  const saveProfile = async () => {
    try {
      const v = await profileForm.validateFields();
      await apiAxios.put("/api/v1/me/profile", v);
      message.success(t("layout.updateSuccess"));
      setProfileOpen(false);
    } catch (e: any) {
      if (e?.errorFields) return;
      message.error(backendErrorMessage(e));
    }
  };

  const openMyLoginLog = async (p: number) => {
    const name = String((identity as any)?.name ?? "");
    if (!name) return;
    setLoginLogOpen(true);
    setMyLogsLoading(true);
    try {
      const res = await apiAxios.get("/api/v1/sys/login-log", {
        params: { username__eq: name, page: p, pageSize: 10 },
      });
      const d = (res.data as any)?.data ?? {};
      setMyLogs(d.list ?? []);
      setMyLogsTotal(d.total ?? 0);
      setMyLogsPage(p);
    } catch (e: any) {
      message.error(backendErrorMessage(e));
    } finally {
      setMyLogsLoading(false);
    }
  };

  const savePassword = async () => {
    try {
      const v = await pwdForm.validateFields();
      await apiAxios.put("/api/v1/me/password", {
        oldPassword: v.oldPassword,
        newPassword: v.newPassword,
      });
      message.success(t("layout.pwdSuccess"));
      pwdForm.resetFields();
      setPwdOpen(false);
    } catch (e: any) {
      if (e?.errorFields) return;
      message.error(backendErrorMessage(e));
    }
  };

  const confirmLogout = () => {
    Modal.confirm({
      title: t("layout.logoutConfirmTitle"),
      content: t("layout.logoutConfirmContent"),
      okText: t("layout.logout"),
      cancelText: t("layout.cancel"),
      okType: "danger",
      onOk: () => logout(),
    });
  };

  const menuItems = buildMenu(visibleNavItems, (key, fallback) => t(key, { defaultValue: fallback }));
  const userMenu = {
    items: [
      { key: "profile", label: t("layout.profile") },
      { key: "password", label: t("layout.changePassword") },
      { key: "loginlog", label: t("layout.loginLog") },
      { type: "divider" as const },
      { key: "logout", label: t("layout.logout"), danger: true },
    ],
    onClick: ({ key }: { key: string }) => {
      if (key === "profile") void openProfile();
      else if (key === "loginlog") void openMyLoginLog(1);
      else if (key === "password") {
        pwdForm.resetFields();
        setPwdOpen(true);
      } else if (key === "logout") confirmLogout();
    },
  };

  return (
    <Layout style={{ height: "100vh" }}>
      <Sider
        theme="dark"
        width={220}
        collapsedWidth={64}
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
          onClick={({ key }) => goPath(key)}
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
          <div style={{ display: "flex", alignItems: "center", gap: 12, minWidth: 0 }}>
            <Button
              type="text"
              icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
              onClick={() => setCollapsed(!collapsed)}
            />
            <Breadcrumb
              items={crumbs.map((c, i) => ({
                title:
                  c.path && i < crumbs.length - 1 ? (
                    // eslint-disable-next-line jsx-a11y/anchor-is-valid
                    <a
                      onClick={() => goPath(c.path!)}
                      style={{ display: "inline-flex", alignItems: "center", gap: 4 }}
                    >
                      {i === 0 && <HomeOutlined />}
                      {c.title}
                    </a>
                  ) : c.path ? (
                    <a onClick={() => goPath(c.path!)}>{c.title}</a>
                  ) : (
                    c.title
                  ),
              }))}
              style={{ whiteSpace: "nowrap" }}
            />
            <Tooltip title={t("layout.searchTip")}>
              <Button type="text" icon={<SearchOutlined />} onClick={() => { setSearchText(""); setSearchOpen(true); }} />
            </Tooltip>
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

      {/* 快捷搜索:过滤菜单并跳转 */}
      <Modal
        title={t("layout.searchTitle")}
        open={searchOpen}
        footer={null}
        onCancel={() => setSearchOpen(false)}
        width={520}
        destroyOnClose
      >
        <Input
          autoFocus
          allowClear
          prefix={<SearchOutlined />}
          placeholder={t("layout.searchPlaceholder")}
          value={searchText}
          onChange={(e) => setSearchText(e.target.value)}
          onPressEnter={() => {
            if (searchResults.length) {
              setSearchOpen(false);
              goPath(searchResults[0].path);
            }
          }}
        />
        <List
          style={{ marginTop: 12, maxHeight: 360, overflow: "auto" }}
          dataSource={searchResults}
          locale={{ emptyText: t("layout.searchEmpty") }}
          renderItem={(n) => {
            const label = t(n.titleKey, { defaultValue: n.title });
            const group = n.group ? t(n.groupKey ?? n.group, { defaultValue: n.group }) : "";
            return (
              <List.Item
                style={{ cursor: "pointer", padding: "8px 12px" }}
                onClick={() => {
                  setSearchOpen(false);
                  goPath(n.path);
                }}
              >
                <List.Item.Meta
                  avatar={iconOf(n.icon)}
                  title={label}
                  description={
                    <Typography.Text type="secondary" style={{ fontSize: 12 }}>
                      {group ? `${group} · ` : ""}{n.path}
                    </Typography.Text>
                  }
                />
              </List.Item>
            );
          }}
        />
      </Modal>

      {/* 个人信息 */}
      <Modal
        title={t("layout.profileTitle")}
        open={profileOpen}
        onOk={saveProfile}
        onCancel={() => setProfileOpen(false)}
        okText={t("layout.save")}
        cancelText={t("layout.cancel")}
        destroyOnClose
      >
        <div style={{ display: "flex", alignItems: "center", gap: 12, marginBottom: 16 }}>
          <Avatar size={48}>
            {String(profile?.username ?? (identity as any)?.name ?? "U")[0].toUpperCase()}
          </Avatar>
          <div>
            <Typography.Text strong>{profile?.username ?? (identity as any)?.name}</Typography.Text>
            <br />
            <Typography.Text type="secondary" style={{ fontSize: 12 }}>
              {profile?.roleName ?? profile?.roleKey ?? ""}
            </Typography.Text>
          </div>
        </div>
        <Form form={profileForm} layout="vertical" preserve={false}>
          <Form.Item label={t("layout.nickName")} name="nickName">
            <Input />
          </Form.Item>
          <Form.Item label={t("layout.phone")} name="phone">
            <Input />
          </Form.Item>
          <Form.Item label={t("layout.email")} name="email">
            <Input />
          </Form.Item>
          <Form.Item label={t("layout.remark")} name="remark">
            <Input.TextArea rows={2} />
          </Form.Item>
        </Form>
      </Modal>

      {/* 我的登录记录:仅当前用户的登录日志(只读) */}
      <Modal
        title={t("layout.loginLogTitle")}
        open={loginLogOpen}
        footer={null}
        onCancel={() => setLoginLogOpen(false)}
        width={860}
        destroyOnClose
      >
        <Table
          rowKey="id"
          size="small"
          loading={myLogsLoading}
          dataSource={myLogs}
          scroll={{ x: "max-content" }}
          pagination={{
            current: myLogsPage,
            pageSize: 10,
            total: myLogsTotal,
            onChange: (p) => openMyLoginLog(p),
          }}
        >
          <Table.Column
            dataIndex="createdAt"
            title={t("loginlog.createdAt")}
            width={170}
            render={(v: string) => (v ? new Date(v).toLocaleString() : "")}
          />
          <Table.Column dataIndex="ip" title={t("loginlog.ip")} width={140} />
          <Table.Column dataIndex="location" title={t("loginlog.location")} width={120} />
          <Table.Column dataIndex="browser" title={t("loginlog.browser")} width={110} />
          <Table.Column dataIndex="os" title={t("loginlog.os")} width={110} />
          <Table.Column
            dataIndex="status"
            title={t("loginlog.status")}
            width={90}
            render={(v: number) =>
              v === 0 ? <Tag color="green">{t("loginlog.success")}</Tag> : <Tag color="red">{t("loginlog.fail")}</Tag>
            }
          />
          <Table.Column dataIndex="msg" title={t("loginlog.msg")} ellipsis />
        </Table>
      </Modal>

      {/* 修改密码 */}
      <Modal
        title={t("layout.passwordTitle")}
        open={pwdOpen}
        onOk={savePassword}
        onCancel={() => setPwdOpen(false)}
        okText={t("layout.save")}
        cancelText={t("layout.cancel")}
        destroyOnClose
      >
        <Form form={pwdForm} layout="vertical" preserve={false}>
          <Form.Item
            label={t("layout.oldPassword")}
            name="oldPassword"
            rules={[{ required: true, message: t("layout.required") }]}
          >
            <Input.Password autoComplete="current-password" />
          </Form.Item>
          <Form.Item
            label={t("layout.newPassword")}
            name="newPassword"
            rules={[
              { required: true, message: t("layout.required") },
              { min: 6, message: t("layout.pwdTooShort") },
            ]}
          >
            <Input.Password autoComplete="new-password" />
          </Form.Item>
          <Form.Item
            label={t("layout.confirmPassword")}
            name="confirmPassword"
            dependencies={["newPassword"]}
            rules={[
              { required: true, message: t("layout.required") },
              ({ getFieldValue }) => ({
                validator: (_, value) =>
                  !value || getFieldValue("newPassword") === value
                    ? Promise.resolve()
                    : Promise.reject(new Error(t("layout.pwdMismatch"))),
              }),
            ]}
          >
            <Input.Password autoComplete="new-password" />
          </Form.Item>
        </Form>
      </Modal>
    </Layout>
  );
}

// TabsView is imported lazily to keep this file focused.
import TabsView from "./TabsView";
function TabsViewMount() {
  return <TabsView />;
}

function buildMenu(items: typeof navItems, t: (key: string, fallback: string) => string) {
  const groups = new Map<string, { label: string; key: string; children: any[] }>();
  const roots: any[] = [];
  items.forEach((it) => {
    // 一级菜单(无分组)与叶子菜单都带 icon;折叠后 Sider 只展示 icon.
    const mi = { key: it.path, label: t(it.titleKey, it.title), icon: iconOf(it.icon) };
    if (it.group) {
      const gk = it.groupKey ?? it.group;
      if (!groups.has(gk)) groups.set(gk, { label: t(gk, it.group), key: `g-${gk}`, children: [] });
      groups.get(gk)!.children.push(mi);
    } else {
      roots.push(mi);
    }
  });
  const out: any[] = [...roots];
  groups.forEach(({ label, children, key }) =>
    out.push({ key, label, icon: groupIcon(key.slice(2)), children })
  );
  return out;
}
