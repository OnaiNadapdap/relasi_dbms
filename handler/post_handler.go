package handler

import (
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onainadapdap1/kampus_tutor/coba_many_to_many/helper"
	"github.com/onainadapdap1/kampus_tutor/coba_many_to_many/model"
	"github.com/onainadapdap1/kampus_tutor/coba_many_to_many/service"
)

type PostHandler interface {
	CreatePostWithTags(c *gin.Context)
	DeleteTagFromPost(c *gin.Context)
	UpdatePostWithTags(c *gin.Context)
	GetAllPosts(c *gin.Context)
	GetPostByTagID(c *gin.Context)
	AppendTagToPost(c *gin.Context)
}

type postHandler struct {
	postService service.PostService
}

func NewPostHandler(postService service.PostService) PostHandler {
	return &postHandler{postService: postService}
}

func (h *postHandler) AppendTagToPost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tagNames []string
	var inputData model.AppendTagToPost
	if err := c.ShouldBindJSON(&inputData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(inputData.TagNames) > 0 {
		tagNames = append(tagNames, inputData.TagNames...)
	}

	if err := h.postService.AppendTagToPost(uint(postID), tagNames); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success add tag to post"})
}

func (h *postHandler) GetPostByTagID(c *gin.Context) {
	tagID, err := strconv.Atoi(c.Param("tagID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	posts, err := h.postService.GetPostByTagID(uint(tagID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// func (h *postHandler) GetPostByTagID(c *gin.Context) {
//     tagID, err := strconv.Atoi(c.Param("tagID"))
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     tag, err := h.postService.GetPostByTagID(uint(tagID))
//     if err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     c.JSON(http.StatusOK, tag)
// }

func (h *postHandler) DeleteTagFromPost(c *gin.Context) {
	postID, err := strconv.Atoi(c.Param("postID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tagID, err := strconv.Atoi(c.Param("tagID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.postService.DeleteTagFromPost(uint(tagID), uint(postID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "success delete tag from post"})
}

func (h *postHandler) CreatePostWithTags(c *gin.Context) {
	var post model.Post

	var tagNames []string

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(post.TagNames) > 0 {
		tagNames = append(tagNames, post.TagNames...)
	}

	if err := h.postService.CreatePostWithTags(&post, tagNames); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *postHandler) UpdatePostWithTags(c *gin.Context) {
	id := c.Param("id")
	var updatedPost model.Post

	if err := c.ShouldBindJSON(&updatedPost); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	postID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.postService.UpdatePostByIDWithTags(uint(postID), &updatedPost); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	c.JSON(http.StatusOK, updatedPost)
}

func (h *postHandler) GetAllPosts(c *gin.Context) {
	perPage, _ := strconv.Atoi(c.DefaultQuery("perPage", "10"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	if perPage < 1 || perPage > 10 {
		perPage = 10
	}

	if page < 1 {
		page = 1
	}
	posts, totalData, err := h.postService.GetPosts(perPage, page)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "failed",
		})
		return
	}

	// totalPage := (totalData / int64(perPage)) + 1
	totalPage := math.Ceil(float64(totalData) / float64(perPage))

	response := helper.APIResponse("success to fetch data", http.StatusOK, "success", posts, perPage, page, totalData, int(totalPage))

	c.JSON(http.StatusOK, response)
}

// func (h *PostHandler) CreatePost(c *gin.Context) {
//     var post model.Post
//     var tagNames []string

//     if err := c.ShouldBindJSON(&post); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

// 	if len(post.TagsName) > 0 {
// 		tagNames = append(tagNames, post.TagsName...)
// 	}

//     if err := h.PostService.CreatePostWithTags(&post, tagNames); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
//         return
//     }

//     c.JSON(http.StatusCreated, post)
// }

// func (h *PostHandler) UpdatePost(c *gin.Context) {
//     var post model.Post
//     var tagNames []string

//     if err := c.ShouldBindJSON(&post); err != nil {
//         c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//         return
//     }

//     if len(post.TagsName) > 0 {
// 		tagNames = append(tagNames, post.TagsName...)
// 	}

//     if err := h.PostService.UpdatePostWithTags(&post, tagNames); err != nil {
//         c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
//         return
//     }

//     c.JSON(http.StatusCreated, post)
// }
