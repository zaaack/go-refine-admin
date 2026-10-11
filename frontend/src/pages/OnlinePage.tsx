import React, { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Table, Form, Input, Button, Space, Card, Tag, Popconfirm, message } from "antd";
import { apiAxios, unwrap } from "../api/client";
import { backendErrorMessage } from "../i18n";

const BASE = "/api/v1/auth/online";

// OnlinePage lists active JWT sessions from sys_user_token. Force-offline is
// admin-only on the backend (403 otherwise); the button stays visible and the
// server error surfaces for non-admin users.
export default function OnlinePage() {
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [rows, setRows] = useState<any[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);
  const [loading, setLoading] = useState(false);

  const fetchList = useCallback(async (p: number, ps: number, filters: any = {}) => {
    setLoading(true);
    try {
      const params: Record<string, any> = { page: p, pageSize: ps };
      if (filters.username) params.username__like = filters.username;
      const d = await unwrap<{ list: any[]; total: number }>(apiAxios.get(BASE, { params }));
      setRows(d.list ?? []);
      setTotal(d.total ?? 0);
    } catch (e: any) {
      message.error(backendErrorMessage(e));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchList(page, pageSize, form.getFieldsValue());
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize]);

  const onSearch = (v: any) => {
    setPage(1);
    fetchList(1, pageSize, v);
  };

  const onForceOffline = async (id: any) => {
    try {
      await unwrap(apiAxios.delete(`${BASE}/${id}`));
      message.success(t("online.offlineSuccess"));
      fetchList(page, pageSize, form.getFieldsValue());
    } catch (e: any) {
      message.error(backendErrorMessage(e));
    }
  };

  return (
    <Card
      title={t("online.title")}
      extra={<Button onClick={() => fetchList(page, pageSize, form.getFieldsValue())}>{t("online.refresh")}</Button>}
    >
      <Form form={form} layout="inline" onFinish={onSearch} style={{ marginBottom: 16, rowGap: 8 }}>
        <Form.Item name="username" label={t("online.username")}>
          <Input allowClear style={{ width: 160 }} />
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
        <Table.Column dataIndex="username" title={t("online.username")} width={130} />
        <Table.Column dataIndex="ip" title={t("online.ip")} width={140} />
        <Table.Column dataIndex="location" title={t("online.location")} width={130} ellipsis />
        <Table.Column dataIndex="browser" title={t("online.browser")} width={120} />
        <Table.Column dataIndex="os" title={t("online.os")} width={120} />
        <Table.Column dataIndex="device" title={t("online.device")} width={110} />
        <Table.Column
          dataIndex="loginTime"
          title={t("online.loginTime")}
          width={170}
          render={(v: string) => (v ? new Date(v).toLocaleString() : "")}
        />
        <Table.Column
          dataIndex="expireAt"
          title={t("online.expireAt")}
          width={170}
          render={(v: string) => (v ? new Date(v).toLocaleString() : "")}
        />
        <Table.Column
          title={t("resource.actions")}
          fixed="right"
          width={110}
          render={(_: any, r: any) => (
            <Popconfirm title={t("online.offlineConfirm")} onConfirm={() => onForceOffline(r.id)}>
              <Button size="small" type="link" danger>
                {t("online.forceOffline")}
              </Button>
            </Popconfirm>
          )}
        />
      </Table>
    </Card>
  );
}
