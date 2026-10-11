import React, { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Card, Row, Col, Descriptions, Progress, Button, message, Spin } from "antd";
import { apiAxios, unwrap } from "../api/client";
import { backendErrorMessage } from "../i18n";

// ServerPage renders the std-only host/runtime metrics from GET /server-monitor.
// Shape mirrors go-admin-ui monitor.vue: os/mem/cpu/disk (+swap/net stubs).
export default function ServerPage() {
  const { t } = useTranslation();
  const [info, setInfo] = useState<any>(null);
  const [loading, setLoading] = useState(false);

  const fetchInfo = useCallback(async () => {
    setLoading(true);
    try {
      const d = await unwrap<any>(apiAxios.get("/api/v1/server-monitor"));
      setInfo(d);
    } catch (e: any) {
      message.error(backendErrorMessage(e));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchInfo();
  }, [fetchInfo]);

  const mem = info?.mem ?? {};
  const cpu = info?.cpu ?? {};
  const disk = info?.disk ?? {};
  const os = info?.os ?? {};

  return (
    <Card title={t("server.title")} extra={<Button onClick={fetchInfo}>{t("server.refresh")}</Button>}>
      <Spin spinning={loading}>
        <Row gutter={16}>
          <Col span={12}>
            <Card size="small" title={t("server.os")} style={{ marginBottom: 16 }}>
              <Descriptions column={1} size="small" bordered>
                <Descriptions.Item label={t("server.hostName")}>{os.hostName}</Descriptions.Item>
                <Descriptions.Item label={t("server.goOs")}>{os.goOs} / {os.arch}</Descriptions.Item>
                <Descriptions.Item label={t("server.goVersion")}>{os.version}</Descriptions.Item>
                <Descriptions.Item label={t("server.ip")}>{os.ip}</Descriptions.Item>
                <Descriptions.Item label={t("server.bootTime")}>{info?.bootTime ?? "-"}</Descriptions.Item>
                <Descriptions.Item label={t("server.time")}>{os.time}</Descriptions.Item>
              </Descriptions>
            </Card>
          </Col>
          <Col span={12}>
            <Card size="small" title={t("server.cpu")} style={{ marginBottom: 16 }}>
              <Descriptions column={1} size="small" bordered>
                <Descriptions.Item label={t("server.cpuNum")}>{cpu.cpuNum}</Descriptions.Item>
                <Descriptions.Item label={t("server.cpuModel")}>
                  {(cpu.cpuInfo?.[0]?.modelName ?? "-") as string}
                </Descriptions.Item>
              </Descriptions>
            </Card>
            <Card size="small" title={t("server.mem")} style={{ marginBottom: 16 }}>
              <p>{`${mem.used ?? 0} MB / ${mem.total ?? 0} MB`}</p>
              <Progress percent={Number(mem.percent ?? 0)} status="active" />
            </Card>
            <Card size="small" title={t("server.disk")}>
              <p>{`${disk.used ?? 0} GB / ${disk.total ?? 0} GB`}</p>
              <Progress percent={Number(disk.percent ?? 0)} status="active" />
            </Card>
          </Col>
        </Row>
      </Spin>
    </Card>
  );
}
