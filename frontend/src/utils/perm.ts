// Permission helpers. Perms come from the login response (the role's granted
// F-menu nodes): ["sys:user:create", ...] or ["*"] for admin.

export function loadPerms(): string[] {
  try {
    return JSON.parse(localStorage.getItem("perms") || "[]");
  } catch {
    return [];
  }
}

// can is fail-open while the permission list is still loading (undefined).
export function can(perms: string[] | undefined, perm: string): boolean {
  if (!perms) return true;
  return perms.includes("*") || perms.includes(perm);
}
