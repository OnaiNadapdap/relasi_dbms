package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/onainadapdap1/kampus_tutor/coba_many_to_many/model"
	"github.com/onainadapdap1/kampus_tutor/coba_many_to_many/service"
)

type PostHandler interface {
    CreatePostWithTags(c *gin.Context)
    UpdatePostWithTags(c *gin.Context) 
    GetAllPosts(c *gin.Context) 
}

type postHandler struct {
    postService service.PostService
}

func NewPostHandler(postService service.PostService) PostHandler {
    return &postHandler{postService: postService}
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
    posts, err := h.postService.GetPosts()
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "failed",
        })
        return
    }

    c.JSON(http.StatusOK, posts)
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
