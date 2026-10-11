import React from "react";
import { Refine, Authenticated } from "@refinedev/core";
import routerProvider, { CatchAllNavigate } from "@refinedev/react-router-v6";
import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { ConfigProvider } from "antd";
import zhCN from "antd/locale/zh_CN";

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

export default function App() {
  return (
    <BrowserRouter>
      <ConfigProvider locale={zhCN}>
        <Refine
          dataProvider={dataProvider}
          authProvider={authProvider}
          routerProvider={routerProvider}
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
      </ConfigProvider>
    </BrowserRouter>
  );
}
