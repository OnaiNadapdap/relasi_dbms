package repository

import (
	"github.com/onainadapdap1/kampus_tutor/coba_many_to_many/helper"
	"github.com/onainadapdap1/kampus_tutor/coba_many_to_many/model"
	"gorm.io/gorm"
)

type PostRepository interface {
	UpdatePost(post *model.Post) error
	CreateTag(tag *model.Tag) error
	FindPostByID(id uint) (*model.Post, error)
	FindPosts(perPage int, page int) ([]model.Post, error)
	CountAllPost() (int64, error)
	// FindTagByName(name string) (model.Tag, error)
	DeletePostTagsByPostID(postID uint) error
	DeleteTagsByPostID(postID uint) error

	// this for add post with tags and delete tag from post
	CreatePost(post *model.Post) error
	AppendTagToPost(post model.Post, tag model.Tag) error
	GetTagByID(tagID uint) (model.Tag, error)
	GetPostByID(postID uint) (model.Post, error)
	DeleteTagFromPost(post model.Post, tag model.Tag) error
	// this is for retrieve all post base id tag
	// GetPostByTagID(tagID uint) (model.Tag, error)
	GetPostByTagID2(tagID uint) (model.Tag, error)
	GetPostAssociationByTag(tag model.Tag) ([]model.Post, error)
}

type postRepository struct {
	DB *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{DB: db}
}

func (r *postRepository) CountAllPost() (int64, error) {
	var posts []model.Post
	var totalRows int64

	tx := r.DB.Begin()
	if err := tx.Debug().Find(&posts).Count(&totalRows).Error; err != nil {
		return totalRows, err
	}

	return totalRows, nil
}

func (r *postRepository) GetPostByTagID(tagID uint) (model.Tag, error) {
	var tag model.Tag
	tx := r.DB.Begin()
	if err := tx.Debug().Preload("Posts").Take(&tag, "id = ?", tagID).Error; err != nil {
		return tag, err
	}
	return tag, nil
}

func (r *postRepository) GetPostByTagID2(tagID uint) (model.Tag, error) {
	var tag model.Tag
	tx := r.DB.Begin()
	if err := tx.Debug().Take(&tag, "id = ?", tagID).Error; err != nil {
		return tag, err
	}

	return tag, nil
	// var product Product
	// err := db.Take(&product, "id = ?", "P001").Error
	// assert.Nil(t, err)

	// var users []User
	// err = db.Model(&product).Where("users.first_name LIKE ?", "User%").Association("LikedByUsers").Find(&users)
	// assert.Nil(t, err)
	// assert.Equal(t, 1, len(users))
}

func (r *postRepository) GetPostAssociationByTag(tag model.Tag) ([]model.Post, error) {
	tx := r.DB.Begin()
	var posts []model.Post

	// Explicitly join the posts and tags tables
	if err := tx.Debug().Raw("SELECT p.id, p.title, p.content, t.id AS tag_id, t.name AS tag_name FROM posts AS p INNER JOIN post_tags AS pt ON p.id = pt.post_id INNER JOIN tags AS t ON pt.tag_id = t.id WHERE t.name LIKE ?", "%"+tag.Name+"%").
		Scan(&posts).Error; err != nil {
		return posts, err
	}

	return posts, nil
}

func (r *postRepository) DeleteTagFromPost(post model.Post, tag model.Tag) error {
	tx := r.DB.Begin()
	if err := tx.Debug().Model(&post).Association("Tags").Delete(&tag); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (r *postRepository) GetPostByID(postID uint) (model.Post, error) {
	var post model.Post
	tx := r.DB.Begin()
	if err := tx.Debug().Take(&post, "id = ?", postID).Error; err != nil {
		return post, err
	}

	return post, nil
}

func (r *postRepository) GetTagByID(tagID uint) (model.Tag, error) {
	var tag model.Tag
	tx := r.DB.Begin()
	if err := tx.Debug().Take(&tag, "id = ?", tagID).Error; err != nil {
		return tag, err
	}
	return tag, nil
}

func (r *postRepository) AppendTagToPost(post model.Post, tag model.Tag) error {
	tx := r.DB.Begin()
	if err := tx.Debug().Model(&post).Association("Tags").Append(&tag); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *postRepository) FindPosts(perPage int, page int) ([]model.Post, error) {
	var posts []model.Post
	tx := r.DB.Begin()
	err := tx.Debug().Scopes(helper.PaginationScopes(page, perPage)).Preload("Tags").Find(&posts).Error
	if err != nil {
		return posts, err
	}

	return posts, nil
}

// func (r *postRepository) FindTagByName(name string) (model.Tag, error) {
//     tx := r.DB.Begin()
//     var tag model.Tag
//     if err := tx.Debug().Where("name = ?", name).First(&tag).Error; err != nil {
//         return tag, err
//     }
//     return tag, nil
// }

func (r *postRepository) CreatePost(post *model.Post) error {
	tx := r.DB.Begin()
	// return tx.Debug().Create(post).Error
	if err := tx.Debug().Create(post).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()

	return nil
}

func (r *postRepository) UpdatePost(post *model.Post) error {
	return r.DB.Save(post).Error
}

func (r *postRepository) FindPostByID(id uint) (*model.Post, error) {
	var post model.Post
	if err := r.DB.Preload("Tags").First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) CreateTag(tag *model.Tag) error {
	return r.DB.Create(tag).Error
}

func (r *postRepository) DeletePostTagsByPostID(postID uint) error {
	return r.DB.Debug().Where("post_id = ?", postID).Delete(&model.PostTag{}).Error
}

func (r *postRepository) DeleteTagsByPostID(postID uint) error {
	tx := r.DB.Begin()

	if err := tx.Debug().Exec("DELETE FROM post_tags WHERE post_id = ?", postID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Debug().Exec("DELETE FROM tags WHERE id NOT IN (SELECT DISTINCT tag_id FROM post_tags)").Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// func NewPostRepository(db *gorm.DB) *PostRepository {
//     return &PostRepository{DB: db}
// }

// func (r *PostRepository) CreatePost(post *model.Post) error {
//     return r.DB.Create(post).Error
// }

// func (r *PostRepository) UpdatePostTags(post *model.Post, tagNames []string) error {
//     return r.DB.Transaction(func(tx *gorm.DB) error {
//         for _, tagName := range tagNames {
//             tag := model.Tag{Name: tagName}
//             if err := tx.FirstOrCreate(&tag, model.Tag{Name: tagName}).Error; err != nil {
//                 return err
//             }

//             post.Tags = append(post.Tags, tag)
//             pt := model.PostTag{PostID: post.ID, TagID: tag.ID}
//             if err := tx.Create(&pt).Error; err != nil {
//                 return err
//             }
//         }

//         // Remove any tags not in the updated list
//         var tagIDs []uint
//         for _, tag := range post.Tags {
//             tagIDs = append(tagIDs, tag.ID)
//         }
//         if err := tx.Where("post_id = ? AND tag_id NOT IN ?", post.ID, tagIDs).Delete(&model.PostTag{}).Error; err != nil {
//             return err
//         }

//         return nil
//     })
// }
