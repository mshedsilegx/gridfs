gridfs high performance extraction utility, for legacy mongodb, to current directory of a local file system. A configuration  file (ex: gridfs_extract_<instance>.properties) is needed of the following format:

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
```gridfs <config_file> <error_file_list>```

Example of execution:\
```gridfs /etc/gridfs_extract_dev.properties /var/lib/gridfs/meta/dev_20250307-105801.error```
