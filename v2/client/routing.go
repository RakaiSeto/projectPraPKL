package main

import (
	"net/http"
	"strconv"

	proto "github.com/RakaiSeto/projectPraPKL/v2/proto"
	"github.com/gin-gonic/gin"
)

func Tes(ctx *gin.Context) {
	req := &proto.EmptyStruct{}
	if response, err := Client.Tes(ctx, req); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func AllUser(ctx *gin.Context) {
	req := &proto.EmptyStruct{}
	if response, err := Client.AllUser(ctx, req); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
	}
}

func OneUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Parameter Id", "error": err.Error()})
		return
	}

	req := &proto.Id{Id: int64(id)}
	if response, err := Client.OneUser(ctx, req); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func PostUser(ctx *gin.Context) {
	var user proto.User

	if err := ctx.BindJSON(&user); err != nil {
        return
    }

	if response, err := Client.AddUser(ctx, &user); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		ctx.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
	}
}

func PatchUser(ctx *gin.Context) {
	var user proto.User

	if err := ctx.BindJSON(&user); err != nil {
        return
    }

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Parameter Id", "error": err.Error()})
		return
	}

	user.Id = int64(id)

	if response, err := Client.UpdateUser(ctx, &user); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
	}
}

func DeleteUser(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Parameter Id", "error": err.Error()})
		return
	}

	if response, err := Client.DeleteUser(ctx, &proto.Id{Id: int64(id)}); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
	}
}

