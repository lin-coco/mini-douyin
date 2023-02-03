package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"interaction.rpc/dal/model"
)

func main() {
	g := gen.NewGenerator(gen.Config{
		OutPath:       "../../dal/query",
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		FieldNullable: true,
	})
	dsn := "xys:232020ctt@@tcp(rm-uf6e4xr978w748b9w7o.mysql.rds.aliyuncs.com:3306)/sql_test?charset=utf8mb4&parseTime=true&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{ //连接数据库
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		AllowGlobalUpdate:      true,
	})
	if err != nil {
		panic("failed to connect database")
	}
	err = db.AutoMigrate(&model.UserFavourite{}, &model.Comment{})
	if err != nil {
		panic("autoMigrate failed")
	}
	g.UseDB(db)
	g.ApplyBasic(model.UserFavourite{}, model.Comment{})

	g.Execute()
}
