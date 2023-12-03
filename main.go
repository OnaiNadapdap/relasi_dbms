package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/onainadapdap1/kampus_tutor/coba_many_to_many/handler"
	"github.com/onainadapdap1/kampus_tutor/coba_many_to_many/model"
	"github.com/onainadapdap1/kampus_tutor/coba_many_to_many/repository"
	"github.com/onainadapdap1/kampus_tutor/coba_many_to_many/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func AutoMigrate(db *gorm.DB) {
    db.AutoMigrate(&model.Post{})
    db.AutoMigrate(&model.Tag{})

    // Tambahkan kaskade ke post_tags
    // db.Exec("ALTER TABLE post_tags ADD CONSTRAINT fk_post_tags_tags FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE;")
	// here is an update
	// updated again
}


func main() {
	var err error

	// MySQL database connection
	dsn := "root:my-secret-pw-23@tcp(127.0.0.1:3306)/coba_many?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}

	// AutoMigrate creates tables
    AutoMigrate(db)
	// err = db.AutoMigrate(&model.Post{}, 
    //     &model.Tag{})
	// if err != nil {
	// 	panic("Failed to auto migrate tables")
	// }

	postRepo := repository.NewPostRepository(db)
	postServ := service.NewPostService(postRepo)
	postHandler := handler.NewPostHandler(postServ)

	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		v1.POST("/posts", postHandler.CreatePostWithTags)
		v1.DELETE("/posts/:postID/tags/:tagID", postHandler.DeleteTagFromPost)
		v1.GET("/posts/tags/:tagID", postHandler.GetPostByTagID)
		v1.POST("/posts/:postID", postHandler.AppendTagToPost)


		// v1.PUT("/posts/:id", postHandler.UpdatePostWithTags)
        v1.GET("/posts", postHandler.GetAllPosts)

	}

	log.Println("starting point yeah")
	log.Println("again")
	log.Println("starting line up")
	log.Println("try git revert")
	router.Run(":8080")

}
