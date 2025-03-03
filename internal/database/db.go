package database

import (
	"fmt"
	"log"
	"sync"

	"daily-report/internal/config"
	"daily-report/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// GetDB 返回数据库连接实例
func GetDB() *gorm.DB {
	once.Do(func() {
		cfg := config.GetConfig()

		log.Println("正在初始化数据库连接...")

		// 首先连接MySQL服务器（不指定数据库）
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s&parseTime=true&loc=%s&collation=utf8mb4_unicode_ci",
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Charset,
			cfg.Database.Loc,
		)

		tempDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatal("连接MySQL服务器失败:", err)
		}
		log.Println("成功连接到MySQL服务器")

		// 检查数据库是否存在
		var count int64
		row := tempDB.Raw(fmt.Sprintf("SELECT COUNT(*) FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '%s'", cfg.Database.Name)).Scan(&count)
		if row.Error != nil {
			log.Fatal("检查数据库是否存在时出错:", row.Error)
		}

		// 如果数据库不存在，则创建
		if count == 0 {
			log.Printf("数据库 %s 不存在，正在创建...", cfg.Database.Name)
			// 创建数据库
			sql := fmt.Sprintf("CREATE DATABASE %s CHARACTER SET %s COLLATE %s_unicode_ci;",
				cfg.Database.Name, cfg.Database.Charset, cfg.Database.Charset)

			if err := tempDB.Exec(sql).Error; err != nil {
				log.Fatal("创建数据库失败:", err)
			}
			log.Printf("数据库 %s 创建成功", cfg.Database.Name)
		} else {
			log.Printf("数据库 %s 已存在，跳过创建", cfg.Database.Name)
		}

		// 连接到目标数据库
		log.Printf("正在连接到数据库 %s...", cfg.Database.Name)
		database, err := gorm.Open(mysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=true&loc=%s&collation=utf8mb4_unicode_ci",
			cfg.Database.Username,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
			cfg.Database.Charset,
			cfg.Database.Loc,
		)), &gorm.Config{})
		if err != nil {
			log.Fatal("连接数据库失败:", err)
		}
		log.Printf("成功连接到数据库 %s", cfg.Database.Name)

		// 自动迁移数据库结构
		log.Println("正在迁移数据库结构...")
		err = database.AutoMigrate(
			&model.User{},
			&model.Project{},
			&model.Report{},
			&model.Task{},
		)
		if err != nil {
			log.Fatal("数据库迁移失败:", err)
		}
		log.Println("数据库迁移成功完成")

		// 检查并创建默认管理员用户
		var adminUser model.User
		if err := database.Where("email = ?", "admin@blingsec.cn").First(&adminUser).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Println("创建默认管理员用户...")
				adminUser = model.User{
					Email:    "admin@blingsec.cn",
					Username: "admin",
					Role:     "admin",
				}
				// 使用16位随机强密码
				if err := adminUser.SetPassword("Kj#9mP$2nL5vB@8x"); err != nil {
					log.Fatal("设置管理员密码失败:", err)
				}
				if err := database.Create(&adminUser).Error; err != nil {
					log.Fatal("创建管理员用户失败:", err)
				}
				log.Println("默认管理员用户创建成功")
				log.Println("默认管理员密码: Kj#9mP$2nL5vB@8x")
			} else {
				log.Fatal("检查管理员用户时出错:", err)
			}
		} else {
			log.Println("默认管理员用户已存在")
		}

		db = database
	})
	return db
}
