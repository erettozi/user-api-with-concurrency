# RESTful API with Concurrency and Data Processing

This project is a RESTful API built in Golang that performs CRUD operations on user data, fetches additional user information concurrently from an external API, processes the data, and writes the results to a CSV file. It also includes a CLI tool to interact with the API.

---

## API Endpoints

- **`POST /users`**: Create a new user.
- **`GET /users`**: Get a list of all users.
- **`GET /users/{id}`**: Get a user by ID.
- **`PUT /users/{id}`**: Update a user by ID.
- **`DELETE /users/{id}`**: Delete a user by ID.

---

## CLI Commands

The CLI tool allows you to interact with the API. To use it, navigate to the `cmd/cli` directory first.

### Building the CLI:
1. Navigate to the `cmd/cli` directory:
   ```bash
   cd cmd/cli
   ```
2. Build the CLI tool:
   ```bash
   go build -o cli
   ```

### Running the CLI:
After building the CLI, you can run it using:
```bash
./cli
```

### CLI Help:
Running the CLI without arguments will display the help message:
```bash
./cli
Usage: cli <command> [options]
Commands:
  fetch-additional-info  Fetch additional information for a user

Use './cli <command> --help' for more information on a specific command.
```

For example, to fetch additional user information:
```bash
./cli fetch-additional-info -id 1
```

---

## Running the Project

### Starting the Server
1. Run the server:
   ```bash
   go run main.go
   ```
   By default, the server will start on port `3000`. You can change the port by setting the `PORT` environment variable.

### Running Tests
To run the tests without overwriting the `users.csv` file, use the `ENV=test` flag:
```bash
ENV=test go test ./...
```
The `ENV=test` flag ensures that a unique filename is generated for the CSV file during tests, preventing overwrites.

---

## Profiling with `pprof`

The project includes a `pprof` server running at `http://localhost:6060/debug/pprof`. This server provides profiling data that can be used to analyze the performance of the application, including CPU usage, memory allocation, and goroutine blocking.

### Accessing `pprof` Data:
1. Start the server as usual.
2. Open your browser and navigate to `http://localhost:6060/debug/pprof`.
3. From there, you can access various profiling endpoints, such as:
   - `/debug/pprof/heap`: Memory allocation profiling.
   - `/debug/pprof/profile`: CPU profiling.
   - `/debug/pprof/goroutine`: Goroutine stack traces.

### Using `pprof` Tool:
You can also use the `go tool pprof` command to analyze the profiling data. For example:
```bash
go tool pprof http://localhost:6060/debug/pprof/profile
```

This will start an interactive session where you can analyze CPU usage.

---

## Environment Variables

The following environment variables can be configured:

- **`PORT`**: The port on which the API server will run. Default: `3000`.
  ```bash
  export PORT=3000
  ```

- **`EXTERNAL_API_URL`**: The URL of the external API used to fetch additional user information. Default: `http://localhost:3000`.
  ```bash
  export EXTERNAL_API_URL=http://localhost:3000
  ```

If these variables are not set, the default values will be used.

---

## Request Payload Examples

### **Create a User (`POST /users`)**
```json
{
    "name": "Erick Rettozi",
    "age": 48,
    "email": "erettozi@tolkien.com"
}
```

### **Update a User (`PUT /users/{id}`)**
```json
{
    "name": "Aragorn Elessar",
    "age": 37,
    "email": "aragorn@tolkien.com"
}
```

---

## Project Structure

```
/project
  /api          # API handlers and routes
  /models       # Data models (e.g., User)
  /services     # Business logic (e.g., fetching external data)
  /utils        # Utility functions (e.g., CSV processing)
  /cmd/cli      # CLI tool to interact with the API
  main.go       # Entry point for the API server
  README.md     # Project documentation
```

---

## How It Works

1. **API Endpoints**:
   - The API supports CRUD operations for user data, stored in memory.
   - Additional user information is fetched concurrently from an external API using Goroutines and channels.

2. **Data Processing**:
   - Users under 18 years old are filtered out.
   - The names of the remaining users are capitalized.
   - The processed data is written to a CSV file.

3. **CLI Tool**:
   - The CLI tool provides commands to interact with the API, such as fetching additional user information.

---

## Contributing

Feel free to open issues or submit pull requests for improvements or bug fixes.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.