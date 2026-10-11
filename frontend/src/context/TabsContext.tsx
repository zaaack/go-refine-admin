import React, { createContext, useCallback, useContext, useMemo, useState } from "react";

export interface TabItem {
  key: string; // route path
  title: string; // fallback text; UI renders t(titleKey, { defaultValue: title })
  titleKey?: string;
  closable?: boolean;
  pinned?: boolean;
}

interface TabsContextValue {
  opened: TabItem[];
  active: string;
  openTab: (tab: TabItem) => void;
  closeTab: (key: string) => void;
  closeTabs: (keys: string[]) => void;
  togglePin: (key: string) => void;
  refreshSeq: Record<string, number>;
  refreshTab: (key: string) => void;
  maximized: boolean;
  toggleMaximized: () => void;
  setActive: (key: string) => void;
}

const TabsContext = createContext<TabsContextValue | null>(null);

// TabsProvider holds the set of open tabs for the multi-tab admin layout.
// Tabs persist across navigation (the layout keeps every opened tab mounted and
// toggles visibility), which preserves each page's table/filter state.
// Titles are stored as titleKey + fallback so switching languages re-renders
// every open tab without stale text.
// Pinned tabs stay at the front of the array and are never closable.
export function TabsProvider({ children }: { children: React.ReactNode }) {
  const [opened, setOpened] = useState<TabItem[]>([
    { key: "/dashboard", title: "仪表盘", titleKey: "nav.dashboard", closable: false },
  ]);
  const [active, setActive] = useState<string>("/dashboard");
  const [refreshSeq, setRefreshSeq] = useState<Record<string, number>>({});
  const [maximized, setMaximized] = useState<boolean>(false);

  const openTab = useCallback((tab: TabItem) => {
    setOpened((prev) => {
      const exists = prev.some((t) => t.key === tab.key);
      if (!exists) return [...prev, tab];
      // Refresh the stored titleKey/title so language switches propagate.
      return prev.map((t) => (t.key === tab.key ? { ...t, title: tab.title, titleKey: tab.titleKey ?? t.titleKey } : t));
    });
    setActive(tab.key);
  }, []);

  const closeTab = useCallback((key: string) => {
    setOpened((prev) => prev.filter((t) => t.key !== key));
  }, []);

  const closeTabs = useCallback((keys: string[]) => {
    const set = new Set(keys);
    setOpened((prev) => prev.filter((t) => !set.has(t.key)));
  }, []);

  const togglePin = useCallback((key: string) => {
    setOpened((prev) => {
      const target = prev.find((t) => t.key === key);
      if (!target) return prev;
      if (target.pinned) {
        // Unpin: move to the end of the unpinned section, keeping pinned block intact.
        const rest = prev.filter((t) => t.key !== key);
        return [...rest, { ...target, pinned: false }];
      }
      // Pin: move to the end of the pinned block at the front.
      const others = prev.filter((t) => t.key !== key);
      const pinnedCount = others.filter((t) => t.pinned).length;
      const pinned = [...others.slice(0, pinnedCount), { ...target, pinned: true }];
      const unpinned = others.slice(pinnedCount);
      return [...pinned, ...unpinned];
    });
  }, []);

  const refreshTab = useCallback((key: string) => {
    setRefreshSeq((prev) => ({ ...prev, [key]: (prev[key] ?? 0) + 1 }));
  }, []);

  const toggleMaximized = useCallback(() => setMaximized((v) => !v), []);

  const value = useMemo(
    () => ({ opened, active, openTab, closeTab, closeTabs, togglePin, refreshSeq, refreshTab, maximized, toggleMaximized, setActive }),
    [opened, active, openTab, closeTab, closeTabs, togglePin, refreshSeq, refreshTab, maximized, toggleMaximized]
  );

  return (
    <TabsContext.Provider value={value}>
      {children}
    </TabsContext.Provider>
  );
}

export function useTabs(): TabsContextValue {
  const ctx = useContext(TabsContext);
  if (!ctx) throw new Error("useTabs must be used within TabsProvider");
  return ctx;
}
