import React, { useCallback, useEffect, useState } from "react";
import { useLogin } from "@refinedev/core";
import { Card, Form, Input, Button, Typography } from "antd";
import { fetchCaptcha } from "../api/client";

interface Captcha {
  id: string;
  img: string;
}

// Login page with a captcha, mirroring go-admin:
//   GET /api/v1/captcha -> { code, data: <base64 image>, id }
//   POST /api/v1/login  -> { username, password, code, uuid } (via generated AuthApi)
export default function Login() {
  const { mutate: login } = useLogin();
  const [form] = Form.useForm();
  const [captcha, setCaptcha] = useState<Captcha | null>(null);

  const refreshCaptcha = useCallback(() => {
    fetchCaptcha()
      .then((c) => setCaptcha(c))
      .catch(() => {});
  }, []);

  useEffect(() => {
    refreshCaptcha();
  }, [refreshCaptcha]);

  const onFinish = (values: any) => {
    login({ ...values, uuid: captcha?.id ?? "" });
    // Refresh the captcha either way: one captcha can only be used once,
    // and a failed attempt needs a fresh image.
    refreshCaptcha();
  };

  return (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        height: "100vh",
        background: "#f0f2f5",
      }}
    >
      <Card style={{ width: 360 }} title={<Typography.Title level={4}>Enterprise Admin</Typography.Title>}>
        <Form form={form} layout="vertical" onFinish={onFinish}>
          <Form.Item name="username" label="用户名" rules={[{ required: true, message: "请输入用户名" }]}>
            <Input placeholder="admin" />
          </Form.Item>
          <Form.Item name="password" label="密码" rules={[{ required: true, message: "请输入密码" }]}>
            <Input.Password placeholder="admin123" />
          </Form.Item>
          <Form.Item label="验证码" required style={{ marginBottom: 24 }}>
            <Input.Group compact style={{ display: "flex" }}>
              <Form.Item
                name="code"
                noStyle
                rules={[{ required: true, message: "请输入验证码" }]}
              >
                <Input placeholder="4 位数字" maxLength={4} style={{ flex: 1 }} />
              </Form.Item>
              <img
                src={captcha?.img}
                alt="验证码"
                title="点击刷新"
                onClick={refreshCaptcha}
                style={{
                  height: 40,
                  marginLeft: 8,
                  borderRadius: 6,
                  cursor: "pointer",
                  border: "1px solid #d9d9d9",
                }}
              />
            </Input.Group>
          </Form.Item>
          <Button type="primary" htmlType="submit" block>
            登录
          </Button>
        </Form>
      </Card>
    </div>
  );
}
