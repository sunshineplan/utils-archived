package mongodb

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Config contains mongodb basic configure.
type Config struct {
	Server     string
	Port       int
	Database   string
	Collection string
	Username   string
	Password   string
}

// Open opens a mongodb database.
func (c *Config) Open() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		fmt.Sprintf("mongodb://%s:%s@%s:%d/%s", c.Username, c.Password, c.Server, c.Port, c.Database)))
	if err != nil {
		return nil, err
	}
	ctx, cancelPing := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelPing()
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Backup backups mongodb database to file.
func (c *Config) Backup(file string) error {
	args := []string{}
	args = append(args, fmt.Sprintf("-h%s:%d", c.Server, c.Port))
	args = append(args, fmt.Sprintf("-d%s", c.Database))
	args = append(args, fmt.Sprintf("-c%s", c.Collection))
	args = append(args, fmt.Sprintf("-u%s", c.Username))
	args = append(args, fmt.Sprintf("-p%s", c.Password))
	args = append(args, "--gzip")
	args = append(args, fmt.Sprintf("--archive=%s", file))

	command := exec.Command("mongodump", args...)
	var stderr bytes.Buffer
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("Failed to backup database: %s\n%v", stderr.String(), err)
	}
	return nil
}

// Restore restores mongodb database from file.
func (c *Config) Restore(file string) error {
	args := []string{}
	args = append(args, fmt.Sprintf("-h%s:%d", c.Server, c.Port))
	args = append(args, fmt.Sprintf("-d%s", c.Database))
	args = append(args, fmt.Sprintf("-c%s", c.Collection))
	args = append(args, fmt.Sprintf("-u%s", c.Username))
	args = append(args, fmt.Sprintf("-p%s", c.Password))
	args = append(args, "--gzip")
	args = append(args, "--drop")
	args = append(args, fmt.Sprintf("--archive=%s", file))

	command := exec.Command("mongorestore", args...)
	var stderr bytes.Buffer
	command.Stderr = &stderr
	if err := command.Run(); err != nil {
		return fmt.Errorf("Failed to restore database: %s\n%v", stderr.String(), err)
	}
	return nil
}
