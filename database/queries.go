package database

const (
	insert_task              = `INSERT INTO tasks(user_id, title, body, status) VALUES(?,?,?,?)`
	get_all_tasks_for_userId = `SELECT * FROM tasks WHERE user_id=?`
	get_username_from_userId = `SELECT name FROM users WHERE id=?`
	get_userId_from_username = `SELECT id FROM users WHERE name=?`
	get_task_by_id           = `SELECT * FROM tasks WHERE id=?`
	update_task_by_id        = `UPDATE tasks SET title=?, body=?, status=? WHERE id=?`
	delete_task_by_id        = `UPDATE tasks SET status="del" WHERE id=?`
)
