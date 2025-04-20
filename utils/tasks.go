package utils

import (
	"encoding/json"
	"go-web-server/models"
	"os"
)

var Tasks []models.Task
var TaskID int = 1

const dataFile = "tasks.json"

func LoadTasksFromFile() error {
	file, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			Tasks = []models.Task{}
			return nil
		}
		return err
	}
	return json.Unmarshal(file, &Tasks)
}

func SaveTasksToFile() error {
	data, err := json.MarshalIndent(Tasks, " ", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, data, 0644)
}

func UpdateTaskID() {
	maxId := 0
	for _, task := range Tasks {
		if task.ID > maxId {
			maxId = task.ID
		}
	}
	TaskID = maxId + 1
}
