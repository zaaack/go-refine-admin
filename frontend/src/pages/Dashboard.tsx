import React, { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Card, Row, Col, Statistic, Spin } from "antd";
import { dashboardApi, unwrap } from "../api/client";

// Dashboard pulls headline counts from the hand-written /dashboard/stats endpoint.
export default function Dashboard() {
  const { t } = useTranslation();
  const [data, setData] = useState<any>(null);

  useEffect(() => {
    unwrap<any>(dashboardApi.dashboardStatsGet())
      .then((d) => setData(d))
      .catch(() => {});
  }, []);

  if (!data) return <Spin style={{ margin: 40 }} />;

  const c = data.counts ?? {};
  return (
    <Card title={t("dashboard.title")}>
      <Row gutter={16}>
        <Col span={4}>
          <Card>
            <Statistic title={t("dashboard.user")} value={c.user ?? 0} />
          </Card>
        </Col>
        <Col span={4}>
          <Card>
            <Statistic title={t("dashboard.role")} value={c.role ?? 0} />
          </Card>
        </Col>
        <Col span={4}>
          <Card>
            <Statistic title={t("dashboard.product")} value={c.product ?? 0} />
          </Card>
        </Col>
        <Col span={4}>
          <Card>
            <Statistic title={t("dashboard.order")} value={c.order ?? 0} />
          </Card>
        </Col>
        <Col span={4}>
          <Card>
            <Statistic title={t("dashboard.dict")} value={c.dict ?? 0} />
          </Card>
        </Col>
        <Col span={4}>
          <Card>
            <Statistic title={t("dashboard.revenue")} value={data.revenue ?? 0} precision={2} prefix="¥" />
          </Card>
        </Col>
      </Row>
    </Card>
  );
}
