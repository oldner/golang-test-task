# Project Name

Golang Test Task

## Overview

This project implements a simple messaging system using a microservices architecture. It consists of three main components:

1. API Service
2. Message Processor
3. Reporting API

The system allows users to send messages, process them asynchronously, and retrieve message history.

## Components

### 1. API Service

- Exposes an HTTP endpoint for sending messages
- Publishes messages to a RabbitMQ queue

### 2. Message Processor

- Consumes messages from the RabbitMQ queue
- Processes and stores messages in Redis

### 3. Reporting API

- Provides an HTTP endpoint for retrieving message history
- Fetches processed messages from Redis

## Technologies Used

- Go (Golang)
- Docker and Docker Compose
- RabbitMQ
- Redis
- Gin (Web Framework)

## Prerequisites

- Docker
- Docker Compose

## Getting Started

1. Clone the repository:
   ```
   git clone [repository-url]
   cd [project-directory]
   ```

2. Create a `.env` file in the project root and add the necessary environment variables:
   ```
   API_PORT=8080
   REPORTING_API_PORT=8081
   RABBITMQ_URL=amqp://rabbitmq_user:rabbitmq_password@rabbitmq:5672/
   REDIS_URL=redis:6379
   RABBITMQ_DEFAULT_USER=rabbitmq_user
   RABBITMQ_DEFAULT_PASS=rabbitmq_password
   ```

3. Build and start the services:
   ```
   docker-compose up --build
   ```

4. The services will be available at:
   - API Service: http://localhost:8080
   - Reporting API: http://localhost:8081

## Usage

### Sending a Message

Send a POST request to `http://localhost:8080/v1/message` with a JSON body:

```json
{
  "sender": "user1",
  "receiver": "user2",
  "message": "Hello, World!"
}
```

### Retrieving Messages

Send a GET request to `http://localhost:8081/v1/message/list?sender=user1&receiver=user2`