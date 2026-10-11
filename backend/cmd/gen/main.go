// Command gen generates a gorm/gen type-safe query package from the project
// models. Run with:  go run ./cmd/gen
//
// gorm/gen is gorm's built-in code generator: it emits typed field accessors
// (e.g. genquery.Order.Amount.Gte(100)) so query conditions are checked at
// compile time instead of being raw strings.
package main

import (
	"log"

	"go-refine-admin/backend/config"
	"go-refine-admin/backend/db"
	"go-refine-admin/backend/models"

	"gorm.io/gen"
)

func main() {
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	gdb, err := db.Open(cfg.Database)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath: "./genquery",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	g.UseDB(gdb)
	g.ApplyBasic(
		models.SysUser{},
		models.SysRole{},
		models.SysMenu{},
		models.SysDept{},
		models.SysDict{},
		models.SysDictDetail{},
		models.SysOperLog{},
		models.Upload{},
		models.Product{},
		models.Order{},
	)
	g.Execute()
	log.Println("generated type-safe query package at ./genquery")
}
