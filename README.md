# Database Converter

This project was create for study purposes. Your main goal is being able to convert
any relational databases into no-relational databases and vice-versa, using 
GO as language.

## Used libraries
* [MySQL](https://github.com/go-sql-driver/mysql) driver
* [MongoDB](https://github.com/mongodb/mongo-go-driver) driver
* [Viper](https://github.com/spf13/viper) for configuration management

## Supported drivers
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

    `./database-converter sample.yml`