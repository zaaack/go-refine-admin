import React from "react";
import { Tabs } from "antd";
import { useLocation, useNavigate } from "react-router-dom";
import { useTabs } from "../context/TabsContext";
import { findNav } from "../registry";

// TabsView renders every opened tab's page component, keeping them all mounted
// and toggling visibility with CSS `display`. This is the React equivalent of
// Vue's <keep-alive>: switching tabs preserves each page's table state, filters
// and form input instead of unmounting.
export default function TabsView() {
  const { opened, active, closeTab, setActive } = useTabs();
  const navigate = useNavigate();
  const location = useLocation();

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
      const remaining = opened.filter((t) => t.key !== key);
      if (remaining.length) navigate(remaining[remaining.length - 1].key);
      else navigate("/dashboard");
    }
  };

  return (
    <div style={{ display: "flex", flexDirection: "column", height: "100%" }}>
      <Tabs
        type="editable-card"
        hideAdd
        activeKey={active}
        onChange={onTabChange}
        onEdit={onTabEdit as any}
        items={opened.map((t) => ({
          key: t.key,
          label: t.title,
          closable: t.closable !== false,
        }))}
        style={{ marginBottom: 0 }}
      />
      <div style={{ flex: 1, overflow: "auto", padding: 16, position: "relative" }}>
        {opened.map((t) => {
          const nav = findNav(t.key);
          const visible = t.key === active;
          return (
            <div
              key={t.key}
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
