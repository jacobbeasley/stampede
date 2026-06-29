---
name: manage-background-jobs
description: Guides the creation, registration, and dispatching of background jobs using the Buffalo worker system.
---

# Manage Background Jobs Skill

This skill outlines how to implement background processing in the Buffalo framework, which is critical for long-running tasks like email delivery, data processing, and report generation.

## 1. Context and Planning
Before implementing a background job:
- Ensure the task genuinely requires background execution (e.g., it blocks the HTTP response for > 500ms).
- Identify the data payload required by the job. Keep it minimal (e.g., pass an ID instead of a whole database record).

## 2. Defining the Worker Handler
Workers in Buffalo use the `worker.Handler` signature.
- Create a new file in the `workers/` directory (e.g., `workers/email.go`).
- Define the handler function:
  ```go
  func SendEmailHandler(args worker.Args) error {
      // Extract args
      userID, ok := args["user_id"].(string)
      if !ok {
          return errors.New("user_id is missing or invalid")
      }

      // Perform task (e.g., fetching user from DB, sending email)
      return nil
  }
  ```

## 3. Registering the Worker
You must register the handler with the worker adapter so the system knows how to process the queue.
- Open `workers/workers.go` (or `actions/app.go` if the project structure dictates).
- Register the handler during application initialization:
  ```go
  w.Register("send_email", SendEmailHandler)
  ```

## 4. Dispatching the Job
To enqueue the job for processing from within a Buffalo Action or Pop ORM model hook:
- Use the `worker` instance (often available via the Buffalo context `c` or globally).
- Enqueue the job:
  ```go
  err := app.Worker.Perform(worker.Job{
      Queue:   "default",
      Handler: "send_email",
      Args: worker.Args{
          "user_id": user.ID.String(),
      },
  })
  ```

## 5. Testing Workers
- Write unit tests for your worker handlers directly in `workers/*_test.go` by invoking the handler function with mocked `worker.Args`.
- Use Buffalo's testing tools to ensure jobs are queued correctly when an action is hit.
