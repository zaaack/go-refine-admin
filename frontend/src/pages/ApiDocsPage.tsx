import React from "react";
import { useTranslation } from "react-i18next";
import { Card, Button, Space } from "antd";

// ApiDocsPage embeds the backend swagger UI through the Vite dev proxy
// (/swagger -> http://localhost:8080) so docs work both in dev and behind
// the Go static server (same origin, no CORS). No new deps: plain iframe.
export default function ApiDocsPage() {
  const { t } = useTranslation();
  return (
    <Card
      title={t("nav.apiDocs", { defaultValue: "接口文档" })}
      extra={
        <Space>
          <Button href="/swagger/doc.json" target="_blank">
            doc.json
          </Button>
          <Button type="primary" href="/swagger/index.html" target="_blank">
            {t("apidocs.openNewTab", { defaultValue: "新窗口打开" })}
          </Button>
        </Space>
      }
      styles={{ body: { padding: 0, height: "calc(100vh - 220px)", minHeight: 480 } }}
    >
      <iframe
        title="swagger-ui"
        src="/swagger/index.html"
        style={{ width: "100%", height: "100%", border: 0 }}
      />
    </Card>
  );
}
