import React from "react";
import { Refine, Authenticated } from "@refinedev/core";
import routerProvider, { CatchAllNavigate } from "@refinedev/react-router-v6";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { ConfigProvider } from "antd";
import zhCN from "antd/locale/zh_CN";
import enUS from "antd/locale/en_US";
import { useTranslation } from "react-i18next";

import "./i18n";
import { refineI18nProvider } from "./i18n";
import { dataProvider } from "./providers/dataProvider";
import { authProvider } from "./providers/authProvider";
import { TabsProvider } from "./context/TabsContext";
import AdminLayout from "./layout/AdminLayout";
import Login from "./pages/Login";

// refine resource declarations (used for notifications / permissions). The actual
// routing is handled by our multi-tab layout, not refine's route generator.
const resources = [
  { name: "demo/product", list: "/demo/product", canDelete: true },
  { name: "order", list: "/demo/order", canDelete: true },
  { name: "sys/user", list: "/system/user", canDelete: true },
  { name: "sys/role", list: "/system/role", canDelete: true },
  { name: "sys/dept", list: "/system/dept", canDelete: true },
  { name: "sys/dict", list: "/system/dict", canDelete: true },
];

// AntdLocale follows the react-i18next language so antd built-ins
// (pagination, table, dates) switch zh-CN/en together with our UI.
function AntdLocale({ children }: { children: React.ReactNode }) {
  const { i18n } = useTranslation();
  return <ConfigProvider locale={i18n.language === "en" ? enUS : zhCN}>{children}</ConfigProvider>;
}

export default function App() {
  return (
    <BrowserRouter>
      <AntdLocale>
        <Refine
          dataProvider={dataProvider}
          authProvider={authProvider}
          routerProvider={routerProvider}
          i18nProvider={refineI18nProvider}
          resources={resources}
          options={{
            syncWithLocation: false,
            warnWhenUnsavedChanges: false,
            disableTelemetry: true,
          }}
        >
          <TabsProvider>
            <Routes>
              {/* Logged-in users are bounced away from /login. */}
              <Route
                path="/login"
                element={
                  <Authenticated key="login" fallback={<Login />}>
                    <Navigate to="/" />
                  </Authenticated>
                }
              />
              {/* Everything else requires authentication; otherwise redirect to /login. */}
              <Route
                path="/*"
                element={
                  <Authenticated key="app" fallback={<CatchAllNavigate to="/login" />}>
                    <AdminLayout />
                  </Authenticated>
                }
              />
            </Routes>
          </TabsProvider>
        </Refine>
      </AntdLocale>
    </BrowserRouter>
  );
}
