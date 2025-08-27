# GridFS Extractor

## Overview

This utility is a high-performance tool for extracting files from a MongoDB GridFS collection to a local file system. It is designed to be fast and efficient, especially when dealing with a large number of files or very large files.

The application is built in Go and leverages concurrency to download multiple files in parallel. It also uses a memory-efficient streaming approach for large files to minimize its resource footprint.

## Features

- **Concurrent Downloads:** Downloads multiple files from GridFS concurrently to maximize download speed.
- **Memory Efficient:** Uses a streaming approach for large files to minimize memory usage. For smaller files, it uses an in-memory approach for maximum speed.
- **Configurable:** Performance-related parameters like the number of concurrent workers and the threshold for large files can be easily configured.
- **Command-line Interface:** Simple and easy-to-use command-line interface.

## Usage

The application is run from the command line with the following syntax:

```bash
gridfs -config <path_to_config_file> -bloblist <path_to_blob_list_file> -blobpath <path_to_store_blobs>
```

### Command-line Flags

- `-config`: (Required) Path to the configuration properties file.
- `-bloblist`: (Required) Path to a text file containing a list of blob filenames to be retrieved (one filename per line).
- `-blobpath`: (Required) Path to the directory where the retrieved blob files will be stored.
- `-version`: (Optional) Display the application version.

## Configuration

The application is configured using a `.properties` file. You can specify the path to this file using the `-config` flag.

### Configuration File Structure

The configuration file should be a simple key-value store. Here are the available options:

| Key                      | Description                                                                                                | Default Value |
| ------------------------ | ---------------------------------------------------------------------------------------------------------- | ------------- |
| `MONGO_URI`              | The connection URI for the MongoDB server.                                                                 | (none)        |
| `MONGO_USER`             | The username for authentication.                                                                           | (none)        |
| `MONGO_PASS`             | The password for authentication.                                                                           | (none)        |
| `MONGO_DB`               | The name of the database to connect to.                                                                    | (none)        |
| `MONGO_GRIDFS_PREFIX`    | The prefix for the GridFS collections (e.g., `fs` for `fs.files` and `fs.chunks`).                           | `fs`          |
| `NUM_WORKERS`            | The number of concurrent workers for downloading files.                                                    | `10`          |
| `LARGE_FILE_THRESHOLD_MB` | The file size in megabytes (MB) above which files are considered "large" and will be streamed to disk.     | `20`          |

### Example Configuration File

Here is an example of a `gridfs_extract.properties` file:

```properties
# MongoDB Connection Details
MONGO_URI=mongodb://localhost:27017
MONGO_USER=myuser
MONGO_PASS=mypassword
MONGO_DB=mydatabase
MONGO_GRIDFS_PREFIX=fs

# Performance Tuning
NUM_WORKERS=20
LARGE_FILE_THRESHOLD_MB=50
```

## Example Usage

Here is a full example of how to run the application.

1.  **Create a blob list file** named `my_blobs.txt`:

    ```
    file1.pdf
    image_archive.zip
    large_video.mp4
    ```

2.  **Create a configuration file** named `config.properties` with your MongoDB details.

3.  **Run the application:**

    ```bash
    gridfs -config config.properties -bloblist my_blobs.txt -blobpath /data/extracted_files
    ```

This command will read the list of files from `my_blobs.txt`, connect to the MongoDB server specified in `config.properties`, and download the files to the `/data/extracted_files` directory.
