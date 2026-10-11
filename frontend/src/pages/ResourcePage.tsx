import React, { useEffect, useMemo, useState } from "react";
import { usePermissions } from "@refinedev/core";
import { useTable, useForm, CreateButton, EditButton, DeleteButton } from "@refinedev/antd";
import { Table, Form, Input, InputNumber, Select, Switch, Drawer, Button, Space, Card } from "antd";
import MenuTreeSelect from "../components/MenuTreeSelect";
import { can } from "../utils/perm";
import type { FieldMeta, ResourceMeta } from "../resources/meta";

// ResourcePage is the auto-generated CRUD UI driven entirely by ResourceMeta.
// It is the front-end half of the backend's Actions mode: one component renders
// the list, search and create/edit form for any resource with no per-resource code.
export default function ResourcePage({ meta }: { meta: ResourceMeta }) {
  const [drawer, setDrawer] = useState<{ open: boolean; id?: string }>({ open: false });
  const [searchForm] = Form.useForm();
  const { data: perms } = usePermissions<string[]>();
  const pre = meta.permPrefix;
  // Write buttons follow the role's F-node permissions (read stays menu-driven).
  const canCreate = !pre || can(perms, `${pre}:create`);
  const canEdit = !pre || can(perms, `${pre}:edit`);
  const canDelete = !pre || can(perms, `${pre}:delete`);

  const { tableProps, tableQueryResult, setFilters } = useTable({
    resource: meta.name,
    pagination: { current: 1, pageSize: 10 },
    syncWithLocation: false,
  });

  const { formProps, saveButtonProps, mutationResult } = useForm({
    resource: meta.name,
    action: drawer.id ? "edit" : "create",
    id: drawer.id,
    redirect: false,
    queryOptions: { enabled: !!drawer.id },
  });

  // Close the drawer and refresh the table after a successful save.
  useEffect(() => {
    if (mutationResult?.isSuccess) {
      setDrawer({ open: false });
      tableQueryResult.refetch();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [mutationResult?.isSuccess]);

  const columns = useMemo(
    () =>
      meta.fields
        .filter((f) => !f.hideInTable)
        .map((f) => ({
          title: f.label,
          dataIndex: f.name,
          width: f.width,
          sorter: f.sortable || undefined,
          render: (val: any) => renderCell(f, val),
        })),
    [meta.fields]
  );

  const onSearch = (values: any) => {
    const filters = meta.fields
      .filter((f) => f.searchable)
      .map((f) => {
        const v = values[f.name];
        if (v === undefined || v === null || v === "") return null;
        const operator = f.type === "text" ? "contains" : "eq";
        return { field: f.name, operator, value: v };
      })
      .filter(Boolean);
    setFilters(filters as any, "replace");
  };

  return (
    <Card
      title={meta.title}
      extra={
        canCreate ? (
          <CreateButton resource={meta.name} onClick={() => setDrawer({ open: true })}>
            新建
          </CreateButton>
        ) : undefined
      }
    >
      <Form form={searchForm} layout="inline" onFinish={onSearch} style={{ marginBottom: 16, rowGap: 8 }}>
        {meta.fields
          .filter((f) => f.searchable)
          .map((f) => (
            <Form.Item key={f.name} name={f.name} label={f.label}>
              {renderSearchInput(f)}
            </Form.Item>
          ))}
        <Form.Item>
          <Space>
            <Button type="primary" htmlType="submit">
              搜索
            </Button>
            <Button
              onClick={() => {
                searchForm.resetFields();
                onSearch({});
              }}
            >
              重置
            </Button>
          </Space>
        </Form.Item>
      </Form>

      <Table {...tableProps} rowKey="id" scroll={{ x: "max-content" }} dataSource={tableProps.dataSource as any}>
        {columns.map((c) => (
          <Table.Column key={c.dataIndex as string} {...c} />
        ))}
        <Table.Column
          title="操作"
          fixed="right"
          width={140}
          render={(_: any, record: any) => (
            <Space>
              {canEdit && (
                <EditButton
                  resource={meta.name}
                  hideText
                  size="small"
                  recordItemId={record.id}
                  onClick={() => setDrawer({ open: true, id: String(record.id) })}
                />
              )}
              {canDelete && (
                <DeleteButton
                  resource={meta.name}
                  hideText
                  size="small"
                  recordItemId={record.id}
                />
              )}
            </Space>
          )}
        />
      </Table>

      <Drawer
        title={drawer.id ? `编辑${meta.title}` : `新建${meta.title}`}
        width={520}
        open={drawer.open}
        onClose={() => setDrawer({ open: false })}
        extra={
          <Button {...saveButtonProps} type="primary">
            保存
          </Button>
        }
      >
        <Form {...formProps} layout="vertical">
          {meta.fields
            .filter((f) => !f.hideInForm)
            .map((f) => (
              <Form.Item
                key={f.name}
                name={f.name}
                label={f.label}
                valuePropName={f.type === "switch" ? "checked" : "value"}
                rules={f.required ? [{ required: true, message: `请输入${f.label}` }] : []}
              >
                {renderFormInput(f)}
              </Form.Item>
            ))}
        </Form>
      </Drawer>
    </Card>
  );
}

function renderCell(f: FieldMeta, val: any) {
  if (f.type === "switch") {
    return <Switch checked={!!val} disabled />;
  }
  if (f.type === "menutree") {
    const count = String(val ?? "")
      .split(",")
      .filter(Boolean).length;
    return count ? `${count} 项菜单` : "-";
  }
  if (f.type === "select" && f.options) {
    return f.options.find((o) => o.value === val)?.label ?? val;
  }
  if (f.type === "datetime" && val) {
    return new Date(val).toLocaleString();
  }
  return val;
}

function renderSearchInput(f: FieldMeta) {
  if (f.type === "select" && f.options) {
    return <Select allowClear options={f.options} style={{ width: 160 }} />;
  }
  if (f.type === "number") {
    return <InputNumber style={{ width: 160 }} />;
  }
  return <Input allowClear style={{ width: 160 }} />;
}

function renderFormInput(f: FieldMeta) {
  switch (f.type) {
    case "number":
      return <InputNumber style={{ width: "100%" }} />;
    case "textarea":
      return <Input.TextArea rows={3} />;
    case "password":
      return <Input.Password />;
    case "switch":
      return <Switch />;
    case "select":
      return <Select options={f.options} allowClear />;
    case "menutree":
      return <MenuTreeSelect />;
    default:
      return <Input />;
  }
}
