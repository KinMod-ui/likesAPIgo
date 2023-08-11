package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	db "likesApi/db/sqlc"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createLikeParams struct {
	UserID    json.Number `json:"user_id" binding:"required"`
	ContentID json.Number `json:"content_id" binding:"required"`
}

const (
	SHORT_EXPIRY = 300 * 1e9  // 5 mins
	LONG_EXPIRY  = 3600 * 1e9 // 1 hour
)

// Add like with a given userId and ContentId if not liked already
func (server *Server) addLike(ctx *gin.Context) {

	var req createLikeParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	UserID, err := req.UserID.Int64()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	ContentID, err := req.ContentID.Int64()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	alreadyLiked := server.client.SIsMember(ctx, fmt.Sprintf("user:%d:likes", UserID), ContentID)
	if alreadyLiked.Val() == true {

		log.Println("Taken from cache")
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "The user already likes the content"},
		)
		return
	}

	arg := db.CreateLikeParams{
		UserID:    UserID,
		ContentID: ContentID,
	}

	likeAdd, err := server.store.CreateLike(ctx, arg)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "likes_pkey"` {

			server.client.SAdd(ctx, fmt.Sprintf("user:%d:likes", UserID), ContentID)
			server.client.Expire(ctx, fmt.Sprintf("user:%d:likes", UserID), SHORT_EXPIRY)

			ctx.JSON(
				http.StatusBadRequest,
				gin.H{"error": "The user already likes the content"},
			)
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	server.client.SAdd(ctx, fmt.Sprintf("user:%d:likes", UserID), ContentID)

	server.client.Expire(ctx, fmt.Sprintf("user:%d:likes", UserID), SHORT_EXPIRY)
	server.client.Incr(ctx, fmt.Sprintf("content:%d:likes", ContentID))
	server.client.Expire(ctx, fmt.Sprintf("content:%d:likes", ContentID) , LONG_EXPIRY)

	ctx.JSON(http.StatusOK, likeAdd)
}

type getLikeRequest struct {
	UserID    int64 `uri:"user_id" binding:"required,min=1"`
	ContentID int64 `uri:"content_id" binding:"required,min=1"`
}

// getLike returns true for whether the UserID has liked ContentID
func (server *Server) getLike(ctx *gin.Context) {
	var req getLikeRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	alreadyLiked := server.client.SIsMember(ctx, fmt.Sprintf("user:%d:likes", req.UserID), req.ContentID)
	if alreadyLiked.Val() == true {
		log.Println("taken from cache")
		ctx.JSON(http.StatusOK, gin.H{"liked": true})
	}

	arg := db.GetLikeParams{
		UserID:    req.UserID,
		ContentID: req.ContentID,
	}

	likes, err := server.store.GetLike(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusOK, gin.H{"liked": false})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	server.client.SAdd(ctx, fmt.Sprintf("user:%d:likes", req.UserID), req.ContentID)
	server.client.Expire(ctx, fmt.Sprintf("user:%d:likes", req.UserID), SHORT_EXPIRY)

	ctx.JSON(http.StatusOK, gin.H{"liked": likes.Bool})
}

type getTotalLikesRequest struct {
	ContentID int64 `uri:"content_id" binding:"required,min=1"`
}

func (server *Server) getTotalLikes(ctx *gin.Context) {
	var req getTotalLikesRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	likeCount := server.client.Get(ctx, fmt.Sprintf("content:%d:likes", req.ContentID))
	if likeCount.Val() != "" {
		log.Println("From cache")
		ctx.JSON(http.StatusOK, gin.H{"likes": likeCount.Val()})
		return
	}

	likes, err := server.store.TotalLikesForContent(
		ctx,
		req.ContentID,
	)

	if err != nil {

		if err == sql.ErrNoRows {
			server.client.SetEx(ctx, fmt.Sprintf("content:%d:likes", req.ContentID), 0, LONG_EXPIRY)
			ctx.JSON(http.StatusOK, gin.H{"likes": 0})
			return
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	server.client.SetEx(ctx, fmt.Sprintf("content:%d:likes", req.ContentID), likes, LONG_EXPIRY)
	ctx.JSON(http.StatusOK, gin.H{"likes": likes})

}
