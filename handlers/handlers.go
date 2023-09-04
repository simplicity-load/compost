package handlers

import (
	"compost/database"
	"compost/models"
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

func Account(c *fiber.Ctx) error {
	session, err := loggedInSession(c)
	if err != nil {
		return err
	}

	userId := session.Get("u").(int)
	username, err := db.GetUserById(userId)
	if err != nil {
		// User does not exist
		return fiber.ErrNotFound
	}
	log.Printf("User %v, logged in", username)

	req := struct {
		Username string `json:"username"`
	}{}
	req.Username = username

	err = extendAndSaveSession(session)
	if err != nil {
		return err
	}
	return c.JSON(req)
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

// The two functions below are expensive since actions on store are somehow expensive?
func loggedInSession(c *fiber.Ctx) (*session.Session, error) {
	session, err := store.Get(c) // Expensive ~ 1ms - 1.5ms
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
	// bench := time.Now()
	session.SetExpiry(time.Minute * database.COOKIE_LIFESPAN)
	// log.Printf(">> Time after setExpiry: %v", time.Now().Sub(bench))
	// bench = time.Now()
	if err := session.Save(); err != nil { // Super Expensive ~ 9ms
		log.Printf("Failed to save cookie %v", session)
		return fiber.ErrInternalServerError
	}
	// log.Printf(">> Time after Save: %v", time.Now().Sub(bench))
	return nil
}

func GetTasks(c *fiber.Ctx) error {
	// bench := time.Now()
	session, err := loggedInSession(c)
	if err != nil {
		return err
	}
	// log.Printf("Time after loggedInSession: %v", time.Now().Sub(bench))
	// bench = time.Now()

	userId := session.Get("u").(int)
	// log.Printf("Time after Get(\"u\"): %v", time.Now().Sub(bench))
	// bench = time.Now()
	tasks, err := db.GetAllTasksForUser(userId)
	if err != nil {
		log.Printf("Failed getting tasks for user: %v", userId)
		return fiber.ErrInternalServerError
	}
	// log.Printf("Time after getTasks: %v", time.Now().Sub(bench))
	// bench = time.Now()
	err = extendAndSaveSession(session)
	if err != nil {
		return err
	}
	// log.Printf("Time after getTasks: %v", time.Now().Sub(bench))
	// bench = time.Now()
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
