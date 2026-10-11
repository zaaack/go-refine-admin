package db

import (
	"strings"

	"go-refine-admin/backend/models"
	"gorm.io/gorm"
)

// allModels lists every table to auto-migrate.
func allModels() []interface{} {
	return []interface{}{
		&models.SysUser{},
		&models.SysRole{},
		&models.SysMenu{},
		&models.SysDept{},
		&models.SysDict{},
		&models.SysDictDetail{},
		&models.SysOperLog{},
		&models.SysLoginLog{},
		&models.SysUserToken{},
		&models.SysConfig{},
		&models.Upload{},
		&models.Product{},
		&models.Order{},
	}
}

// Migrate auto-migrates all tables and seeds initial data when empty.
func Migrate(gdb *gorm.DB) error {
	if err := gdb.AutoMigrate(allModels()...); err != nil {
		return err
	}
	return seed(gdb)
}

func seed(gdb *gorm.DB) error {
	// Admin role
	var roleCount int64
	gdb.Model(&models.SysRole{}).Count(&roleCount)
	if roleCount == 0 {
		admin := models.SysRole{Name: "超级管理员", RoleKey: "admin", Status: 1, Sort: 1, DataScope: 1, Remark: "super"}
		if err := gdb.Create(&admin).Error; err != nil {
			return err
		}
		// Menus (a small tree) granted to admin. Type: M=directory, C=menu
		// (read access), F=button (write permission). A role granted only the
		// C node is read-only; granting the F nodes makes it read-write.
		menus := []models.SysMenu{
			{Name: "Dashboard", Title: "仪表盘", Icon: "dashboard", Path: "/dashboard", Component: "dashboard", ParentID: 0, Sort: 1, Type: "C", Status: 1},
			{Name: "System", Title: "系统管理", Icon: "setting", Path: "/system", Component: "system", ParentID: 0, Sort: 2, Type: "M", Status: 1},
			{Name: "SysUser", Title: "用户管理", Icon: "user", Path: "/system/user", Component: "system/user", ParentID: 2, Sort: 1, Type: "C", Status: 1},
			{Name: "SysUserCreate", Title: "用户新增", Permission: "sys:user:create", ParentID: 3, Sort: 1, Type: "F", Status: 1},
			{Name: "SysUserEdit", Title: "用户编辑", Permission: "sys:user:edit", ParentID: 3, Sort: 2, Type: "F", Status: 1},
			{Name: "SysUserDelete", Title: "用户删除", Permission: "sys:user:delete", ParentID: 3, Sort: 3, Type: "F", Status: 1},
			{Name: "SysRole", Title: "角色管理", Icon: "role", Path: "/system/role", Component: "system/role", ParentID: 2, Sort: 2, Type: "C", Status: 1},
			{Name: "SysMenu", Title: "菜单管理", Icon: "menu", Path: "/system/menu", Component: "system/menu", ParentID: 2, Sort: 3, Type: "C", Status: 1},
			{Name: "SysDept", Title: "部门管理", Icon: "dept", Path: "/system/dept", Component: "system/dept", ParentID: 2, Sort: 4, Type: "C", Status: 1},
			{Name: "SysDict", Title: "字典管理", Icon: "dict", Path: "/system/dict", Component: "system/dict", ParentID: 2, Sort: 5, Type: "C", Status: 1},
			{Name: "Demo", Title: "业务示例", Icon: "example", Path: "/demo", Component: "demo", ParentID: 0, Sort: 3, Type: "M", Status: 1},
			{Name: "Product", Title: "商品(Actions)", Icon: "goods", Path: "/demo/product", Component: "demo/product", ParentID: 11, Sort: 1, Type: "C", Status: 1},
			{Name: "ProductCreate", Title: "商品新增", Permission: "demo:product:create", ParentID: 12, Sort: 1, Type: "F", Status: 1},
			{Name: "ProductEdit", Title: "商品编辑", Permission: "demo:product:edit", ParentID: 12, Sort: 2, Type: "F", Status: 1},
			{Name: "ProductDelete", Title: "商品删除", Permission: "demo:product:delete", ParentID: 12, Sort: 3, Type: "F", Status: 1},
			{Name: "Order", Title: "订单(常规)", Icon: "order", Path: "/demo/order", Component: "demo/order", ParentID: 11, Sort: 2, Type: "C", Status: 1},
			{Name: "OrderCreate", Title: "订单新增", Permission: "demo:order:create", ParentID: 16, Sort: 1, Type: "F", Status: 1},
			{Name: "OrderPaid", Title: "订单付款", Permission: "demo:order:paid", ParentID: 16, Sort: 2, Type: "F", Status: 1},
			{Name: "OrderDelete", Title: "订单删除", Permission: "demo:order:delete", ParentID: 16, Sort: 3, Type: "F", Status: 1},
			{Name: "SysRoleCreate", Title: "角色新增", Permission: "sys:role:create", ParentID: 7, Sort: 1, Type: "F", Status: 1},
			{Name: "SysRoleEdit", Title: "角色编辑", Permission: "sys:role:edit", ParentID: 7, Sort: 2, Type: "F", Status: 1},
			{Name: "SysRoleDelete", Title: "角色删除", Permission: "sys:role:delete", ParentID: 7, Sort: 3, Type: "F", Status: 1},
			{Name: "SysMenuCreate", Title: "菜单新增", Permission: "sys:menu:create", ParentID: 8, Sort: 1, Type: "F", Status: 1},
			{Name: "SysMenuEdit", Title: "菜单编辑", Permission: "sys:menu:edit", ParentID: 8, Sort: 2, Type: "F", Status: 1},
			{Name: "SysMenuDelete", Title: "菜单删除", Permission: "sys:menu:delete", ParentID: 8, Sort: 3, Type: "F", Status: 1},
			{Name: "SysDeptCreate", Title: "部门新增", Permission: "sys:dept:create", ParentID: 9, Sort: 1, Type: "F", Status: 1},
			{Name: "SysDeptEdit", Title: "部门编辑", Permission: "sys:dept:edit", ParentID: 9, Sort: 2, Type: "F", Status: 1},
			{Name: "SysDeptDelete", Title: "部门删除", Permission: "sys:dept:delete", ParentID: 9, Sort: 3, Type: "F", Status: 1},
			{Name: "SysDictCreate", Title: "字典新增", Permission: "sys:dict:create", ParentID: 10, Sort: 1, Type: "F", Status: 1},
			{Name: "SysDictEdit", Title: "字典编辑", Permission: "sys:dict:edit", ParentID: 10, Sort: 2, Type: "F", Status: 1},
			{Name: "SysDictDelete", Title: "字典删除", Permission: "sys:dict:delete", ParentID: 10, Sort: 3, Type: "F", Status: 1},
		}
		if err := gdb.Create(&menus).Error; err != nil {
			return err
		}
		ids := ""
		for _, m := range menus {
			ids += uintToStr(m.ID) + ","
		}
		gdb.Model(&admin).Update("menu_ids", trimComma(ids))

		// Dict
		dict := models.SysDict{Name: "订单状态", Type: "order_status", Status: 1, Remark: "order status"}
		if err := gdb.Create(&dict).Error; err != nil {
			return err
		}
		details := []models.SysDictDetail{
			{DictID: dict.ID, Label: "待付款", Value: "pending", Sort: 1, Status: 1},
			{DictID: dict.ID, Label: "已付款", Value: "paid", Sort: 2, Status: 1},
			{DictID: dict.ID, Label: "已发货", Value: "shipped", Sort: 3, Status: 1},
			{DictID: dict.ID, Label: "已完成", Value: "done", Sort: 4, Status: 1},
			{DictID: dict.ID, Label: "已取消", Value: "cancel", Sort: 5, Status: 1},
		}
		if err := gdb.Create(&details).Error; err != nil {
			return err
		}
	}

	// Admin user
	var userCount int64
	gdb.Model(&models.SysUser{}).Count(&userCount)
	if userCount == 0 {
		var adminRole models.SysRole
		gdb.Where("role_key = ?", "admin").First(&adminRole)
		u := models.SysUser{
			Username: "admin", NickName: "管理员", Phone: "13800000000",
			Email: "admin@example.com", Status: 1, RoleID: adminRole.ID, DeptID: 0, Remark: "seed",
		}
		if err := u.SetPassword("admin123"); err != nil {
			return err
		}
		if err := gdb.Create(&u).Error; err != nil {
			return err
		}
		// Sample products & orders
		products := []models.Product{
			{Name: "机械键盘", SKU: "KB-001", Category: "外设", Price: 399.00, Stock: 120, Status: 1},
			{Name: "人体工学椅", SKU: "CH-001", Category: "家具", Price: 1299.00, Stock: 30, Status: 1},
			{Name: "4K显示器", SKU: "MN-001", Category: "显示器", Price: 1899.00, Stock: 0, Status: 0},
		}
		if err := gdb.Create(&products).Error; err != nil {
			return err
		}
	}

	// System configs (参考 go-admin sys_config)
	var cfgCount int64
	gdb.Model(&models.SysConfig{}).Count(&cfgCount)
	if cfgCount == 0 {
		configs := []models.SysConfig{
			{ConfigName: "系统名称", ConfigKey: "sys.name", ConfigValue: "Go Refine Admin", ConfigType: "Y", IsFrontend: 1, Remark: "前台展示"},
			{ConfigName: "系统版本", ConfigKey: "sys.version", ConfigValue: "1.0", ConfigType: "Y", IsFrontend: 1, Remark: ""},
			{ConfigName: "登录验证码", ConfigKey: "sys.login.captcha", ConfigValue: "1", ConfigType: "N", IsFrontend: 0, Remark: "1开启 0关闭(预留)"},
		}
		if err := gdb.Create(&configs).Error; err != nil {
			return err
		}
	}

	// Monitor + config menus: idempotent, so pre-existing databases get them too.
	if err := ensureMonitorMenus(gdb); err != nil {
		return err
	}
	// Log retention parameters (same idempotent rule).
	if err := ensureRetentionConfigs(gdb); err != nil {
		return err
	}
	return nil
}

