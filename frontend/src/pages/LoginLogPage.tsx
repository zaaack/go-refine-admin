import React, { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Table, Form, Input, Select, Button, Space, Card, Tag, message } from "antd";
import { apiAxios, unwrap } from "../api/client";
import { backendErrorMessage } from "../i18n";

const BASE = "/api/v1/sys/login-log";

// LoginLogPage shows every login attempt (success/fail) with ip/device/location.
// Same hand-written shape as OperLogPage: list + delete/clean, no create/edit.
export default function LoginLogPage() {
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
      if (filters.ip) params.ip = filters.ip;
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
  }, []);

  useEffect(() => {
    fetchList(page, pageSize, form.getFieldsValue());
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [page, pageSize]);

  const onSearch = (v: any) => {
    setPage(1);
    fetchList(1, pageSize, v);
  };

  return (
    <Card title={t("loginlog.title")}>
      <Form form={form} layout="inline" onFinish={onSearch} style={{ marginBottom: 16, rowGap: 8 }}>
        <Form.Item name="username" label={t("loginlog.username")}>
          <Input allowClear style={{ width: 160 }} />
        </Form.Item>
        <Form.Item name="ip" label={t("loginlog.ip")}>
          <Input allowClear style={{ width: 160 }} />
        </Form.Item>
        <Form.Item name="status" label={t("loginlog.status")}>
          <Select
            allowClear
            style={{ width: 140 }}
            options={[
              { label: t("loginlog.success"), value: 0 },
              { label: t("loginlog.fail"), value: 1 },
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
        <Table.Column dataIndex="username" title={t("loginlog.username")} width={130} />
        <Table.Column dataIndex="ip" title={t("loginlog.ip")} width={140} />
        <Table.Column dataIndex="location" title={t("loginlog.location")} width={130} ellipsis />
        <Table.Column dataIndex="browser" title={t("loginlog.browser")} width={120} />
        <Table.Column dataIndex="os" title={t("loginlog.os")} width={120} />
        <Table.Column dataIndex="device" title={t("loginlog.device")} width={110} />
        <Table.Column
          dataIndex="status"
          title={t("loginlog.status")}
          width={90}
          render={(v: number) =>
            v === 0 ? <Tag color="green">{t("loginlog.success")}</Tag> : <Tag color="red">{t("loginlog.fail")}</Tag>
          }
        />
        <Table.Column dataIndex="msg" title={t("loginlog.msg")} ellipsis />
        <Table.Column
          dataIndex="createdAt"
          title={t("loginlog.createdAt")}
          width={170}
          render={(v: string) => (v ? new Date(v).toLocaleString() : "")}
        />
      </Table>
    </Card>
  );
}
