# knot-virtualthing-SQL
 
KNoT SQL aims to provide an abstraction that allows DBMS to interact with KNoT Cloud services, virtualizing a KNoT device. This service virtualizes DBMS as KNoT device. Currently, the following DBMS are supported:
- CosmosDB;
- SQL Server;
- Oracle.
 
# Basic installation and usage
 
### System Installation and Usage
 
- This project requires Go 1.8 or higher;
- gocosmos library 0.1.8 or higher (https://github.com/AntonioJanael/gocosmos). The gocosmos presented some bugs that made it impossible to capture data from CosmosDB. Therefore, it was necessary to create a fork and implement the bug fixes in it;
- go-mssqldb library 0.12.2 or higher;
- go-ora* library 2.4.25 or higher.
 
 *This library is especially useful because it does not require the installation of the "Oracle Instant Client".

### Testing

To run the complete suite of tests, you will need to provide instances of CosmosDB, SQL Server and Oracle, as well as configure the tests with access credentials and queries. 

go test -v -coverprofile=coverage.out .\\...

### Configuration
- The internal/configuration/application_configuration.yaml file contains the general specifications of the application. Parameters:
    - intervalBetweenRequestInSeconds: interval in seconds between successive requests to the database. This interval is unique for all sensors;
    - datesPersistenceFilepath: internal/configuration/latest_record_timestamp.yaml;
    - dataRecoveryPeriodInHours: hours prior to the execution of the application in which historical data will be collected. By default, this value is 0, which means that the application will start collecting from the moment of its execution;
    - pertinentTags: mapping between the sensor id and a name that identifies it in the database. Example:
        - 1: "XPTO1"
        - 2: "XPTO2"
    - logFilepath: file path where log messages are stored;
    - numberParallelTags: number of concurrent requests made to the database. This value depends on the capacity of the DBMS instance used by the customer. This parameter requires special attention to avoid overloading the client's DBMS;
    - context: context in which KNoT SQL will operate. Options: cosmosdb, sqlserver and oracle.  
- For cosmosDB with 2000 RUs, the number of parallel tags, based on experiments, the recommended value is 40. The number of parallel tags depends on the throughput configured in the DBMS. Therefore, experiments are needed to identify the appropriate number of tags in parallel. These experiments consist of evaluating different amounts of tags and verifying which one returns acceptable results. For cosmosDB, we selected the largest amount of tags where the percentage of 429 messages returned by cosmosDB was less than 5%, according to a Microsoft recommendation;
- Internal/configuration/database_configuration.yaml: contains the information needed to connect to the DBMSs. Parameters:
    - driver (required): driver name used to connect to the DBMS. The names of currently used drivers are available at the end of this document;
    - connectionString (required): database connection string or Data Source Name (DSN);
    - IP (optional);
    - port (optional);
    - username (optional);
    - password (optional);
    - database (optional).
- Internal/configuration/device_config.yaml: KNoT device configuration;
- Internal/configuration/knot_setup.yaml: KNoT cloud configuration;
- Internal/configuration/latest_record_timestamp.yaml: mapping between sensor identifier and timestamp. Example:
    - 1: 2022-09-30 10:00:00
    - 2: 2022-09-28 07:44:01 
- Internal/scripts/queries: contains a mapping between the sensor id and the SQL query used to get the value and timesamp of this sensor. Example:
    - 1: "SELECT value, timestamp FROM sensor_data WHERE tagname='XPTO1' AND timestamp > '2022-09-30 10:00:00' ORDER BY timestamp ASC"
    - 2: "SELECT value, timestamp FROM sensor_data WHERE tagname='XPTO2' AND timestamp > '2022-09-28 07:44:01' ORDER BY timestamp ASC"
- <span style="color:red">**Caution: the sensor id needs to be consistent across all files.**</span>

## Environment variable
```sh
export DEVICE_CONFIG_FILEPATH=internal/configuration/device_config.yaml
```
 
### Special tip
When working with gocosmos, to avoid problems with creating unnecessary goroutines, comment out line 793 (go db.connectionOpener(ctx)) of file Go\src\database\sql\sql.go. These goroutines are created and wait for a termination signal sent by the database driver. The driver used, gocosmos, does not send this signal as this operation is irrelevant to it. Thus, created goroutines are never terminated and would lead KNoT SQL to a goroutine leak. Other drivers may dispense with this configuration.
 
# Project Structure and Contribution
This application is designed to continually evolve and is extensible to incorporate new DBMSs and clients. Therefore, integration with DBMSs occurs through repositories and the components that capture data from specific clients (collectors) were developed with the strategy design pattern. This arrangement allows changes to collectors and repositories without changing the application core, while facilitating application extensibility.
 
To include new DBMSs, implement a suitable repository. Similarly, to embed a collector for another client, follow the interface used with the strategy pattern. The business rules of each client must be in the collectors.


## Guide to add support for other DBMSs
If you want to add support for new DBMSs, we recommend following the guide below in order to leverage existing code.
- Import the proper library in connection.go file;
- Create a structure composed of "SQL" in database_repositories.go file;
- Implement a collector inside the "application" folder. Examples are in the cosmos.go and sql_server.go files;
- In the "collector_utils.go" file, add a constant with the "context" (examples on lines 12-14);
- In the "collector_utils.go" file, in the "newSQLPlaceholderHandlerMapping" function, update the mapping that indicates how to manage the sql placeholder ':' in the specific DBMS;
- In the "builder.go" file, in the "NewBuilderMapping" function, update the mapping that indicates how to build the structure for reading data from the new DBMS;
- Still in the "builder.go" file, implement the "builder" for the new DBMS. Examples between lines 130 and 159;
- Adjust the configuration according to the specifics of the new DBMS.

It is worth mentioning that the step-by-step guide above is a general guide and that adaptations may be necessary to adapt to a specific DBMS.

# Docker installation and usage
## Building and running
Change the working directory to the project root:
```bash
$ cd <path/to/knot-virtualthing-sql>
```

With Docker installed, to avoid problems caused by conflicts with the host's IP range, create a new Docker network. In the homologation environment, avoid networks with the following prefixes: 172.24, 172.25, and 172.26. This step needs to be done regardless of creating the image or downloading from the remote repository.
```bash
docker network create -d bridge --subnet=172.50.0.0/24 sql
```

To create the docker image:
```bash
$ docker build . --file docker/Dockerfile --tag knot-sql
```

To run a container with the created image:
```bash
$ docker run --restart unless-stopped --name knot-sql -v $(pwd)/internal/configuration/:/app/internal/configuration/ -v $(pwd)/internal/scripts/:/app/internal/scripts/ --net=sql -it knot-sql
```

If you want to use an already created docker image available on ECR AWS, use the command below:

```bash
docker pull 239644445781.dkr.ecr.us-east-1.amazonaws.com/knot-sql:v01.00-rc01
```

Some files need to be persisted outside the container to avoid data loss. Create knot-sql-configuration folder:
```bash
$ mkdir knot-sql-configuration && sudo chmod 777 knot-sql-configuration && cd knot-sql-configuration
```

Create temporary docker image to extract configuration files:
```bash
$ docker run -d --name knot-sql -v $(pwd)/internal/configuration/:/app/internal/configuration/ -v $(pwd)/internal/scripts/:/app/internal/scripts/ --net=sql -it 239644445781.dkr.ecr.us-east-1.amazonaws.com/knot-sql
```

Get configuration files from the container:
```bash
$ docker cp knot-sql:/app/internal/configuration .
```

Stop temporary Docker container:
```bash
$ docker stop knot-sql
```

Remove temporary container:
```bash
$ docker rm knot-sql
```

Run container with data persistence
```bash
$ docker run -d --restart unless-stopped --name knot-sql -v $(pwd)/internal/configuration/:/app/internal/configuration/ -v $(pwd)/internal/scripts/:/app/internal/scripts/ --net=sql -it knot-sql
```

To look at knot-sql logs in container:
```bash
$ docker logs -f <docker-container>
```

# DBMS configuration
For all DBMSs, the order of queried columns must be: value and timestamp. To switch between DBMSs, you just need to:
- Specify the value of the "context" parameter in application_configuration.yaml. Current options: cosmosdb, sqlserver, and oracle;
- Inform the proper driver, connection string and SQL query as described below. Also, you need to specify "the latest record timestamp" in the file "latest_record_timestamp.yaml". This file contains a mapping between the sensor id and the timestamp, in the following format: id: timestamp;
- <span style="color:red">**Caution: for CosmosDB, the ":" in timestamp must be enclosed by '\`', such as "00\`:\`00\`:\`00". On the other hand, for SQL Server and Oracle, the '\`' must be avoided: "00:00:00"**</span>;
- <span style="color:red">**Caution: queries should be sorted based on timestamp in ascending order**</span>.
- <span style="color:red">**Caution: the query needs to filter records by timestamp, in the following example format: WHERE timestamp > '2022-09-28 07:44:01'. This filter ensures that KNoT SQL always processes the most recent data. In the query, the timestamp must be specified as a placeholder, %s, where KNoT SQL will replace the placeholder with the correct value of the timestamp.**</span>.
- <span style="color:red">**Caution: the following datetime format is expected: "YYYY-MM-DD HH:mm:ss". SQL functions can be used to guarantee this datetime format**</span>.




## CosmosDB
Driver name:
```bash
gocosmos
```

Connection string: 
```bash
AccountEndpoint=<endpoint>;AccountKey=<key>;DefaultDb=<database>
```
For some reason, the data returned from CosmosDB is not in the expected order: value and timestamp. These records are in reverse order: timestamp and value. To work around this issue, just specify an alias for the value column. Also, avoid using "value" as an alias.
Query format:
```bash
SELECT <table>.<value_column_name> AS <alias>, <table>.<timestamp_column_name> FROM <container_name> <table> WHERE <table>.<tag_column_name>='%s' AND <table>.<timestamp_column_name> > '%s' ORDER BY <table>.<timestamp_column_name> ASC
```

## SQL Server
Driver name:
```bash
mssql
```

Connection string: 
```bash
server=<server>;user id=<username>;password=<password>;port=<port>;database=<database>
```

Query format:
```bash
SELECT <value_column_name>, <timestamp_column_name> FROM <table/view> WHERE <timestamp_column_name> > '%s' ORDER BY <timestamp_column_name> ASC
```

## Oracle
Driver name:
```bash
oracle
```

Connection string: 
```bash
(DESCRIPTION=
    (ADDRESS_LIST=
        (ADDRESS=(protocol=TCP)(host=<server_ip>)(port=<port>))
    )
    (CONNECT_DATA=
        (SERVICE_NAME=<service_name>)
    )
)
```

Query format:
```bash
SELECT <value_column_name>, <timestamp_column_name> FROM <table/view> WHERE <timestamp_column_name> > '%s' ORDER BY <timestamp_column_name> ASC
```
# KNoT Cloud configuration
For this you need the KNoT cloud Token generated by KNoT, KNoT Cloud AMQP IP, AMQP user and AMQP password. So, go to internal/configuration/knot_setup.yaml or create a file in the same path and use the folow tamplete:
```yaml
url: "amqp://knoAMQPtUser:knotAMQPPassword@IP_Knot_AMQP:5672"
user_token : "KNoT_token"
```

# Sensor configuration
On KNoT SQL you can use ich query as a sensor doing queries with JOIN and sub Selects, fo that, first you have to device on internal/configuration/application_configuration.yaml and set the pertinentTags as demostrated on the folow example:
```yaml
pertinentTags:
	1:sensor_1_name
	2:sensor_2_name
```
Here were define the sensor with id 1 as sensor_1_name and id 2 as sensor_2_name.

after that go to internal/scripts/queries.yaml and use the tamplete showing to define ich query for ich sensor id.

```yaml
mapping:
	1:"SELECT value, timestamp FROM table_1"
	2:"SELECT value, timestamp FROM table_2"
```
Here we define the where the data of sensor_1_name and sensor_2_name with id 1 and 2 comes from, it is table_1 and table_2. Folow this patter, alway use value first and timestamp as second. Your DataBase must have this information or use the system time to create the timestamp.

Now you know the sensor data type, lets config the KNoT device.The KNoT sql as reconized as one KNoT device, all your data will be spread as KNoT sensors, for this you must config ich sensor with its ID, name unit type and value type, go to internal/configuration/device_config.yaml, there is a map of devices and ich device has multiplos sensors. If there is no file devices_config.yaml, create one with this tamplate and change the name, a new id and token will be generated when the registration with KNoT cloud was done:

The possible combinations of these parameters are specified in the link below:
https://knot-devel.cesar.org.br/doc/thing/unit-type-value.html

```yaml
 #device_config.yaml
0d7cd9d221385e1f:
  id: 0d7cd9d221385e1f
  token: ""
  name: "name"
  config:
  - sensorId: 1
    schema:
      valueType: 2
      unit: 1
      typeId: 65296 
      name: name
    event:
      change: true
      timeSec: 0
      lowerThreshold: null 
      upperThreshold: null
  state: new
  data: []
  error: ""
```