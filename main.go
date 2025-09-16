package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// User represents a user in our system
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

// Response represents a standard API response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    int         `json:"code,omitempty"`
}

// In-memory storage
var users = []User{
	{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30},
	{ID: 2, Name: "Jane Smith", Email: "jane@example.com", Age: 25},
	{ID: 3, Name: "Bob Wilson", Email: "bob@example.com", Age: 35},
}
var nextID = 4

func main() {
	router := gin.Default()

	// GET /users - Get all users
	router.GET("/users", getAllUsers)

	// GET /users/:id - Get user by ID
	router.GET("/users/:id", getUserByID)

	// POST /users - Create new user
	router.POST("/users", createUser)

	// PUT /users/:id - Update user
	router.PUT("/users/:id", updateUser)

	// DELETE /users/:id - Delete user
	router.DELETE("/users/:id", deleteUser)

	// GET /users/search - Search users by name
	router.GET("/users/search", searchUsers)

	router.Run(":8080")
}

// TODO: Implement handler functions

// getAllUsers handles GET /users
func getAllUsers(c *gin.Context) {
	req := Response{
		Success: true,
		Data:    users,
		Message: "successfully retrieved all users",
	}
	c.JSON(http.StatusOK, req)
}

// getUserByID handles GET /users/:id
func getUserByID(c *gin.Context) {
	// 1. 解析ID参数并处理格式错误
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid user ID format",
			Code:    400,
		})
		return
	}

	// 3. 根据查找结果返回响应
	if foundUser, index := findUserByID(id); index != -1 {
		c.JSON(http.StatusOK, Response{
			Success: true,
			Data:    foundUser,
			Message: "user found successfully",
		})
	} else {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "user not found",
			Code:    404,
		})
	}
}

// createUser handles POST /users
func createUser(c *gin.Context) {
	var userInput struct {
		Name  string `json:"name" binding:"required"`
		Email string `json:"email" binding:"required,email"`
		Age   int    `json:"age" binding:"omitempty,min=0,max=150"`
	}

	// 解析传入的json
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	for _, em := range users {
		if em.Email == userInput.Email {
			c.JSON(http.StatusConflict, Response{
				Success: false,
				Error:   "Email already in users",
				Code:    http.StatusConflict,
			})
			return
		}
	}

	newUser := User{
		ID:    nextID,
		Name:  userInput.Name,
		Email: userInput.Email,
		Age:   userInput.Age,
	}
	users = append(users, newUser)
	nextID++

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    newUser,
		Message: "User Create successfully",
		Code:    http.StatusCreated,
	})
}

// updateUser handles PUT /users/:id
func updateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid user ID format",
			Code:    400,
		})
		return
	}

	user, index := findUserByID(id)
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "user not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	var updateData struct {
		Name  string `json:"name" binding:"omitempty,required"`
		Email string `json:"email" binding:"omitempty,required,email"`
		Age   int    `json:"age" binding:"omitempty,min=0,max=150"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid update data: " + err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}
	if updateData.Email != "" && updateData.Email != user.Email {
		for _, u := range users {
			if u.Email == updateData.Email {
				c.JSON(http.StatusConflict, Response{
					Success: false,
					Error:   "email already in use",
					Code:    http.StatusConflict,
				})
				return
			}
		}
	}

	if updateData.Name != "" {
		user.Name = updateData.Name
	}
	if updateData.Email != "" {
		user.Email = updateData.Email
	}
	if updateData.Age != 0 {
		user.Age = updateData.Age
	}

	users[index] = *user
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    user,
		Message: "user updated successfully",
		Code:    http.StatusOK,
	})
}

// deleteUser handles DELETE /users/:id
func deleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "invalid user ID format",
			Code:    400,
		})
		return
	}

	_, index := findUserByID(id)
	if index == -1 {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "user not found",
			Code:    http.StatusNotFound,
		})
		return
	}

	users = append(users[:index], users[index+1:]...)

	// 4. 返回删除成功的响应
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "user deleted successfully",
		Code:    http.StatusOK,
	})
}

// searchUsers handles GET /users/search?name=value
func searchUsers(c *gin.Context) {
	queryLower := strings.ToLower(c.Query("name"))

	var matchUser []User
	for _, user := range users {
		// 将用户名转换为小写后判断是否包含查询字符串
		userNameLower := strings.ToLower(user.Name)
		if strings.Contains(userNameLower, queryLower) {
			matchUser = append(matchUser, user)
		}
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    matchUser,
		Message: fmt.Sprintf("found %d matching users", len(matchUser)),
		Code:    http.StatusOK,
	})
}

// Helper function to find user by ID
func findUserByID(id int) (*User, int) {
	var foundUser *User
	for i, user := range users {
		if user.ID == id {
			foundUser = &users[i]
			return foundUser, i
		}
	}
	return nil, -1
}

// Helper function to validate user data
func validateUser(user User) error {
	if strings.TrimSpace(user.Name) == "" {
		return fmt.Errorf("name is required")
	}

	if strings.TrimSpace(user.Email) == "" {
		return fmt.Errorf("email is required")
	}

	if !strings.Contains(user.Email, "@") {
		return fmt.Errorf("invalid email format (missing @)")
	}

	return nil
}
