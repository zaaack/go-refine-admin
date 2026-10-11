import React, { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Table, Form, Input, Select, Button, Space, Card, Tag, Drawer, Descriptions, message } from "antd";
import { apiAxios, unwrap } from "../api/client";
import { backendErrorMessage } from "../i18n";

const BASE = "/api/v1/sys/oper-log";

// OperLogPage is a hand-written page for the audit log (regular-mode backend
// route, not an Actions resource). Read-only list + detail + delete/clean:
// logs are written by the server middleware, never created by hand.
export default function OperLogPage() {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [rows, setRows] = useState<any[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [loading, setLoading] = useState(false);
  const [detail, setDetail] = useState<any>(null);

  const fetchList = useCallback(
    async (p: number, ps: number, filters: any = {}) => {
      setLoading(true);
      try {
        const params: Record<string, any> = { page: p, pageSize: ps };
        if (filters.username) params.username__like = filters.username;
        if (filters.path) params.path__like = filters.path;
        if (filters.status !== undefined && filters.status !== null && filters.status !== "") {
          params.status__eq = filters.status;
        }
        const d = await unwrap<{ list: any[]; total: number }>(apiAxios.get(BASE, { params }));
        setRows(d.list ?? []);
        setTotal(d.total ?? 0);
      } catch (e: any) {
        message.error(backendErrorMessage(e));
      } finally {
        setLoading(false);
      }
    },
    []
  );

  useEffect(() => {
    fetchList(page, pageSize, form.getFieldsValue());
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize]);

  const onSearch = (v: any) => {
    setPage(1);
    fetchList(1, pageSize, v);
  };

  const onDetail = async (id: any) => {
    try {
      const d = await unwrap<any>(apiAxios.get(`${BASE}/${id}`));
      setDetail(d);
    } catch (e: any) {
      message.error(backendErrorMessage(e));
    }
  };

  return (
    <Card title={t("operlog.title")}>
      <Form form={form} layout="inline" onFinish={onSearch} style={{ marginBottom: 16, rowGap: 8 }}>
        <Form.Item name="username" label={t("operlog.username")}>
          <Input allowClear style={{ width: 160 }} />
        </Form.Item>
        <Form.Item name="path" label={t("operlog.path")}>
          <Input allowClear style={{ width: 200 }} />
        </Form.Item>
        <Form.Item name="status" label={t("operlog.status")}>
          <Select
            allowClear
            style={{ width: 140 }}
            options={[
              { label: t("operlog.success"), value: 0 },
              { label: t("operlog.fail"), value: 1 },
            ]}
          />
        </Form.Item>
        <Form.Item>
          <Space>
            <Button type="primary" htmlType="submit">
              {t("resource.search")}
            </Button>
            <Button onClick={() => { form.resetFields(); onSearch({}); }}>
              {t("resource.reset")}
            </Button>
          </Space>
        </Form.Item>
      </Form>

      <Table
        rowKey="id"
        loading={loading}
        dataSource={rows}
        scroll={{ x: "max-content" }}
        pagination={{
          current: page,
          pageSize,
          total,
          showSizeChanger: true,
          onChange: (p, ps) => { setPage(p); setPageSize(ps); },
        }}
      >
        <Table.Column dataIndex="id" title="ID" width={70} />
        <Table.Column dataIndex="title" title={t("operlog.title_")} ellipsis />
        <Table.Column dataIndex="method" title={t("operlog.method")} width={90} />
        <Table.Column dataIndex="path" title={t("operlog.path")} ellipsis />
        <Table.Column dataIndex="username" title={t("operlog.username")} width={120} />
        <Table.Column dataIndex="ip" title={t("operlog.ip")} width={140} />
        <Table.Column dataIndex="location" title={t("operlog.location")} width={120} ellipsis />
        <Table.Column
          dataIndex="status"
          title={t("operlog.status")}
          width={90}
          render={(v: number) =>
            v === 0 ? <Tag color="green">{t("operlog.success")}</Tag> : <Tag color="red">{t("operlog.fail")}</Tag>
          }
        />
        <Table.Column dataIndex="statusCode" title={t("operlog.statusCode")} width={90} />
        <Table.Column dataIndex="duration" title={t("operlog.duration")} width={90} render={(v: number) => `${v}ms`} />
        <Table.Column
          dataIndex="createdAt"
          title={t("operlog.createdAt")}
          width={170}
          render={(v: string) => (v ? new Date(v).toLocaleString() : "")}
        />
        <Table.Column
          title={t("resource.actions")}
          fixed="right"
          width={90}
          render={(_: any, r: any) => (
            <Button size="small" type="link" onClick={() => onDetail(r.id)}>
              {t("operlog.detail")}
            </Button>
          )}
        />
      </Table>

      <Drawer title={t("operlog.detail")} width={560} open={!!detail} onClose={() => setDetail(null)}>
        {detail && (
          <Descriptions column={1} bordered size="small">
            <Descriptions.Item label={t("operlog.title_")}>{detail.title}</Descriptions.Item>
            <Descriptions.Item label={t("operlog.method")}>{detail.method}</Descriptions.Item>
            <Descriptions.Item label={t("operlog.path")}>{detail.path}</Descriptions.Item>
            <Descriptions.Item label={t("operlog.username")}>{detail.username}</Descriptions.Item>
            <Descriptions.Item label={t("operlog.ip")}>{detail.ip}</Descriptions.Item>
            <Descriptions.Item label={t("operlog.location")}>{detail.location}</Descriptions.Item>
            <Descriptions.Item label="User-Agent">{detail.userAgent}</Descriptions.Item>
            <Descriptions.Item label={t("operlog.statusCode")}>{detail.statusCode}</Descriptions.Item>
            <Descriptions.Item label={t("operlog.duration")}>{`${detail.duration}ms`}</Descriptions.Item>
            <Descriptions.Item label={t("operlog.body")}>
              <pre style={{ whiteSpace: "pre-wrap", margin: 0 }}>{detail.body || "-"}</pre>
            </Descriptions.Item>
            <Descriptions.Item label={t("operlog.errMsg")}>{detail.errMsg || "-"}</Descriptions.Item>
            <Descriptions.Item label={t("operlog.createdAt")}>
              {detail.createdAt ? new Date(detail.createdAt).toLocaleString() : ""}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Drawer>
    </Card>
  );
}
