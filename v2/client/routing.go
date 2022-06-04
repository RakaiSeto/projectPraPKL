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
		ctx.IndentedJSON(http.StatusOK, gin.H{"error": err.Error()})
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
		ctx.IndentedJSON(http.StatusOK, gin.H{"error": err.Error()})
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
	var user proto.User

	if err := ctx.BindJSON(&user); err != nil {
        return
    }
	
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Parameter Id", "error": err.Error()})
		return
	}

	user.Id = id

	if response, err := Client.DeleteUser(ctx, &user); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "user not found or password false"})
	}
}

func AllProduct(ctx *gin.Context) {
	req := &proto.EmptyStruct{}
	if response, err := Client.AllProduct(ctx, req); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
	}
}

func OneProduct(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Parameter Id", "error": err.Error()})	
		return
	}

	req := &proto.Id{Id: id}
	if response, err := Client.OneProduct(ctx, req); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func AllOrder(ctx *gin.Context) {
	var user proto.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.IndentedJSON(http.StatusOK, gin.H{"error": err.Error()})
        return
    }

	if response, err := Client.AllOrder(ctx, &user); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
	}
}

func OneOrder(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Parameter Id", "error": err.Error()})	
		return
	}

	var order proto.Order

	if err := ctx.BindJSON(&order); err != nil {
		ctx.IndentedJSON(http.StatusOK, gin.H{"error": err.Error()})
        return
    }

	order.Id = id

	if response, err := Client.OneOrder(ctx, &order); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err})
	}
}

func PostOrder(ctx *gin.Context) {
	var order proto.Order

	if err := ctx.BindJSON(&order); err != nil {
		ctx.IndentedJSON(http.StatusOK, gin.H{"error": err.Error()})
        return
    }

	if response, err := Client.AddOrder(ctx, &order); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		ctx.JSON(http.StatusConflict, gin.H{"error": err})
	}
}

func PatchOrder(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Parameter Id", "error": err.Error()})	
		return
	}

	var order proto.Order

	if err := ctx.BindJSON(&order); err != nil {
		ctx.IndentedJSON(http.StatusOK, gin.H{"error": err.Error()})
        return
    }

	order.Id = id

	if response, err := Client.UpdateOrder(ctx, &order); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		ctx.JSON(http.StatusConflict, gin.H{"error": "order doesn't exist or password wrong"})
	}
}

func DeleteOrder(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Parameter Id", "error": err.Error()})	
		return
	}

	var order proto.Order

	if err := ctx.BindJSON(&order); err != nil {
		ctx.IndentedJSON(http.StatusOK, gin.H{"error": err.Error()})
        return
    }

	order.Id = id

	if response, err := Client.DeleteOrder(ctx, &order); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		ctx.JSON(http.StatusConflict, gin.H{"error": "order doesn't exist or password wrong"})
	}
}