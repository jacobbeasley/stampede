# Buffalo Testing Setup & Execution

When writing or executing tests in this project, AI agents must ensure the environment is correctly prepared.

## 1. Database Preparation
Buffalo test suites (`ActionSuite`, `ModelSuite`) are tightly coupled to the database. They require an active PostgreSQL connection to the `_test` database defined in `database.yml`.
If the database doesn't exist, the test binary will fail to compile or run.

**Before running any tests, you MUST run:**
```bash
buffalo pop create -a
buffalo pop migrate
```
*Note: If `buffalo` CLI is broken or missing, use `soda create -a` and `soda migrate`.*

## 2. Running Tests
You can run all tests, or scope them to a package:
```bash
# Run everything
buffalo test

# Run only model tests
buffalo test ./models/...

# Run only action tests
buffalo test ./actions/...
```

## 3. Testing Isolated Middleware
If you need to test standard Go middleware without requiring a database connection, you can instantiate a minimal, isolated Buffalo app:
```go
func Test_MyMiddleware(t *testing.T) {
    app := buffalo.New(buffalo.Options{})
    app.Use(MyMiddleware)
    app.GET("/", func(c buffalo.Context) error {
        return c.Render(http.StatusOK, r.String("OK"))
    })

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/", nil)
    app.ServeHTTP(w, req)

    assert.Equal(t, http.StatusOK, w.Code)
}
```
