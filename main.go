package main

import (
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
		v1.PUT("/posts/:id", postHandler.UpdatePostWithTags)
        v1.GET("/posts", postHandler.GetAllPosts)
	}

	// Establish many-to-many relationship

	router.Run(":8080")

}

// func CreatePost(c *gin.Context) {
//     // Ambil data post dari request body
//     var post model.Post
//     if err := c.ShouldBindJSON(&post); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     // Buat relasi post dengan tag
//     var tagNames []string
//     if len(post.TagNames) > 0 {
// 		tagNames = append(tagNames, post.TagNames...)
// 	}

//     // Simpan post dan relasi post dengan tag

//     // Buat slice untuk menyimpan tags
//     var tags []model.Tag

//     // Looping tag names
//     for _, tagName := range tags {
//         // Simpan tag ke databasef
//         db.Create(&tagName)

//         // Append tag ke slice tags
//         tags = append(tags, tagName)
//     }

//     // Simpan post dan relasi post dengan tag
//     db.Create(&post)
//     for _, tag := range tags {
//         db.Model(&post).Association("Tags").Append(&tag)
//     }

//     // Kirim respon
//     c.JSON(http.StatusOK, gin.H{"data": post})
// }
