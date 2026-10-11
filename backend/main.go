// Package main is the enterprise-admin backend server.
//
//	@title			Enterprise Admin API
//	@version		1.0
//	@description	Go + gin + gorm enterprise admin scaffold with Actions / regular dual-track and type-safe queries.
//	@BasePath		/api/v1
//	@schemes		http
//	@securityDefinitions.apikey	Bearer
//	@in				header
//	@name			Authorization
//	@description	"Type 'Bearer <token>' to authorize."
//
//	@contact.name	enterprise-admin
//
//go:generate swag init -g main.go -o docs
package main

import (
	"log"

	"enterprise-admin/backend/config"
	"enterprise-admin/backend/db"
	"enterprise-admin/backend/genquery"
	"enterprise-admin/backend/router"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	gdb, err := db.Open(cfg.Database)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	if err := db.Migrate(gdb); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	// Bind gorm/gen's generated query package to the live database so the
	// type-safe accessors (genquery.Order.Amount.Gte, ...) operate on it.
	genquery.SetDefault(gdb)
	r := router.Build(gdb)
	log.Printf("enterprise-admin listening on %s", cfg.Server.Addr())
	if err := r.Run(cfg.Server.Addr()); err != nil {
		log.Fatalf("server: %v", err)
	}
}
