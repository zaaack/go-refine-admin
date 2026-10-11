import React, { useCallback, useEffect, useState } from "react";
import { useLogin } from "@refinedev/core";
import { useTranslation } from "react-i18next";
import { Card, Form, Input, Button, Typography, message } from "antd";
import { fetchCaptcha } from "../api/client";
import { backendErrorMessage } from "../i18n";

interface Captcha {
  id: string;
  img: string;
}

// Login page with a captcha, mirroring go-admin:
//   GET /api/v1/captcha -> { code, data: <base64 image>, id }
//   POST /api/v1/login  -> { username, password, code, uuid } (via generated AuthApi)
export default function Login() {
  const { mutate: login } = useLogin();
  const { t } = useTranslation();
  const [form] = Form.useForm();
  const [captcha, setCaptcha] = useState<Captcha | null>(null);

  const refreshCaptcha = useCallback(() => {
    fetchCaptcha()
      .then((c) => setCaptcha(c))
      .catch((e) => message.error(backendErrorMessage(e)));
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
      <Card style={{ width: 360 }} title={<Typography.Title level={4}>{t("login.title")}</Typography.Title>}>
        <Form form={form} layout="vertical" onFinish={onFinish}>
          <Form.Item name="username" label={t("login.username")} rules={[{ required: true, message: t("login.usernameRequired") }]}>
            <Input placeholder={t("login.usernamePh")} />
          </Form.Item>
          <Form.Item name="password" label={t("login.password")} rules={[{ required: true, message: t("login.passwordRequired") }]}>
            <Input.Password placeholder={t("login.passwordPh")} />
          </Form.Item>
          <Form.Item label={t("login.captcha")} required style={{ marginBottom: 24 }}>
            <Input.Group compact style={{ display: "flex" }}>
              <Form.Item
                name="code"
                noStyle
                rules={[{ required: true, message: t("login.captchaRequired") }]}
              >
                <Input placeholder={t("login.captchaPh")} maxLength={4} style={{ flex: 1 }} />
              </Form.Item>
              <img
                src={captcha?.img}
                alt={t("login.captchaAlt")}
                title={t("login.captchaTitle")}
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
            {t("login.submit")}
          </Button>
        </Form>
      </Card>
    </div>
  );
}
