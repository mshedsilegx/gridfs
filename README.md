gridfs high performance extraction utility, for legacy mongodb, to current directory of a local file system. Since compatibility with mongodb 3.x is required, mongo-driver is frozen to release 1.10.6. A configuration  file (ex: ```gridfs_extract_<instance>.properties```) is needed of the following format:

```
# Application - On Prem DEV Environment
MONGO_URI=<proto>://<hostname>:<port>
MONGO_USER=<service_account_readonly>
MONGO_PASS=<password>
MONGO_DB=<database>
MONGO_GRIDFS_PREFIX=<grid fs collection, default is: default.fs>
```

```proto```: mongodb for onprem, mongodb+srv for Atlas\
```port```: 27017 by default, cannot be specified for Atlas

Call syntax:\
```gridfs -config <config_file> -bloblist <error_file_list> -blobpath <Stored_blob_path>```

Example of execution:\
```gridfs -config /etc/gridfs_extract_dev.properties -bloblist /var/lib/gridfs/meta/dev_20250307-105801.txt -blobpath /var/lib/gridfs/data/20240101-20241231```
