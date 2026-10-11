import React from "react";
import { Dropdown, Tabs } from "antd";
import { PushpinFilled } from "@ant-design/icons";
import { useLocation, useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useTabs, type TabItem } from "../context/TabsContext";
import { findNav } from "../registry";

// TabsView renders every opened tab's page component, keeping them all mounted
// and toggling visibility with CSS `display`. This is the React equivalent of
// Vue's <keep-alive>: switching tabs preserves each page's table state, filters
// and form input instead of unmounting.
// Tab labels re-render through i18n (titleKey + fallback) on language switch.
// Right-clicking a label opens a context menu: refresh / pin / close current,
// close others, close right, close all / maximize.
export default function TabsView() {
  const { opened, active, closeTab, closeTabs, togglePin, refreshTab, refreshSeq, maximized, toggleMaximized, setActive } =
    useTabs();
  const navigate = useNavigate();
  const location = useLocation();
  const { t } = useTranslation();

  React.useEffect(() => {
    if (location.pathname !== active) setActive(location.pathname);
  }, [location.pathname, active, setActive]);

  const onTabChange = (key: string) => {
    setActive(key);
    navigate(key);
  };

  const onTabEdit = (key: any, action: "add" | "remove") => {
    if (action === "remove") {
      closeTab(key);
      const remaining = opened.filter((tab) => tab.key !== key);
      if (remaining.length) navigate(remaining[remaining.length - 1].key);
      else navigate("/dashboard");
    }
  };

  const closableOf = (tab: TabItem) => tab.closable !== false && !tab.pinned;

  // Close the given keys, then navigate if the active tab was closed.
  const closeAndNavigate = (keys: string[]) => {
    closeTabs(keys);
    const set = new Set(keys);
    if (set.has(active)) {
      const remaining = opened.filter((tab) => !set.has(tab.key));
      if (remaining.length) navigate(remaining[remaining.length - 1].key);
      else navigate("/dashboard");
    }
  };

  const contextMenuFor = (tab: TabItem) => {
    const closableTabs = opened.filter(closableOf);
    const idx = opened.findIndex((it) => it.key === tab.key);
    const rightClosable = opened.slice(idx + 1).filter(closableOf);
    const othersClosable = closableTabs.filter((it) => it.key !== tab.key);
    return {
      items: [
        { key: "refresh", label: t("tab.refresh", { defaultValue: "刷新" }) },
        {
          key: "pin",
          label: tab.pinned
            ? t("tab.unpin", { defaultValue: "取消固定" })
            : t("tab.pin", { defaultValue: "固定" }),
          disabled: tab.key === "/dashboard",
        },
        { type: "divider" as const },
        { key: "closeCurrent", label: t("tab.closeCurrent", { defaultValue: "关闭当前" }), disabled: !closableOf(tab) },
        { key: "closeOthers", label: t("tab.closeOthers", { defaultValue: "关闭其他" }), disabled: !othersClosable.length },
        { key: "closeRight", label: t("tab.closeRight", { defaultValue: "关闭右侧" }), disabled: !rightClosable.length },
        { key: "closeAll", label: t("tab.closeAll", { defaultValue: "关闭全部" }), disabled: !closableTabs.length },
        { type: "divider" as const },
        {
          key: "maximize",
          label: maximized
            ? t("tab.exitMaximize", { defaultValue: "退出最大化" })
            : t("tab.maximize", { defaultValue: "最大化" }),
        },
      ],
      onClick: ({ key, domEvent }: { key: string; domEvent: { stopPropagation: () => void } }) => {
        domEvent.stopPropagation();
        if (key === "refresh") refreshTab(tab.key);
        else if (key === "pin") togglePin(tab.key);
        else if (key === "closeCurrent") closeAndNavigate([tab.key]);
        else if (key === "closeOthers") closeAndNavigate(othersClosable.map((it) => it.key));
        else if (key === "closeRight") closeAndNavigate(rightClosable.map((it) => it.key));
        else if (key === "closeAll") closeAndNavigate(closableTabs.map((it) => it.key));
        else if (key === "maximize") toggleMaximized();
      },
    };
  };

  return (
    <div style={{ display: "flex", flexDirection: "column", height: "100%" }}>
      <Tabs
        type="editable-card"
        hideAdd
        activeKey={active}
        onChange={onTabChange}
        onEdit={onTabEdit as any}
        items={opened.map((tab) => ({
          key: tab.key,
          label: (
            <Dropdown trigger={["contextMenu"]} menu={contextMenuFor(tab)}>
              <span className="tab-label-full" onContextMenu={(e) => e.preventDefault()}>
                {tab.pinned && <PushpinFilled style={{ marginRight: 4, fontSize: 12 }} />}
                {tab.titleKey ? t(tab.titleKey, { defaultValue: tab.title }) : tab.title}
              </span>
            </Dropdown>
          ),
          closable: closableOf(tab),
        }))}
        style={{ marginBottom: 0 }}
      />
      <div style={{ flex: 1, overflow: "auto", padding: 16, position: "relative" }}>
        {opened.map((tab) => {
          const nav = findNav(tab.key);
          const visible = tab.key === active;
          return (
            <div
              key={`${tab.key}#${refreshSeq[tab.key] ?? 0}`}
              style={{ display: visible ? "block" : "none" }}
            >
              {nav?.element}
            </div>
          );
        })}
      </div>
    </div>
  );
}
