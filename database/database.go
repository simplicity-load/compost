package database

import (
	"fiber-proj1/models"

	"database/sql"
	_ "embed"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed task.sql
var db_schema_file string

const db_file_name = "tasks.sqlite3"

type DBAct struct {
	db *sql.DB
	mu sync.Mutex
}

func createIfNotExists(db *sql.DB, file string) error {
	query := file
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

// Connect with database
func Connect() (*DBAct, error) {
	db, err := sql.Open("sqlite3", db_file_name)
	if err != nil {
		return nil, err
	}
	if err := createIfNotExists(db, db_schema_file); err != nil {
		log.Printf("Failed creating schema of database\nErr: %v", err)
		return nil, err
	}
	log.Println("Connected with Database")
	return &DBAct{
		db: db,
	}, nil
}

func (a *DBAct) GetAllTasksForUser(userId int) ([]models.Task, error) {
	tasks := make([]models.Task, 0, 100)
	a.mu.Lock()
	defer a.mu.Unlock()
	rows, err := a.db.Query(get_all_tasks_for_userId, userId)
	defer rows.Close()
	if err != nil {
		log.Printf("Failed to select rows\nErr: %v", err)
		return nil, err
	}

	for rows.Next() {
		task := models.Task{}
		err = rows.Scan(
			&task.Id,
			&task.UserId,
			&task.Title,
			&task.Body,
			&task.Status)
		if err != nil {
			log.Printf("Failed scanning rows\nErr: %v", err)
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (a *DBAct) GetUserByString(username string) (id int, err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	row := a.db.QueryRow(get_userId_from_username, username)
	if err = row.Scan(&id); err == sql.ErrNoRows {
		log.Printf("Id not found %v, %v, %v", id, username, "jnig" == username)
		return 0, err
	}
	return id, nil
}

func (a *DBAct) taskExistsByID(task *models.Task) (bool, error) {
	row := a.db.QueryRow(get_task_by_id, task.Id)
	err := row.Scan(
		&task.Id,
		&task.UserId,
		&task.Title,
		&task.Body,
		&task.Status)
	log.Printf("taskid: %v", task.Id)

	switch {
	case err == sql.ErrNoRows:
		return false, nil
	case err != nil:
		log.Printf("Failed to get count of task, %v", err)
		return false, fiber.ErrInternalServerError
	}
	return true, nil
}

func (a *DBAct) taskExistsByStatus(task *models.Task) bool {
	// TODO change from str "del" to enum or smth
  return task.Status != "del"
}

func (a *DBAct) SaveOrEditTask(task models.Task) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	var retrievedTask models.Task
	retrievedTask.Id = task.Id
	isEditing, err := a.taskExistsByID(&retrievedTask)
	if err != nil {
		return err
	}
	exists := a.taskExistsByStatus(&retrievedTask)
	if !exists {
		return fiber.ErrNotFound
	}

	if isEditing {
		// Editing task
		if retrievedTask.UserId != task.UserId {
			log.Printf("User does not own task: %v", retrievedTask.Id)
			return fiber.ErrForbidden
		}
		_, err = a.db.Exec(
			update_task_by_id,
			task.Title,
			task.Body,
			task.Status,
			task.Id)
		if err != nil {
			return fiber.ErrBadRequest
		}
	} else {
		_, err := a.db.Exec(insert_task,
			task.UserId,
			task.Title,
			task.Body,
			task.Status)
		if err != nil {
			return fiber.ErrBadRequest
		}
	}
	return nil
}

func (a *DBAct) DeleteTask(id int, userId int) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	var retrievedTask models.Task
	retrievedTask.Id = id
	exists, err := a.taskExistsByID(&retrievedTask)
	if err != nil {
		return err
	}
	exists = a.taskExistsByStatus(&retrievedTask)
	if !exists {
		return fiber.ErrNotFound
	}

	if exists {
		if retrievedTask.UserId != userId {
			log.Printf("User does not own task: %v", retrievedTask.Id)
			return fiber.ErrForbidden
		}
		_, err := a.db.Exec(delete_task_by_id, id)
		if err != nil {
			return fiber.ErrInternalServerError
		}
	} else {
		log.Println("3 here")
		return fiber.ErrNotFound
	}
	return nil
}

func (a *DBAct) GetTask(id int, userId int) (*models.Task, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	var retrievedTask models.Task
	retrievedTask.Id = id
	exists, err := a.taskExistsByID(&retrievedTask)
	if err != nil {
		return nil, err
	}
	exists = a.taskExistsByStatus(&retrievedTask)
	if !exists {
		return nil, fiber.ErrNotFound
	}

	if exists {
		if retrievedTask.UserId != userId {
			log.Printf("User does not own task: %v", retrievedTask.Id)
			return nil, fiber.ErrForbidden
		}
	} else {
		return nil, fiber.ErrNotFound
	}
	return &retrievedTask, nil
}
