package router
import (
	"github.com/fiber/src/handlers"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GroupRouterUser(v1 fiber.Router, db *gorm.DB) {
	v1.Get("/user/:name",func(c *fiber.Ctx) error {
		return handlers.GetUserTittleByUserID(c,db)
	})

}

func GroupRouterPosts(v1 fiber.Router, db *gorm.DB){
	v1.Post("/post/create", func(c *fiber.Ctx) error {
		return handlers.CreatePost(c,db)
	}) 
	
	v1.Get("/post/all",func(c *fiber.Ctx) error {
		return handlers.GetPost(c,db)
	})
}

func GroupRouterComments(v1 fiber.Router, db *gorm.DB){
	v1.Post("/comments/create", func(c *fiber.Ctx) error {
		return handlers.CreateComments(c,db)
	}) 
	
	v1.Get("/comments/all/:postID", func(c *fiber.Ctx) error {
		return handlers.GetCommentByPostID(c,db)
	}) 

	v1.Put("/comments/update/:ID", func(c *fiber.Ctx) error {
		return handlers.UpdateCommentByID(c,db)
	}) 

	v1.Delete("/comments/delete/:ID", func(c *fiber.Ctx) error {
		return handlers.DeleteCommentByID(c,db)
	}) 


}