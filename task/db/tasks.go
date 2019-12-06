package db

import (
	"encoding/binary"
	"github.com/boltdb/bolt"
	"time"
)

var taskBucket = []byte("tasks")

var db *bolt.DB

type Task struct {
	Key   int
	Value string
}

func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		if err != nil {
			return err
		}

		_, err = tx.CreateBucketIfNotExists(getCompletedTasksBucket())
		if err != nil {
			return err
		}
		return nil
	})
}

func CreateTask(task string) (int, error) {
	var id int
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		id64, _ := b.NextSequence()
		id = int(id64)
		key := itob(id)
		return b.Put(key, []byte(task))
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

func CompleteTask(task Task) error {
	return db.Update(func(tx *bolt.Tx) error {
		err := DeleteTask(task.Key, tx)
		if err != nil {
			return err
		}

		b := tx.Bucket(getCompletedTasksBucket())
		err = b.Put(itob(task.Key), []byte(task.Value))
		if err != nil {
			return err
		}

		return nil
	})
}

func CompletedTasks() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(getCompletedTasksBucket())
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, Task{
				Key:   btoi(k),
				Value: string(v),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func AllTasks() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, Task{
				Key:   btoi(k),
				Value: string(v),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func DeleteTask(id int, tx *bolt.Tx) error {
	if tx == nil {
		return db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(taskBucket)
			return b.Delete(itob(id))
		})
	}
	b := tx.Bucket(taskBucket)
	return b.Delete(itob(id))
}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}

func getCompletedTasksBucket() []byte {
	currentDay := time.Now().Format("2006-01-02")
	return []byte(currentDay)
}
