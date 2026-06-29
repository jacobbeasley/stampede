# Application Architecture and Advanced Features

This document provides architectural guidance and best practices for extending this starter repository with advanced features and integrations.

## 1. Observability, Logging, and Metrics

When deploying to production, observability is crucial for monitoring application health and debugging issues.

### Structured Logging
*   **Recommendation:** Configure Buffalo's logger to output JSON in production. JSON logs are easily parseable by log aggregation tools (e.g., Datadog, ELK stack, Google Cloud Logging).
*   **Implementation:** Buffalo uses `logrus` under the hood. You can configure a JSON formatter in `main.go` or `actions/app.go` based on the environment:
    ```go
    if envy.Get("GO_ENV", "development") == "production" {
        app.Logger = buffalo.NewLogger("production")
        // Note: Buffalo's default production logger outputs JSON.
    }
    ```

### Request Correlation (Tracing)
*   **Recommendation:** Implement a middleware that injects a unique `X-Request-ID` or `Trace-ID` into every incoming HTTP request context.
*   **Implementation:** Buffalo provides the `buffalo.RequestID` middleware which is often included by default or can be easily added to your `app.Use()` stack. Pass this ID along to background workers (e.g., mailers) to trace actions across the system.

### Metrics (Prometheus)
*   **Recommendation:** Expose a `/metrics` endpoint using the Prometheus Go client to track request counts, latency, and memory usage.
*   **Implementation:** Use a community middleware or build a simple custom middleware that records request durations and statuses, and serves the Prometheus registry at a dedicated route.

## 2. API Documentation (Swagger/OpenAPI)

If you intend to expose a public API or build a separate frontend (e.g., mobile app), generating API documentation is highly recommended.

*   **Tooling:** Use `swaggo/swag` (https://github.com/swaggo/swag) to generate Swagger 2.0 / OpenAPI documentation directly from Go code comments.
*   **Workflow:**
    1.  Add declarative comments above your Buffalo action handler functions:
        ```go
        // @Summary Create a task
        // @Description Adds a new task for the current user
        // @Accept json
        // @Produce json
        // @Success 201 {object} models.Todo
        // @Router /api/todos [post]
        func TodosCreate(c buffalo.Context) error { ... }
        ```
    2.  Run `swag init` to generate a `docs/swagger.json` file.
    3.  Mount a Swagger UI route (e.g., `/api/docs`) using a package like `github.com/swaggo/http-swagger` to serve the interactive documentation.

## 3. Handling File Uploads

When your application requires users to upload files (avatars, attachments, etc.), consider the following best practices:

*   **Binding:** Use Buffalo's built-in file binding:
    ```go
    file, err := c.File("avatar")
    ```
*   **Storage Strategy:**
    *   **Avoid Local Disk:** Do not store files on the local disk of the container (e.g., Cloud Run or Kubernetes), as it is ephemeral and changes are lost on restart.
    *   **Cloud Storage:** Stream uploaded files directly to object storage services like Amazon S3, Google Cloud Storage, or Azure Blob Storage.
*   **Implementation:** Use the official Go SDKs (e.g., `aws-sdk-go-v2` or Google Cloud Storage Go client) to handle the upload. Consider processing large files or image resizing in a background worker to avoid blocking the HTTP response.

## 4. Real-time Features (WebSockets)

If you need real-time functionality (e.g., live notifications, chat), you can implement WebSockets.

*   **Tooling:** Use the `gorilla/websocket` package, as Buffalo is built on top of Gorilla Mux.
*   **State Management in Multi-Instance Environments:**
    *   If you deploy to a serverless or load-balanced environment (like Cloud Run), WebSocket connections will be distributed across multiple instances.
    *   **Pub/Sub Required:** To broadcast a message to a user, you cannot rely on in-memory connection maps. You must introduce a Pub/Sub mechanism (e.g., Redis Pub/Sub, NATS, or Google Cloud Pub/Sub) to publish events that all application instances subscribe to, allowing the instance holding the specific WebSocket connection to relay the message.

## 5. Health Checks

The application includes built-in endpoints for infrastructure monitoring (e.g., Kubernetes probes, Load Balancer health checks):
*   `/api/ready` - Returns immediately. Used to verify the application process has started and can accept connections (Liveness probe).
*   `/api/health` - Verifies connectivity to critical dependencies (e.g., the PostgreSQL database). Used to verify the application is fully functional (Readiness probe).
