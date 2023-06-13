package handlers

import (
	"fiber-proj1/database"
	"fiber-proj1/models"
	"strconv"

	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var db *database.DBAct
var store *session.Store

func InitDatasources() error {
	var err error
	bench := time.Now()

	if db, err = database.Connect(); err != nil {
		log.Println("Failed to connect database")
		return err
	}
	if store, err = database.CreateSessions(); err != nil {
		log.Println("Failed to connect to sessions")
		return err
	}
	log.Println(bench.Sub(time.Now()))
	return nil
}

func Login(c *fiber.Ctx) error {
	req := struct {
		Username string `json:"username"`
		// Password string `json:"password"`
	}{}
	if err := c.BodyParser(&req); err != nil {
		return fiber.ErrBadRequest
	}

	userId, err := db.GetUserByString(req.Username)
	if err != nil {
		// User does not exist
		return fiber.ErrNotFound
	}

	// CHECK FOR PASSWORD
	// ---
	//

	// Get value
	session, err := store.Get(c)
	if err != nil {
		log.Printf("Failed to get cookie %v", session)
		return fiber.ErrInternalServerError
	}
	session.Set("u", userId)

	if err := session.Save(); err != nil {
		log.Printf("Failed to save cookie %v", session)
		return fiber.ErrInternalServerError
	}
	return c.JSON(nil)
}

func loggedInSession(c *fiber.Ctx) (*session.Session, error) {
	session, err := store.Get(c)
	if err != nil {
		return nil, fiber.ErrInternalServerError
	}
	if session.Fresh() {
		// Not authenticated yet
		return nil, fiber.ErrForbidden
	}
	return session, nil
}

func extendAndSaveSession(session *session.Session) error {
	session.SetExpiry(time.Minute * database.COOKIE_LIFESPAN)
	if err := session.Save(); err != nil {
		log.Printf("Failed to save cookie %v", session)
		return fiber.ErrInternalServerError
	}
	return nil
}

func GetTasks(c *fiber.Ctx) error {
	session, err := loggedInSession(c)
	if err != nil {
		return err
	}

	userId := session.Get("u").(int)
	tasks, err := db.GetAllTasksForUser(userId)
	if err != nil {
		log.Printf("Failed getting tasks for user: %v", userId)
		return fiber.ErrInternalServerError
	}
	err = extendAndSaveSession(session)
	if err != nil {
		return err
	}
	return c.JSON(tasks)
}

func SaveOrEditTask(c *fiber.Ctx) error {
	session, err := loggedInSession(c)
	if err != nil {
		return err
	}

	var task models.Task
	if err = c.BodyParser(&task); err != nil {
		return fiber.ErrBadRequest
	}

	task.UserId = session.Get("u").(int)

	// Check if recieved values are within criteria, e.g.
	// task.Status should have only some values

	err = db.SaveOrEditTask(task)
	if err != nil {
		return err
	}

	err = extendAndSaveSession(session)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(nil)
}

func DeleteTask(c *fiber.Ctx) error {
	session, err := loggedInSession(c)
	if err != nil {
		return err
	}

	id, e := strconv.Atoi(c.Params("id"))
	log.Printf("id: %v", id)
	if e != nil {
		return fiber.ErrBadRequest
	}

	userId := session.Get("u").(int)

	err = db.DeleteTask(id, userId)
	if err != nil {
		return err
	}

	err = extendAndSaveSession(session)
	if err != nil {
		return err
	}
	return c.JSON(nil)
}

func GetTask(c *fiber.Ctx) error {
	session, err := loggedInSession(c)
	if err != nil {
		return err
	}

	id, e := strconv.Atoi(c.Params("id"))
	if e != nil {
		return fiber.ErrBadRequest
	}

	userId := session.Get("u").(int)

	task, err := db.GetTask(id, userId)
	if err != nil {
		return err
	}

	err = extendAndSaveSession(session)
	if err != nil {
		return err
	}
	return c.JSON(task)
}

// NotFound returns custom 404 page
func NotFound(c *fiber.Ctx) error {
	return c.Status(404).SendFile("./static/private/404.html")
}
