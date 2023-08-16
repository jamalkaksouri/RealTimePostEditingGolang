# RealTimePostEditingGolang

A real-time post editing application using SSE (Server-sent events) in Golang.

## Usage

Follow these steps to get started with the RealTimePostEditingGolang application:

1. **Create a New PostgreSQL Database:**
   Set up a new PostgreSQL database where the application data will be stored.

2. **Enable UUID Extension:**
   Run the following SQL query in your PostgreSQL database to enable the UUID extension:
   ```sql
   CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
   CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name TEXT,
    stock_quantity INTEGER,
    version UUID DEFAULT uuid_generate_v4() NOT NULL
   );
   ```

### Configure the `config.yaml` File

Update the `config.yaml` file on your local machine with the necessary configuration settings for the RealTimePostEditingGolang application. This configuration may include database connection details, API endpoints, and other relevant settings.

### Run the Application

To start the RealTimePostEditingGolang application, open your terminal and execute the following command:

```sh
go run main.go

```

This will launch the application and set it up to handle real-time post editing using SSE (Server-sent events).

### Access the Application

Once the application is up and running, you can interact with its real-time post editing features. To do so, follow these steps:

Open a web browser.
Navigate to the `index.html`` file included in the project.
Use a live server to serve the index.html file. This will allow you to experience the real-time post editing capabilities of the application.
