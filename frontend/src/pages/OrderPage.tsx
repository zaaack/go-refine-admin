import React, { useEffect, useState } from "react";
import { useTable, DeleteButton } from "@refinedev/antd";
import { useTranslation } from "react-i18next";
import { Table, Card, Row, Col, Statistic, Button, Space, Tag, message } from "antd";
import { orderApi, unwrap } from "../api/client";
import { backendErrorMessage } from "../i18n";

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
  const { t } = useTranslation();
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
      message.success(t("order.markedPaid"));
      tableQueryResult.refetch();
    } catch (e: any) {
      message.error(backendErrorMessage(e));
    }
  };

  return (
    <Card title={t("order.title")}>
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={8}>
          <Card size="small">
            <Statistic title={t("order.total")} value={summary.count} />
          </Card>
        </Col>
        <Col span={8}>
          <Card size="small">
            <Statistic title={t("order.totalAmount")} value={summary.totalAmount} precision={2} prefix="¥" />
          </Card>
        </Col>
        <Col span={8}>
          <Card size="small">
            <Statistic title={t("order.paidAmount")} value={summary.paidAmount} precision={2} prefix="¥" />
          </Card>
        </Col>
      </Row>
      <Table {...tableProps} rowKey="id" dataSource={tableProps.dataSource as any}>
        <Table.Column dataIndex="id" title={t("order.id")} width={70} />
        <Table.Column dataIndex="orderNo" title={t("order.no")} />
        <Table.Column dataIndex="amount" title={t("order.amount")} render={(v: number) => `¥${v}`} />
        <Table.Column
          dataIndex="status"
          title={t("order.status")}
          render={(v: string) => <Tag color={statusColor[v]}>{t(`order.statusValue.${v}`, { defaultValue: v })}</Tag>}
        />
        <Table.Column dataIndex="channel" title={t("order.channel")} />
        <Table.Column
          dataIndex="createdAt"
          title={t("order.createdAt")}
          render={(v: string) => (v ? new Date(v).toLocaleString() : "")}
        />
        <Table.Column
          title={t("order.actions")}
          width={160}
          render={(_: any, r: any) => (
            <Space>
              <Button
                size="small"
                type="link"
                disabled={r.status !== "pending" && r.status !== "cancel"}
                onClick={() => markPaid(r.id)}
              >
                {t("order.markPaid")}
              </Button>
              <DeleteButton resource="order" hideText size="small" recordItemId={r.id} />
            </Space>
          )}
        />
      </Table>
    </Card>
  );
}
