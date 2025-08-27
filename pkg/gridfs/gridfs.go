package gridfs

import (
	"context"
	"criticalsys/gridfs/pkg/config"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Client is a wrapper around the mongo client and bucket.
type Client struct {
	client *mongo.Client
	bucket *gridfs.Bucket
}

// NewClient creates a new GridFS client.
func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	clientOptions := options.Client().ApplyURI(cfg.MongoURI).SetReadPreference(readpref.Secondary()).SetAuth(options.Credential{
		Username: cfg.MongoUser,
		Password: cfg.MongoPass,
	})

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(cfg.MongoDB)
	bucket, err := gridfs.NewBucket(db, options.GridFSBucket().SetName(cfg.MongoGridFSPrefix))
	if err != nil {
		return nil, fmt.Errorf("failed to create GridFS bucket: %w", err)
	}

	return &Client{
		client: client,
		bucket: bucket,
	}, nil
}

// DownloadFile downloads a file from GridFS and saves it to the specified path.
func (c *Client) DownloadFile(fileName, blobPath string, largeFileThreshold int64) error {
	downloadStream, err := c.bucket.OpenDownloadStreamByName(fileName)
	if err != nil {
		return fmt.Errorf("failed to open download stream for file %v: %w", fileName, err)
	}
	defer downloadStream.Close()

	fileSize := downloadStream.GetFile().Length
	filePath := filepath.Join(blobPath, fileName)

	if fileSize > largeFileThreshold {
		// Stream large files
		file, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("failed to create file %v for streaming: %w", filePath, err)
		}
		defer file.Close()

		if _, err := io.Copy(file, downloadStream); err != nil {
			return fmt.Errorf("failed to stream file %v to disk: %w", fileName, err)
		}
	} else {
		// Read small files into memory
		data, err := io.ReadAll(downloadStream)
		if err != nil {
			return fmt.Errorf("failed to read data from download stream: %w", err)
		}

		if err := os.WriteFile(filePath, data, 0644); err != nil {
			return fmt.Errorf("failed to write file %v to disk: %w", filePath, err)
		}
	}

	return nil
}

// Disconnect disconnects the mongo client.
func (c *Client) Disconnect(ctx context.Context) error {
	if err := c.client.Disconnect(ctx); err != nil {
		return fmt.Errorf("failed to disconnect MongoDB: %w", err)
	}
	return nil
}
