import React, { useEffect, useState } from "react";
import { Card, Row, Col, Statistic, Spin } from "antd";
import { dashboardApi, unwrap } from "../api/client";

// Dashboard pulls headline counts from the hand-written /dashboard/stats endpoint.
export default function Dashboard() {
  const [data, setData] = useState<any>(null);

  useEffect(() => {
    unwrap<any>(dashboardApi.dashboardStatsGet())
      .then((d) => setData(d))
      .catch(() => {});
  }, []);

  if (!data) return <Spin style={{ margin: 40 }} />;

  const c = data.counts ?? {};
  return (
    <Card title="仪表盘">
      <Row gutter={16}>
        <Col span={4}>
          <Card>
            <Statistic title="用户" value={c.user ?? 0} />
          </Card>
        </Col>
        <Col span={4}>
          <Card>
            <Statistic title="角色" value={c.role ?? 0} />
          </Card>
        </Col>
        <Col span={4}>
          <Card>
            <Statistic title="商品" value={c.product ?? 0} />
          </Card>
        </Col>
        <Col span={4}>
          <Card>
            <Statistic title="订单" value={c.order ?? 0} />
          </Card>
        </Col>
        <Col span={4}>
          <Card>
            <Statistic title="字典" value={c.dict ?? 0} />
          </Card>
        </Col>
        <Col span={4}>
          <Card>
            <Statistic title="营收" value={data.revenue ?? 0} precision={2} prefix="¥" />
          </Card>
        </Col>
      </Row>
    </Card>
  );
}
