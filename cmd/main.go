package main

import (
	"Mou1ght/consts"
	"Mou1ght/internal/api/handler"
	"Mou1ght/internal/api/service"
	"Mou1ght/internal/config"
	"Mou1ght/internal/repository/instance"
	"Mou1ght/internal/service/router"
	"Mou1ght/pkg/util"
)

func init() {
	config.GetConfig()
	util.CreateDirNotExists(consts.Upload)
}

func runApp() {
	// 数据库连接
	database := instance.NewDatabase()
	postCounter := instance.NewPostCounter(database.GetCounter())
	// 实例
	userRepository := instance.NewUserRepository(database.DB)
	articleRepository := instance.NewArticleRepository(database.DB, postCounter)
	sharingRepository := instance.NewSharingRepository(database.DB, postCounter)
	messageRepository := instance.NewMessageRepository(database.DB, postCounter)
	tagRepository := instance.NewTagRepository(database.DB)
	categoryRepository := instance.NewCategoryRepository(database.DB)
	categoryLinkRepository := instance.NewCategoryLinkRepository(database.DB)
	postRepository := instance.NewPostRepository(postCounter, database.DB)
	attachmentRepository := instance.NewAttachmentRepository(database.DB)
	sharingAttachmentLinkRepository := instance.NewSharingAttachmentLinkRepository(database.DB)

	// 服务
	dtoService := service.NewDTOService(userRepository, articleRepository, tagRepository, categoryRepository, attachmentRepository, sharingAttachmentLinkRepository, postCounter)
	userService := service.NewUserService(userRepository, articleRepository, sharingRepository)
	articleService := service.NewArticleService(articleRepository, categoryRepository, categoryLinkRepository, tagRepository)
	sharingService := service.NewSharingService(sharingRepository, tagRepository, attachmentRepository, sharingAttachmentLinkRepository)
	messageService := service.NewMessageService(messageRepository)
	tagService := service.NewTagService(tagRepository)
	categoryService := service.NewCategoryService(categoryRepository, categoryLinkRepository)
	postService := service.NewPostService(articleRepository, sharingRepository, messageRepository, postRepository)
	attachmentService := service.NewAttachmentService(attachmentRepository)

	deps := router.Deps{
		UserHandler:       handler.NewUserHandler(userService),
		ArticleHandler:    handler.NewArticleHandler(articleService, dtoService),
		SharingHandler:    handler.NewSharingHandler(sharingService, dtoService),
		MessageHandler:    handler.NewMessageHandler(messageService, dtoService),
		AttachmentHandler: handler.NewAttachmentHandler(attachmentService),
		TagHandler:        handler.NewTagHandler(tagService),
		CategoryHandler:   handler.NewCategoryHandler(categoryService, categoryRepository, dtoService),
		PostHandler:       handler.NewPostHandler(articleService, categoryService, sharingService, messageService, tagService, userService, postService, dtoService),
	}

	r := router.InitRouter(deps)
	e := r.Listen(":10420")
	if e != nil {
		panic(e)
	}
}

func main() {
	runApp()
}
