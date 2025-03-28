package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/magiconair/properties"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var version string

// fileExistsAndNotEmpty checks if the file exists and is not empty.
func fileExistsAndNotEmpty(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) || info.Size() == 0 {
		return false
	}
	return true
}

func main() {
	// Define command-line flags
	configFile := flag.String("config", "", "Path to property file")
	blobList := flag.String("bloblist", "", "List of blob files to be retrieved")
	blobPath := flag.String("blobpath", "", "Path where to store the blob files")
	showVersion := flag.Bool("version", false, "Display application version")

	flag.Parse()

	// Argument validation
	if *showVersion {
		fmt.Println("GridFS Data Extractor:", version)
		return
	}
	if *configFile == "" || *blobList == "" || *blobPath == "" {
		log.Fatalf("Usage: %s -config <config_file> -bloblist <list_of_blob_files> -blobpath <Stored_blob_path>", os.Args[0])
		return
	}

	// Read configuration
	err := readConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to read configuration: %v", err)
	}

	uri := viper.GetString("MONGO_URI")
	username := viper.GetString("MONGO_USER")
	password := viper.GetString("MONGO_PASS")
	database := viper.GetString("MONGO_DB")
	gridFSPrefix := viper.GetString("MONGO_GRIDFS_PREFIX")

	// Read file names from the provided file
	fileNames, err := readFileNames(*blobList)
	if err != nil {
		log.Fatalf("Failed to read file names: %v", err)
	}

	// Check destination blob path
	err = os.MkdirAll(*blobPath, 0755)
	if err != nil {
		fmt.Println("Error creating directory: %v", err)
		return
	}

	// Create a new client and connect to the server
	clientOptions := options.Client().ApplyURI(uri).SetReadPreference(readpref.Secondary()).SetAuth(options.Credential{
		Username: username,
		Password: password,
	})
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ensure disconnection at the end
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Failed to disconnect MongoDB: %v", err)
		}
	}()

	// Ping the primary
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	// Select the database and GridFS bucket
	db := client.Database(database)
	bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(gridFSPrefix))
	if err != nil {
		log.Fatalf("Failed to create GridFS bucket: %v", err)
	}

	for _, fileName := range fileNames {
		if fileExistsAndNotEmpty(fileName) {
			fmt.Printf("File %s already exists locally or is empty.\n", fileName)
		} else {
			// Open a download stream by file name
			downloadStream, err := bucket.OpenDownloadStreamByName(fileName)
			if err != nil {
				log.Printf("Failed to open download stream for file %v: %v", fileName, err)
				continue
			}

			// Read the file data
			data, err := io.ReadAll(downloadStream)
			if err != nil {
				log.Printf("Failed to read data from download stream: %v", err)
				downloadStream.Close()
				continue
			}
			downloadStream.Close()

			// Save the file to disk
			filePath := filepath.Join(*blobPath, fileName)
			err = os.WriteFile(filePath, data, 0644)
			if err != nil {
				log.Printf("Failed to write file %v to disk: %v", filePath, err)
				continue
			}

			fmt.Printf("Downloaded file from gridfs to local: %v\n", fileName)
		}
	}
}

// readConfig reads the configuration from the specified file.
func readConfig(filename string) error {
	p, err := properties.LoadFile(filename, properties.UTF8)
	if err != nil {
		return err
	}

	// Iterate through the properties and set them in Viper
	for _, key := range p.Keys() {
		value, ok := p.Get(key)
		if ok {
			viper.Set(key, value)
		}
	}
	// Return nil to indicate successful completion
	return nil
}

// readFileNames reads file names from the given file.
func readFileNames(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var names []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			names = append(names, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return names, nil
}
