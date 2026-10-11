// Package main is the go-refine-admin backend server.
//
//	@title			Go Refine Admin API
//	@version		1.0
//	@description	Go + gin + gorm go-refine-admin scaffold with Actions / regular dual-track and type-safe queries.
//	@BasePath		/api/v1
//	@schemes		http
//	@securityDefinitions.apikey	Bearer
//	@in				header
//	@name			Authorization
//	@description	"Type 'Bearer <token>' to authorize."
//
//	@contact.name	go-refine-admin
//
//go:generate swag init -g main.go -o docs
package main

import (
	"log"

	"go-refine-admin/backend/config"
	"go-refine-admin/backend/db"
	"go-refine-admin/backend/genquery"
	"go-refine-admin/backend/router"
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
	// Expired audit/login logs are purged at boot and every 24h
	// (window from sys.log.*-retention-days, floor 30 days).
	go db.StartLogRetentionJanitor(gdb)
	r := router.Build(gdb)
	log.Printf("go-refine-admin listening on %s", cfg.Server.Addr())
	if err := r.Run(cfg.Server.Addr()); err != nil {
		log.Fatalf("server: %v", err)
	}
}
