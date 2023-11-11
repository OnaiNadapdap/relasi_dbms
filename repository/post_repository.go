package repository

import (
	"github.com/onainadapdap1/kampus_tutor/coba_many_to_many/model"
	"gorm.io/gorm"
)

type PostRepository interface {
    CreatePost(post *model.Post) error
    UpdatePost(post *model.Post) error
    CreateTag(tag *model.Tag) error
    FindPostByID(id uint) (*model.Post, error)
    FindPosts() ([]model.Post, error)
    // FindTagByName(name string) (model.Tag, error)
    DeletePostTagsByPostID(postID uint) error
    DeleteTagsByPostID(postID uint) error
}

type postRepository struct {
    DB *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
    return &postRepository{DB: db}
}

func (r *postRepository) FindPosts() ([]model.Post, error) {
    var posts []model.Post
	tx := r.DB.Begin()
	err := tx.Debug().Preload("Tags").Find(&posts).Error
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
    return r.DB.Create(post).Error
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