// ensureMonitorMenus creates the monitor/config menu nodes when missing and
// grants them to the admin role. Name-keyed (not id-keyed) so it is safe to
// run on every boot against both fresh and existing databases.
func ensureMonitorMenus(gdb *gorm.DB) error {
	var system models.SysMenu
	if err := gdb.Where("name = ?", "System").First(&system).Error; err != nil {
		return nil // no system tree yet; fresh seed will create everything
	}
	defs := []models.SysMenu{
		{Name: "SysOperLog", Title: "操作日志", Icon: "log", Path: "/monitor/oper-log", Component: "monitor/oper-log", ParentID: system.ID, Sort: 6, Type: "C", Status: 1},
		{Name: "SysLoginLog", Title: "登录日志", Icon: "login", Path: "/monitor/login-log", Component: "monitor/login-log", ParentID: system.ID, Sort: 7, Type: "C", Status: 1},
		{Name: "SysOnline", Title: "在线用户", Icon: "online", Path: "/monitor/online", Component: "monitor/online", ParentID: system.ID, Sort: 8, Type: "C", Status: 1},
		{Name: "SysServer", Title: "服务监控", Icon: "server", Path: "/monitor/server", Component: "monitor/server", ParentID: system.ID, Sort: 9, Type: "C", Status: 1},
		{Name: "SysApiDocs", Title: "接口文档", Icon: "api", Path: "/docs/api", Component: "docs/api", ParentID: system.ID, Sort: 11, Type: "C", Status: 1},
		{Name: "SysConfig", Title: "参数管理", Icon: "config", Path: "/system/config", Component: "system/config", ParentID: system.ID, Sort: 10, Type: "C", Status: 1},
	}
	var configID uint
	for i := range defs {
		var m models.SysMenu
		if err := gdb.Where("name = ?", defs[i].Name).First(&m).Error; err != nil {
			if err := gdb.Create(&defs[i]).Error; err != nil {
				return err
			}
			m = defs[i]
		}
		if defs[i].Name == "SysConfig" {
			configID = m.ID
		}
	}
	fdefs := []models.SysMenu{
		{Name: "SysConfigCreate", Title: "参数新增", Permission: "sys:config:create", ParentID: configID, Sort: 1, Type: "F", Status: 1},
		{Name: "SysConfigEdit", Title: "参数编辑", Permission: "sys:config:edit", ParentID: configID, Sort: 2, Type: "F", Status: 1},
		{Name: "SysConfigDelete", Title: "参数删除", Permission: "sys:config:delete", ParentID: configID, Sort: 3, Type: "F", Status: 1},
	}
	var newIDs []string
	for i := range fdefs {
		var m models.SysMenu
		if err := gdb.Where("name = ?", fdefs[i].Name).First(&m).Error; err != nil {
			fdefs[i].ParentID = configID
			if err := gdb.Create(&fdefs[i]).Error; err != nil {
				return err
			}
			m = fdefs[i]
		}
		newIDs = append(newIDs, uintToStr(m.ID))
	}
	for i := range defs {
		var m models.SysMenu
		gdb.Where("name = ?", defs[i].Name).First(&m)
		if m.ID != 0 {
			newIDs = append(newIDs, uintToStr(m.ID))
		}
	}
	var admin models.SysRole
	if err := gdb.Where("role_key = ?", "admin").First(&admin).Error; err != nil {
		return nil
	}
	have := map[string]bool{}
	for _, id := range splitMenuIDs(admin.MenuIDs) {
		have[id] = true
	}
	ids := admin.MenuIDs
	for _, id := range newIDs {
		if !have[id] {
			ids += "," + id
			have[id] = true
		}
	}
	return gdb.Model(&admin).Update("menu_ids", trimComma(ids)).Error
}

func splitMenuIDs(s string) []string {
	var out []string
	for _, p := range strings.Split(s, ",") {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

func uintToStr(u uint) string {
	if u == 0 {
		return "0"
	}
	var b []byte
	for u > 0 {
		b = append([]byte{byte('0' + u%10)}, b...)
		u /= 10
	}
	return string(b)
}

func trimComma(s string) string {
	if len(s) > 0 && s[len(s)-1] == ',' {
		return s[:len(s)-1]
	}
	return s
}
