package main

import (
	"net/http"
	"strconv"
	"strings"

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
		errorHandler(ctx, err, 1)
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
		errorHandler(ctx, err, 1)
	}
}

func PostUser(ctx *gin.Context) {
	var user proto.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	if response, err := Client.AddUser(ctx, &user); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, err, 1)
	}
}

func PatchUser(ctx *gin.Context) {
	var user proto.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		errorHandler(ctx, err, 1)
	}
}

func DeleteUser(ctx *gin.Context) {
	var user proto.User

	if err := ctx.BindJSON(&user); err != nil {
        ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		errorHandler(ctx, err, 1)
	}
}

func AllProduct(ctx *gin.Context) {
	req := &proto.EmptyStruct{}
	if response, err := Client.AllProduct(ctx, req); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, err, 2)
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
		ctx.IndentedJSON(http.StatusInternalServerError, response)
	} else {	
		errorHandler(ctx, err, 2)
	}
}

func PostProduct(ctx *gin.Context) {
	var adminProduct proto.AdminProduct

	if err := ctx.BindJSON(&adminProduct); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	if response, err := Client.AddProduct(ctx, &adminProduct); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, err, 2)
	}
}

func PatchProduct(ctx *gin.Context) {
	var adminProduct proto.AdminProduct

	if err := ctx.BindJSON(&adminProduct); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Parameter Id", "error": err.Error()})
		return
	}

	adminProduct.Id = int64(id)

	if response, err := Client.UpdateProduct(ctx, &adminProduct); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, err, 2)
	}
}

func DeleteProduct(ctx *gin.Context) {
	var adminProduct proto.AdminProduct

	if err := ctx.BindJSON(&adminProduct); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.IndentedJSON(http.StatusBadRequest, gin.H{"status": "Invalid Parameter Id", "error": err.Error()})
		return
	}

	adminProduct.Id = int64(id)

	if response, err := Client.DeleteProduct(ctx, &adminProduct); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, err, 2)
	}
}

func AllOrder(ctx *gin.Context) {
	var user proto.User

	if err := ctx.BindJSON(&user); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	if response, err := Client.AllOrder(ctx, &user); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {
		errorHandler(ctx, err, 3)
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
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	order.Id = id

	if response, err := Client.OneOrder(ctx, &order); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, err, 3)
	}
}

func PostOrder(ctx *gin.Context) {
	var order proto.Order

	if err := ctx.BindJSON(&order); err != nil {
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	if response, err := Client.AddOrder(ctx, &order); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, err, 3)
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
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	order.Id = id

	if response, err := Client.UpdateOrder(ctx, &order); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, err, 3)
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
		ctx.IndentedJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

	order.Id = id

	if response, err := Client.DeleteOrder(ctx, &order); err == nil {
		ctx.IndentedJSON(http.StatusOK, response)
	} else {	
		errorHandler(ctx, err, 3)
	}
}

// ERROR HANDLER

// type : 1 for user, 2 for product, 3 for order
func errorHandler(ctx *gin.Context, err error, errType int) {
	switch errType{
	case 1:
		if strings.Contains(err.Error(), "wrong password for user"){
			ctx.IndentedJSON(http.StatusForbidden, err.Error())
			return
		} else if strings.Contains(err.Error(), "code = NotFound") {
			ctx.IndentedJSON(http.StatusNotFound, err.Error())
			return
		} else if strings.Contains(err.Error(), "please include password in request") {
			ctx.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		} else if strings.Contains(err.Error(), "code = AlreadyExists") {
			ctx.IndentedJSON(http.StatusConflict, err.Error())
			return
		}
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		return


	case 2:
		if strings.Contains(err.Error(), "wrong password for user"){
			ctx.IndentedJSON(http.StatusForbidden, err.Error())
			return
		} else if strings.Contains(err.Error(), "code = NotFound") {
			ctx.IndentedJSON(http.StatusNotFound, err.Error())
			return
		} else if strings.Contains(err.Error(), "please include password in request") {
			ctx.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		} else if strings.Contains(err.Error(), "code = AlreadyExists") {
			ctx.IndentedJSON(http.StatusConflict, err.Error())
			return
		} else if strings.Contains(err.Error(), "code = PermissionDenied") {
			ctx.IndentedJSON(http.StatusForbidden, err.Error())
			return
		}
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		return

	case 3:
		if strings.Contains(err.Error(), "wrong password for user"){
			ctx.IndentedJSON(http.StatusForbidden, err.Error())
			return
		} else if strings.Contains(err.Error(), "code = NotFound") {
			ctx.IndentedJSON(http.StatusNotFound, err.Error())
			return
		} else if strings.Contains(err.Error(), "please include password in request") {
			ctx.IndentedJSON(http.StatusBadRequest, err.Error())
			return
		}
		ctx.IndentedJSON(http.StatusInternalServerError, err.Error())
		return
	}
}