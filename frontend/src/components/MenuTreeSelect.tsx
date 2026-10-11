import React, { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Radio, Space, Tag, Typography } from "antd";
import { authApi, unwrap } from "../api/client";

type Mode = "none" | "read" | "write";

interface MenuNode {
  id: number;
  name?: string; // stable identifier; UI renders t(`menu.${name}`) with title fallback
  title: string;
  path?: string;
  type?: string; // M directory | C menu | F button
  children?: MenuNode[];
}

// MenuTreeSelect is the role's permission editor with three states per menu:
//   无权限   -> menu invisible
//   只读     -> menu visible, no write buttons
//   可读可写 -> menu visible + its F (button) nodes granted
//
// Storage stays the existing sys_role.menu_ids csv: C node id = read, C + F
// child ids = write. The value prop is the single source of truth (derived on
// every render), so edit-mode initialValues load in without sync issues.
export default function MenuTreeSelect({
  value,
  onChange,
}: {
  value?: string;
  onChange?: (v: string) => void;
}) {
  const { t } = useTranslation();
  const [tree, setTree] = useState<MenuNode[]>([]);

  useEffect(() => {
    unwrap<MenuNode[]>(authApi.menuTreeGet())
      .then((d) => setTree(d ?? []))
      .catch(() => {});
  }, []);

  const nodeTitle = (n: MenuNode) =>
    n.name ? t(`menu.${n.name}`, { defaultValue: n.title }) : n.title;

  const modeOptions = [
    { label: t("menutree.none"), value: "none" as Mode },
    { label: t("menutree.read"), value: "read" as Mode },
    { label: t("menutree.write"), value: "write" as Mode },
  ];

  const granted = new Set(
    (value ?? "")
      .split(",")
      .map((s) => Number(s.trim()))
      .filter((n) => n > 0)
  );

  const fKids = (n: MenuNode) => (n.children ?? []).filter((c) => c.type === "F");
  const subMenus = (n: MenuNode) => (n.children ?? []).filter((c) => c.type !== "F");

  const stateOf = (n: MenuNode): Mode => {
    if (fKids(n).some((f) => granted.has(f.id))) return "write";
    return granted.has(n.id) ? "read" : "none";
  };

  const setMode = (target: MenuNode, mode: Mode) => {
    const states: Record<number, Mode> = {};
    const collect = (nodes: MenuNode[]) =>
      nodes.forEach((n) => {
        if (n.type !== "F") {
          states[n.id] = stateOf(n);
          collect(subMenus(n));
        }
      });
    collect(tree);
    states[target.id] = mode;
    emit(states);
  };

  // Serialize tri-states back into the menu_ids csv (ancestor M ids included).
  const emit = (states: Record<number, Mode>) => {
    const ids = new Set<number>();
    const walk = (nodes: MenuNode[], ancestors: number[]) =>
      nodes.forEach((n) => {
        if (n.type === "F") return;
        const st = states[n.id] ?? "none";
        if (st !== "none") {
          ancestors.forEach((a) => ids.add(a));
          ids.add(n.id);
          if (st === "write") fKids(n).forEach((f) => ids.add(f.id));
        }
        walk(subMenus(n), n.type === "M" ? [...ancestors, n.id] : ancestors);
      });
    walk(tree, []);
    onChange?.([...ids].join(","));
  };

  const renderNode = (n: MenuNode, depth: number): React.ReactNode => {
    if (n.type === "F") return null; // buttons render as tags under their menu
    const st = stateOf(n);
    const grantedSub = subMenus(n).some((c) => stateOf(c) !== "none");
    return (
      <div key={n.id} style={{ paddingLeft: depth * 28, marginBottom: 8 }}>
        <Space size={8} wrap>
          <span style={{ minWidth: 88, display: "inline-block" }}>{nodeTitle(n)}</span>
          {n.type === "C" && (
            <Radio.Group
              size="small"
              optionType="button"
              buttonStyle="solid"
              options={modeOptions}
              value={st}
              onChange={(e) => setMode(n, e.target.value as Mode)}
            />
          )}
          {n.type === "M" && st !== "none" && <Tag color="blue">{t("menutree.directory")}</Tag>}
          {n.type === "M" && st === "none" && grantedSub && (
            <Tag color="geekblue">{t("menutree.hasSubmenu")}</Tag>
          )}
          {st === "write" &&
            fKids(n).map((f) => (
              <Tag key={f.id} color="green">
                {nodeTitle(f)}
              </Tag>
            ))}
        </Space>
        {subMenus(n).map((c) => renderNode(c, depth + 1))}
      </div>
    );
  };

  return (
    <div
      style={{
        maxHeight: 360,
        overflow: "auto",
        border: "1px solid #f0f0f0",
        borderRadius: 6,
        padding: 12,
      }}
    >
      <Typography.Text type="secondary" style={{ display: "block", marginBottom: 8 }}>
        {t("menutree.hint")}
      </Typography.Text>
      {tree.length ? tree.map((n) => renderNode(n, 0)) : t("menutree.loading")}
    </div>
  );
}
