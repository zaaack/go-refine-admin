import React, { createContext, useCallback, useContext, useState } from "react";

export interface TabItem {
  key: string; // route path
  title: string;
  closable?: boolean;
}

interface TabsContextValue {
  opened: TabItem[];
  active: string;
  openTab: (tab: TabItem) => void;
  closeTab: (key: string) => void;
  setActive: (key: string) => void;
}

const TabsContext = createContext<TabsContextValue | null>(null);

// TabsProvider holds the set of open tabs for the multi-tab admin layout.
// Tabs persist across navigation (the layout keeps every opened tab mounted and
// toggles visibility), which preserves each page's table/filter state.
export function TabsProvider({ children }: { children: React.ReactNode }) {
  const [opened, setOpened] = useState<TabItem[]>([
    { key: "/dashboard", title: "仪表盘", closable: false },
  ]);
  const [active, setActive] = useState<string>("/dashboard");

  const openTab = useCallback((tab: TabItem) => {
    setOpened((prev) => (prev.some((t) => t.key === tab.key) ? prev : [...prev, tab]));
    setActive(tab.key);
  }, []);

  const closeTab = useCallback((key: string) => {
    setOpened((prev) => prev.filter((t) => t.key !== key));
  }, []);

  return (
    <TabsContext.Provider value={{ opened, active, openTab, closeTab, setActive }}>
      {children}
    </TabsContext.Provider>
  );
}

export function useTabs(): TabsContextValue {
  const ctx = useContext(TabsContext);
  if (!ctx) throw new Error("useTabs must be used within TabsProvider");
  return ctx;
}
