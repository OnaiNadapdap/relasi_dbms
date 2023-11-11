package service

import (
	"log"

	"github.com/onainadapdap1/kampus_tutor/coba_many_to_many/model"
	"github.com/onainadapdap1/kampus_tutor/coba_many_to_many/repository"
)

type PostService interface {
	CreatePostWithTags(post *model.Post, tagNames []string) error
	UpdatePostByIDWithTags(id uint, updatedPost *model.Post) error
	GetPostByID(id uint) (*model.Post, error)
	GetPosts() ([]model.Post, error)
}

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{postRepo: postRepo}
}

func (s *postService) CreatePostWithTags(post *model.Post, tagNames []string) error {
	for _, tagName := range tagNames {
		tag := model.Tag{Name: tagName}
		post.Tags = append(post.Tags, tag)
	}
	return s.postRepo.CreatePost(post)
}

func (s *postService) UpdatePostByIDWithTags(id uint, updatedPost *model.Post) error {
	existingPost, err := s.postRepo.FindPostByID(id)
	if err != nil {
		return err
	}

	// Delete existing tags associated with the post
	if err := s.postRepo.DeleteTagsByPostID(existingPost.ID); err != nil {
		log.Println("error here")
		return err
	}
	log.Println("stop after delete tags")

	for _, tagName := range updatedPost.TagNames {
		tag := model.Tag{Name: tagName}
		err = s.postRepo.CreateTag(&tag)
		if err != nil {
			return err
		}

		existingPost.Tags = append(existingPost.Tags, tag)
	}
	existingPost.Title = updatedPost.Title
	existingPost.Content = updatedPost.Content

	return s.postRepo.UpdatePost(existingPost)
}

func (s *postService) GetPostByID(id uint) (*model.Post, error) {
	return s.postRepo.FindPostByID(id)
}

func (s *postService) GetPosts() ([]model.Post, error) {
	posts, err := s.postRepo.FindPosts()
	if err != nil {
		return posts, err
	}

	return posts, nil
}

// func (s *PostService) CreatePostWithTags(post *model.Post, tagNames []string) error {
//     if err := s.PostRepo.CreatePost(post); err != nil {
//         return err
//     }

//     for _, tagName := range tagNames {
//         tag := model.Tag{Name: tagName}
//         if err := s.PostRepo.CreateTag(&tag); err != nil {
//             return err
//         }
//         post.Tags = append(post.Tags, tag)
//     }
// 	// return nil

// 	if err := s.PostRepo.DB.Save(post).Error; err != nil {
// 		return err
// 	}
// 	return nil
// }

// func (s *PostService) UpdatePostWithTags(post *model.Post, tagNames []string) error {
//     return s.PostRepo.UpdatePostTags(post, tagNames)
// }
