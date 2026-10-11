import React, { useEffect, useState } from "react";
import { useTable, DeleteButton } from "@refinedev/antd";
import { Table, Card, Row, Col, Statistic, Button, Space, Tag, message } from "antd";
import { orderApi, unwrap } from "../api/client";

const statusColor: Record<string, string> = {
  pending: "gold",
  paid: "green",
  shipped: "blue",
  done: "cyan",
  cancel: "red",
};

// OrderPage is a hand-written (regular mode) page. It shows the order list plus
// a revenue summary fetched from a bespoke backend endpoint, and a custom
// "mark paid" action that calls a hand-written backend route. This is the kind
// of logic Actions mode does not cover.
export default function OrderPage() {
  const { tableProps, tableQueryResult } = useTable({
    resource: "order",
    pagination: { pageSize: 10 },
    syncWithLocation: false,
  });
  const [summary, setSummary] = useState({ totalAmount: 0, paidAmount: 0, count: 0 });

  useEffect(() => {
    unwrap<any>(orderApi.orderGet({ page: 1, pageSize: 1 }))
      .then((d) => setSummary(d?.summary ?? { totalAmount: 0, paidAmount: 0, count: 0 }))
      .catch(() => {});
  }, []);

  const markPaid = async (id: any) => {
    try {
      await unwrap(orderApi.orderIdPaidPut({ id: String(id) }));
      message.success("已标记为付款");
      tableQueryResult.refetch();
    } catch (e: any) {
      message.error(e?.message || "操作失败");
    }
  };

  return (
    <Card title="订单管理（常规模式 · 手写复杂逻辑）">
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={8}>
          <Card size="small">
            <Statistic title="订单总数" value={summary.count} />
          </Card>
        </Col>
        <Col span={8}>
          <Card size="small">
            <Statistic title="订单总额" value={summary.totalAmount} precision={2} prefix="¥" />
          </Card>
        </Col>
        <Col span={8}>
          <Card size="small">
            <Statistic title="已收款" value={summary.paidAmount} precision={2} prefix="¥" />
          </Card>
        </Col>
      </Row>
      <Table {...tableProps} rowKey="id" dataSource={tableProps.dataSource as any}>
        <Table.Column dataIndex="id" title="ID" width={70} />
        <Table.Column dataIndex="orderNo" title="订单号" />
        <Table.Column dataIndex="amount" title="金额" render={(v: number) => `¥${v}`} />
        <Table.Column
          dataIndex="status"
          title="状态"
          render={(v: string) => <Tag color={statusColor[v]}>{v}</Tag>}
        />
        <Table.Column dataIndex="channel" title="渠道" />
        <Table.Column
          dataIndex="createdAt"
          title="创建时间"
          render={(v: string) => (v ? new Date(v).toLocaleString() : "")}
        />
        <Table.Column
          title="操作"
          width={160}
          render={(_: any, r: any) => (
            <Space>
              <Button
                size="small"
                type="link"
                disabled={r.status !== "pending" && r.status !== "cancel"}
                onClick={() => markPaid(r.id)}
              >
                标记付款
              </Button>
              <DeleteButton resource="order" hideText size="small" recordItemId={r.id} />
            </Space>
          )}
        />
      </Table>
    </Card>
  );
}
