# Consumer RMQ - Golang

This is the consumer counterpart to `ktor-rmq` consumer. 

Consumes messages from the queue and stores information to a `PostgreSQL` database. 

The easiest way to get `RabbitMQ` locally on machines is to pull the [docker image](https://hub.docker.com/_/rabbitmq).
There are two options:
- Base image: `docker pull rabbitmq:3`
- Base image with `Management plugin`: `docker pull rabbitmq:3-management`

The management plugin provides a web based dashboard to inspect queue information. Accessed at port `15672`.

Since we can expose docker ports and map them to `localhost` ports, in development terms our apps can locally interact with
the `RabbitMQ` framework, and later they can be deployed to their own containers.

As for Postgres, it can also be pulled from `dockerhub`
- `docker pull postgres:14`

If size is a concern, an `alpine` image is also available.

## Run the applicaiton
1. Pull `RabbitMQ` image from dockerhub
2. Create local docker network `docker network create <network_name>`
3. Run `RabbitMQ` image `docker run -d --rm --net <network_name> --hostname <host_name> --name <container_name> <image_name>`
   1. `<host_name>` is **important** as our applications need to know this name to work.
4. If downloaded base `RabbitMQ` image, enable `management` plugin
   ```
    # Run docker image
   docker run -d --rm --net <network_name> --hostname <host_name> --name <container_name> <image_name>
   
    # Access local container terminal
   docker exec -it <container_name> bash
   
    # Enable management plugin
   rabbitmq-plugins enable rabbitmq_management
    ```

### Running PSQL and Consumer

1. Pull `Postgres` image from dockerhub
2. Run `docker run --name <c_name> -e POSTGRES_PASSWORD=<some_password> -d <image_name>`
3. Go inside the container `docker exec -it <c_name> bash`
4. Create database for the consumer application
    ```
    psql -U postgres
    CREATE DATABASE <db_name>; // DON'T forget the semmicolon
    GRANT PRIVILEGES ON DATABASE <db_name> TO postgres;
    ```
5. `psql` CLI allows execution of standard SQL statements to view and interact with tables.
6. Build and run docker image for consumer
    ```
    # Open terminal in app root folder(dockerfile is accessible)
    docker build --tag <name>:<tag>

    # Run the image in container, the app assumes that RabbitMQ and Postgres are using the same <host_name>
    docker run -d -it --rm --net <network_name> --name <c_name> -e RABBIT_HOST=<host_name> -e RABBIT_PORT=5672 -e RABBIT_USER=guest -e RABBIT_PASS=guest -e PSQL_USER=postgres -e PSQL_PASS=<some_password> -e PSQL_PORT=5432 -e PSQL_NAME=<db_name> <name>:<tag>
    ``` 