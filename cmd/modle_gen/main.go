// cmd/gen/main.go
package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"strings"
)

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/xxx?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic("failed to connect database")
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:           "internal/model/gen/gen",
		ModelPkgPath:      "internal/model/gen",
		FieldWithIndexTag: true,
		FieldWithTypeTag:  true,
		Mode:              gen.WithoutContext,
	})

	g.UseDB(db)
	g.WithImportPkgPath("lls_api/internal/model/base")
	// 仅映射 time.Time → NullTime，其他保持原生类型
	dataMap := map[string]func(columnType gorm.ColumnType) (dataType string){
		"datetime": func(columnType gorm.ColumnType) (dataType string) {
			if n, ok := columnType.Nullable(); ok && n {
				return "base.NullTime"
			}
			return "time.Time"
		},
		// bool mapping
		"tinyint": func(columnType gorm.ColumnType) (dataType string) {
			ct, _ := columnType.ColumnType()
			if strings.HasPrefix(ct, "tinyint(1)") {
				return "bool"
			}
			return "byte"
		},
		// 显式处理 int
		"int": func(columnType gorm.ColumnType) (dataType string) {
			if strings.HasSuffix(columnType.Name(), "id") {
				if n, ok := columnType.Nullable(); ok && n {
					return "base.NullID"
				}
				return "base.ID"
			} else {
				if n, ok := columnType.Nullable(); ok && n {
					return "*int32"
				}
				return "int32"
			}
		},
	}
	g.WithDataTypeMap(dataMap)

	// 获取所有表名
	tables, err := db.Migrator().GetTables()
	if err != nil {
		panic(err)
	}

	// 定义需要特殊处理的表
	specialTables := map[string]string{
		"blue_usersaas": "BlueUsersaas",
		"blue_saas":     "BlueSaas",
	}

	var models []any

	for _, table := range tables {
		if newName, ok := specialTables[table]; ok {
			// 使用指定的结构体名
			models = append(models, g.GenerateModelAs(table, newName))
		} else if table == "blue_contract" {
			// 特殊字段处理
			models = append(models, g.GenerateModel(table,
				gen.FieldRename("contractId", "ContractIDStr"),
				gen.FieldRename("contract_templateId", "ContractTemplateIDStr"),
			))
		} else {
			// 默认生成
			models = append(models, g.GenerateModel(table))
		}
	}

	g.ApplyBasic(models...)
	g.Execute()
}
