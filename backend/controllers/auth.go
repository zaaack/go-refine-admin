package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"go-refine-admin/backend/config"
	"go-refine-admin/backend/core/auth"
	"go-refine-admin/backend/core/response"
	"go-refine-admin/backend/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LoginRequest is the login payload (also emitted into swagger models).
type LoginRequest struct {
	Username string `json:"username" example:"admin"`
	Password string `json:"password" example:"admin123"`
	Code     string `json:"code" example:"1234"` // captcha
	UUID     string `json:"uuid" example:"captcha-id"` // captcha id
}

// Login is the regular-mode authentication endpoint. It authenticates a user,
// loads the role and issues a JWT. This is hand-written (not Actions mode)
// because it carries bespoke business rules.
// It is served both at POST /login (root, legacy) and POST /api/v1/login
// (swagger basePath) so the generated client can call it directly.
//
//	@Summary	Login
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		LoginRequest	true	"credentials"
//	@Success	200		{object}	response.Envelope
//	@Router		/login [post]
func Login(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, http.StatusBadRequest, 400, err, "")
			return
		}
		// Captcha check, skipped in dev mode (same rule as go-admin).
		if config.Get().Server.Mode != "dev" {
			if !VerifyCaptcha(req.UUID, req.Code) {
				response.Error(c, http.StatusUnauthorized, 401, errors.New("验证码错误"), "")
				return
			}
		}
		var u models.SysUser
		if err := gdb.Where("username = ?", req.Username).First(&u).Error; err != nil {
			response.Error(c, http.StatusUnauthorized, 401, errors.New("用户名或密码错误"), "")
			return
		}
		if u.Status != 1 {
			response.Error(c, http.StatusForbidden, 403, errors.New("账号已被禁用"), "")
			return
		}
		if !u.ComparePassword(req.Password) {
			response.Error(c, http.StatusUnauthorized, 401, errors.New("用户名或密码错误"), "")
			return
		}
		var role models.SysRole
		gdb.First(&role, u.RoleID)
		perms := rolePerms(gdb, role)

		token, err := auth.GenerateToken(auth.Claims{
			UserID:    u.ID,
			Username:  u.Username,
			RoleKey:   role.RoleKey,
			RoleID:    role.ID,
			DataScope: role.DataScope,
			Perms:     perms,
		})
		if err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		response.OK(c, gin.H{
			"token": token,
			"user": gin.H{
				"id":       u.ID,
				"username": u.Username,
				"nickName": u.NickName,
				"roleKey":  role.RoleKey,
				"avatar":   u.Avatar,
			},
			"perms": perms,
		}, "")
	}
}

// rolePerms collects the Permission strings of the role's granted menus.
// The admin role holds "*" (everything).
func rolePerms(gdb *gorm.DB, role models.SysRole) []string {
	if role.RoleKey == "admin" {
		return []string{"*"}
	}
	var ids []interface{}
	for _, s := range strings.Split(role.MenuIDs, ",") {
		s = strings.TrimSpace(s)
		if id, err := strconv.ParseUint(s, 10, 64); err == nil {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return nil
	}
	var menus []models.SysMenu
	gdb.Where("id IN ?", ids).Find(&menus)
	var perms []string
	for _, m := range menus {
		if m.Permission != "" {
			perms = append(perms, m.Permission)
		}
	}
	return perms
}

// Me returns the current authenticated user profile.
//
//	@Summary	Current user
//	@Tags		Auth
//	@Produce	json
//	@Success	200	{object}	response.Envelope
//	@Router		/me [get]
//	@Security	Bearer
func Me(c *gin.Context) {
	claims := auth.GetClaims(c)
	if claims == nil {
		response.Error(c, http.StatusUnauthorized, 401, errors.New("unauthorized"), "")
		return
	}
	response.OK(c, gin.H{
		"id":       claims.UserID,
		"username": claims.Username,
		"roleKey":  claims.RoleKey,
	}, "")
}

// MenuTree returns the menu tree granted to the current user's role (regular
// mode). The admin role sees everything; other roles are filtered by the
// role's comma-separated menu_ids, with ancestor nodes kept so the tree renders.
//
//	@Summary	Menu tree (role-scoped)
//	@Tags		Auth
//	@Produce	json
//	@Success	200	{object}	response.Envelope
//	@Router		/menu/tree [get]
//	@Security	Bearer
func MenuTree(gdb *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var menus []models.SysMenu
		if err := gdb.Order("sort asc").Find(&menus).Error; err != nil {
			response.Error(c, http.StatusInternalServerError, 500, err, "")
			return
		}
		granted := filterMenusByRole(c, gdb, menus)
		response.OK(c, buildMenuTree(granted, 0), "")
	}
}

// filterMenusByRole reduces menus to those granted to the caller's role.
func filterMenusByRole(c *gin.Context, gdb *gorm.DB, menus []models.SysMenu) []models.SysMenu {
	claims := auth.GetClaims(c)
	if claims == nil || claims.RoleKey == "admin" {
		return menus
	}
	var role models.SysRole
	if err := gdb.First(&role, claims.RoleID).Error; err != nil {
		return nil
	}
	grantedIDs := map[uint]bool{}
	for _, s := range strings.Split(role.MenuIDs, ",") {
		s = strings.TrimSpace(s)
		if id, err := strconv.ParseUint(s, 10, 64); err == nil {
			grantedIDs[uint(id)] = true
		}
	}
	if len(grantedIDs) == 0 {
		return nil
	}
	byID := make(map[uint]models.SysMenu, len(menus))
	for _, m := range menus {
		byID[m.ID] = m
	}
	var granted []models.SysMenu
	inSet := map[uint]bool{}
	// Walk up to keep ancestors of granted nodes in the tree.
	for _, m := range menus {
		if !grantedIDs[m.ID] {
			continue
		}
		for p := m; p.ID != 0 && !inSet[p.ID]; p = byID[p.ParentID] {
			granted = append(granted, p)
			inSet[p.ID] = true
		}
	}
	return granted
}

type MenuNode struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Title     string     `json:"title"`
	Icon      string     `json:"icon"`
	Path      string     `json:"path"`
	Component string     `json:"component"`
	Type      string     `json:"type"`
	Sort      int        `json:"sort"`
	Children  []MenuNode `json:"children,omitempty"`
}

func buildMenuTree(menus []models.SysMenu, parentID uint) []MenuNode {
	var nodes []MenuNode
	for _, m := range menus {
		if m.ParentID != parentID {
			continue
		}
		node := MenuNode{
			ID:        m.ID,
			Name:      m.Name,
			Title:     m.Title,
			Icon:      m.Icon,
			Path:      m.Path,
			Component: m.Component,
			Type:      m.Type,
			Sort:      m.Sort,
			Children:  buildMenuTree(menus, m.ID),
		}
		nodes = append(nodes, node)
	}
	return nodes
}
