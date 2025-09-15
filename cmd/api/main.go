package main

import (
	"context"
	"goblog/internal/handlers"
	"goblog/internal/models"
	"goblog/internal/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type config struct {
	JWT struct {
		Secret      string `yaml:"secret"`
		ExpireHours int    `yaml:"expire_hours"`
	} `yaml:"jwt"`
	MySQL struct {
		DSN string `yaml:"dsn"`
	} `yaml:"mysql"`
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
}

func loadConfig() config {
	var conf config
	file, err := os.Open("configs/apps.yaml")
	if err != nil {
		log.Fatal("打开配置文件失败", err)
	}
	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(&conf); err != nil {
		log.Fatal("解析配置文件失败", err)
	}
	return conf
}

func migrateDB(db *gorm.DB) {
	log.Println("Migrating database...正在数据库迁移")
	err := db.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Comment{},
	)
	if err != nil {
		log.Fatal("数据库迁移失败", err)
	}
	log.Println("数据库迁移完成")
}

func main() {
	conf := loadConfig()
	jwtSecret := []byte(conf.JWT.Secret)
	serverPort := conf.Server.Port

	db, err := gorm.Open(mysql.Open(conf.MySQL.DSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		log.Fatal("数据库连接失败了 :", err)
	}

	//执行迁移
	migrateDB(db)

	//创建service，依赖注入db 和 jwtSecret
	authService := services.NewAuthService(db, jwtSecret)
	postService := services.NewPostService(db)
	commentService := services.NewCommentService(db)

	//创建handler，依赖注入service
	authHandler := handlers.NewAuthHandler(authService)
	postHandler := handlers.NewPostHandler(postService)
	commentHandler := handlers.NewCommentHandler(commentService)

	//初始化gin，注册路由
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE,PATCH")
		c.Header("Access-Control-Allow-Headers", "Authorization,Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		api.POST("/login", authHandler.Login)
		api.POST("/register", authHandler.Register)

		api.GET("/post", postHandler.CreatePost)
		api.GET("/post/:id", postHandler.GetPostByID)

		api.GET("/post/:id/comments", commentHandler.CreateComment)
		api.GET("/post/:id/comments", commentHandler.GetComment)
	}
	server := &http.Server{
		Addr:    ":" + serverPort,
		Handler: r,
	}

	go func() {
		log.Printf("server listening on port %s", serverPort)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("服务启动失败", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...正在关闭")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown failed:", err)
	}

	log.Println("Server Closed Successfully ")
}
