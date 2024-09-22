package handlers

import (
	"fmt"
	"github.com/fiber/src/models"
	"github.com/fiber/src/utils"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"time"
)

type User struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
}

type PostResponse struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Content       string    `json:"content"`
	ImageURL      string    `json:"image_url"`
	LikesCount    int       `json:"likes_count"`
	CommentsCount int       `json:"comments_count"`
	User          User      `json:"user"`
}

func GetUser(c *fiber.Ctx, db *gorm.DB) error {
	id, err := utils.Decoded(c.Cookies("auth"))

	if err != nil {
		fmt.Println(err)
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	var user User
	if err := db.Where("id = ?", id).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}
	return c.Status(200).JSON(user)
}

func GetUserTittleByUserID(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("name")
	if id == "" {
		return c.Status(404).JSON(fiber.Map{"error": "Params is empty"})
	}

	var user models.User
	if err := db.Preload("Followers").Preload("Following").Where("user_id = ?", id).First(&user).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	followersCount := len(user.Followers)
	followingCount := len(user.Following)

	response := fiber.Map{
		"name":           user.Name,
		"email":          user.Email,
		"followersCount": followersCount,
		"followingCount": followingCount,
	}

	return c.Status(200).JSON(response)
}

func CreatePost(c *fiber.Ctx, db *gorm.DB) error {
	var post models.Post
	id, err := utils.Decoded(c.Cookies("auth"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create post",
		})
	}
	post.UserID = id

	if err := c.BodyParser(&post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	if post.Content == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Content is required",
		})
	}
	if err := db.Create(&post).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create post",
		})
	}

	// คืนค่า post ที่เพิ่มไปแล้ว
	return c.Status(fiber.StatusCreated).SendString("create post successful!")
}

func GetPost(c *fiber.Ctx, db *gorm.DB) error {
	var posts []models.Post
	var responsePosts []PostResponse

	// ดึงข้อมูลโพสต์ทั้งหมดพร้อมกับข้อมูลผู้ใช้ รวมถึง Likes และ Comments
	if err := db.Preload("User").
		Preload("Likes").
		Preload("Comments").
		Find(&posts).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch posts",
		})
	}

	for _, post := range posts {
		// นับจำนวน Likes และ Comments
		likesCount := len(post.Likes)
		commentsCount := len(post.Comments)

		// เพิ่มข้อมูลลงใน PostResponse
		responsePost := PostResponse{
			ID:        post.ID,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			Content:   post.Content,
			ImageURL:  post.ImageURL,
			User: User{
				UserID: post.User.UserID,
				Name:   post.User.Name,
				Email:  post.User.Email,
			},
			LikesCount:    likesCount,
			CommentsCount: commentsCount,
		}
		responsePosts = append(responsePosts, responsePost)
	}

	// ส่งกลับข้อมูลโพสต์ทั้งหมดในรูปแบบที่กำหนด
	return c.JSON(responsePosts)
}

func CreateComments(c *fiber.Ctx, db *gorm.DB) error {
	var comments models.Comment
	id, err := utils.Decoded(c.Cookies("auth"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create post",
		})
	}
	comments.UserID = id

	if err := c.BodyParser(&comments); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	if comments.CommentText == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Content is required",
		})
	}
	if err := db.Create(&comments).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create post",
		})
	}

	// คืนค่า post ที่เพิ่มไปแล้ว
	return c.Status(fiber.StatusCreated).SendString("create comments successful!")
}

func GetCommentByPostID(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("postID")
	if id == "" {
		return c.Status(404).JSON(fiber.Map{"error": "Params is empty"})
	}

	var comments []models.Comment
	if err := db.Preload("User").Where("post_id = ?", id).Find(&comments).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Comments not found"})
	}

	// สร้าง response เพื่อส่งกลับไป
	type CommentResponse struct {
		ID          uint      `json:"id"`
		CommentText string    `json:"comment_text"`
		UserName    string    `json:"user_name"`
		CreatedAt   time.Time `json:"createdAt"`
	}

	if len(comments) == 0 {
		return c.Status(200).JSON([]CommentResponse{})
	}

	var response []CommentResponse
	for _, comment := range comments {
		response = append(response, CommentResponse{
			ID:          comment.ID,
			CommentText: comment.CommentText,
			UserName:    comment.User.Name,
			CreatedAt:   comment.CreatedAt,
		})
	}

	return c.Status(200).JSON(response)
}

func UpdateCommentByID(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("ID")
	if id == "" {
		return c.Status(404).JSON(fiber.Map{"error": "Params is empty"})
	}

	// รับข้อมูลที่ส่งมาจาก client
	var updateData struct {
		CommentText string `json:"comment_text"`
	}

	// ตรวจสอบว่าข้อมูลที่ส่งมานั้นถูกต้อง
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input data"})
	}

	// หา comment ที่ต้องการอัปเดต
	var comment models.Comment
	if err := db.First(&comment, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Comment not found"})
	}

	// อัปเดตข้อมูล comment
	comment.CommentText = updateData.CommentText

	// บันทึกการเปลี่ยนแปลงในฐานข้อมูล
	if err := db.Save(&comment).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update comment"})
	}

	return c.Status(200).SendString("Update successful!")
}

func DeleteCommentByID(c *fiber.Ctx, db *gorm.DB) error {
	id := c.Params("ID")
	if id == "" {
		return c.Status(404).JSON(fiber.Map{"error": "Params is empty"})
	}

	// หา comment ที่ต้องการลบ
	var comment models.Comment
	if err := db.First(&comment, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Comment not found"})
	}

	// ลบ comment
	if err := db.Delete(&comment).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete comment"})
	}

	return c.Status(200).SendString("Delete successful!")
}