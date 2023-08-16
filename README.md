# RealTimePostEditingGolang

**SSE (Server-sent events) in Golang**

## Usage:

1. Create new database in postgres
2. Run query `CREATE EXTENSION "uuid-ossp";`
3. Create new table `products`:
   `CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name TEXT,
    stock_quantity INTEGER,
    version UUID DEFAULT uuid_generate_v4() NOT NULL
);`
4. Run `go run main.go`
5. Open `index.html` in browser (using live server) 
