# Database Converter

This project was create for study purposes. Your main goal is being able to convert
any (no-)relational databases into (no-)relational databases, using GO as language.

If you want to join me in this project, be my guest =D

## Used libraries
* [MySQL](https://github.com/go-sql-driver/mysql) driver
* [MongoDB](https://github.com/mongodb/mongo-go-driver) driver
* [Viper](https://github.com/spf13/viper) for configuration management

## Supported drivers so far
* MySQL
* MongoDB

## Sample

The sample.yml file contains a sample of a config file that is configured to
convert a MySQL database into a MongoDB database, both available as containers.

## Running

1. Deploying database containers: 

    `docker-compose up`

2. Building the application

    `go build`
    
3. Running the application

    `./database-converter -config=sample.yml`
    
4. Running unit tests

    `go test ./test/...`    
    
## Todo

- [ ] Unit tests
- [ ] Log centralization
- [ ] Fix date/time data
- [ ] Enable sequential conversion 
- [ ] Enable pagination during conversion
- [ ] Enable conversion from MongoDB -> MySQL
- [ ] Enable truncation of destination table/collection
- [ ] Enable composition of table/collection
- [ ] Enable ID preservation
- [ ] Pipelines and map functions enabling some transformation
- [ ] ...

## Notice that...

- Sometime, in the future, it could become a REST API, so, this is why you should see some 
MarshalJSON implementations - I mean, it's not used right now, but it's there.

- For simplicity's sake, I decided to use only conceptual names from relational databases.
For example, a MongoDB Collection is represented here as a Table object.

- ...