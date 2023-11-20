package service

import (
	"log"

	"github.com/onainadapdap1/kampus_tutor/coba_many_to_many/model"
	"github.com/onainadapdap1/kampus_tutor/coba_many_to_many/repository"
)

type PostService interface {
	CreatePostWithTags(post *model.Post, tagNames []string) error
	DeleteTagFromPost(tagID uint, postID uint) error
	UpdatePostByIDWithTags(id uint, updatedPost *model.Post) error
	GetPostByID(id uint) (model.Post, error)
	GetPosts(perPage int, page int) ([]model.Post, int64, error)
	GetPostByTagID(tagID uint) ([]model.Post, error)
	AppendTagToPost(postID uint, tagNames []string) error 
}

type postService struct {
	postRepo repository.PostRepository
}

func NewPostService(postRepo repository.PostRepository) PostService {
	return &postService{postRepo: postRepo}
}

func (s *postService) AppendTagToPost(postID uint, tagNames []string) error {
	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		return err
	}
	
	for _, tagName := range tagNames {
		tag := model.Tag{Name: tagName}
		s.postRepo.AppendTagToPost(post, tag)
		// post.Tags = append(post.Tags, tag)
	}
	return nil
}

func (s *postService) GetPostByTagID(tagID uint) ([]model.Post, error) {
	tag, err := s.postRepo.GetPostByTagID2(tagID)
	if err != nil {
		return nil, err
	}

	posts, err := s.postRepo.GetPostAssociationByTag(tag)
	if err != nil {
		return posts, err
	}
	
	return posts, nil
}

func (s *postService) DeleteTagFromPost(tagID uint, postID uint) error {
	tag, err := s.postRepo.GetTagByID(tagID)
	if err != nil {
		return err
	}

	post, err := s.postRepo.GetPostByID(postID)
	if err != nil {
		return err
	}

	err = s.postRepo.DeleteTagFromPost(post, tag)
	if err != nil {
		return err
	}

	return nil
}

func (s *postService) CreatePostWithTags(post *model.Post, tagNames []string) error {
	err := s.postRepo.CreatePost(post)
	if err != nil {
		return err
	}
	for _, tagName := range tagNames {
		tag := model.Tag{Name: tagName}
		s.postRepo.AppendTagToPost(*post, tag)
		// post.Tags = append(post.Tags, tag)
	}
	return nil
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

func (s *postService) GetPostByID(id uint) (model.Post, error) {
	return s.postRepo.GetPostByID(id)
}

func (s *postService) GetPosts(perPage int, page int) ([]model.Post, int64, error){
	
	posts, err := s.postRepo.FindPosts(perPage, page)
	if err != nil {
		return posts, 0, err
	}

	totalData, err := s.postRepo.CountAllPost()
	if err != nil {
		return posts, totalData, err
	}

	return posts, totalData, nil
}
