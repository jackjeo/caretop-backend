package database

import (
	"log"

	"caretop-backend/models"
)

func Migrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.BlogPost{},
		&models.BlogComment{},
		&models.ForumBoard{},
		&models.ForumThread{},
		&models.ForumPost{},
		&models.ForumLike{},
		&models.ForumCollection{},
		&models.Ticket{},
		&models.TicketReply{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	seedForumBoards()
}

func seedForumBoards() {
	boards := []models.ForumBoard{
		{Name: "MindLink产品交流", Slug: "mindlink", Description: "MindLink产品使用交流与反馈", SortOrder: 1},
		{Name: "HenryIway产品交流", Slug: "henryiway", Description: "HenryIway产品使用交流与反馈", SortOrder: 2},
		{Name: "RemoteDesktop交流", Slug: "remotedesktop", Description: "远程桌面产品交流区", SortOrder: 3},
		{Name: "经验分享", Slug: "experience", Description: "技术心得与经验分享", SortOrder: 4},
		{Name: "官方公告", Slug: "announcements", Description: "官方新闻与公告", SortOrder: 5},
	}

	for _, board := range boards {
		DB.FirstOrCreate(&board, models.ForumBoard{Slug: board.Slug})
	}
}
