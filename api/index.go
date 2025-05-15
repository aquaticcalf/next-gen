package api

import (
	"net/http"
	"time"

	web "github.com/aquaticcalf/go-web"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// Users is a slice of User
type Users []User

// MockUsers provides some sample users
var MockUsers = Users{
	{
		ID:        "1",
		Name:      "John Doe",
		Email:     "john@example.com",
		CreatedAt: time.Now().Add(-24 * time.Hour),
	},
	{
		ID:        "2",
		Name:      "Jane Smith",
		Email:     "jane@example.com",
		CreatedAt: time.Now().Add(-48 * time.Hour),
	},
}

// UserRequest is the expected structure for creating a user
type UserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Create a single app instance to be reused across requests
var app *web.App

// init initializes the app when the package is loaded
func init() {
	// Initialize our web app
	app = web.New()

	// Add middleware
	app.Use(web.LoggerMiddleware())
	app.Use(web.RecoverMiddleware())

	// Setup API routes
	app.Add("/api", func(g *web.Group) {
		// Users endpoints
		g.GET("/users", getUsers)
		g.POST("/users", createUser)
		g.GET("/users/{id}", getUserByID)

		// Health check
		g.GET("/health", healthCheck)
	})

	// Root endpoint
	app.Add("/api", func(g *web.Group) {
		g.GET("", index)
	})
}

// Handler is the main entry point for Vercel serverless functions
func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}

// index handles the root endpoint, need to replace with swagger or similar
func index(c *web.Context) {
	c.HTML(http.StatusOK, `
		<html>
			<head>
				<title>API Server</title>
				<style>
					body {
						font-family: system-ui, -apple-system, sans-serif;
						max-width: 800px;
						margin: 0 auto;
						padding: 2rem;
						line-height: 1.6;
					}
					h1 { color: #333; }
					a { color: inherit; text-decoration: none; }
					.endpoint {
						background: #f7f7f7;
						padding: 1rem;
						border-radius: 4px;
						margin-bottom: 1rem;
						cursor: pointer;
					}
					.method {
						display: inline-block;
						padding: 0.25rem 0.5rem;
						border-radius: 4px;
						font-weight: bold;
						margin-right: 0.5rem;
					}
					.get { background: #e3f2fd; color: #0277bd; }
					.post { background: #e8f5e9; color: #2e7d32; }
					.response {
						display: none;
						background: #f0f0f0;
						border: 1px solid #ddd;
						border-radius: 4px;
						padding: 1rem;
						margin-top: 0.5rem;
						white-space: pre-wrap;
						font-family: monospace;
					}
					.modal {
						display: none;
						position: fixed;
						z-index: 1;
						left: 0;
						top: 0;
						width: 100%;
						height: 100%;
						background-color: rgba(0,0,0,0.4);
					}
					.modal-content {
						background-color: white;
						margin: 15% auto;
						padding: 20px;
						border-radius: 5px;
						width: 70%;
					}
					.close {
						color: #aaa;
						float: right;
						font-size: 28px;
						font-weight: bold;
						cursor: pointer;
					}
					.close:hover { color: black; }
					textarea {
						width: 100%;
						height: 150px;
						margin-bottom: 10px;
						padding: 8px;
					}
					button {
						background: #0070f3;
						color: white;
						border: none;
						padding: 8px 16px;
						border-radius: 4px;
						cursor: pointer;
					}
					.repo-link { color: #0070f3; text-decoration: none; }
					.repo-link:hover { text-decoration: underline; }
				</style>
			</head>
			<body>
				<h1>API Documentation</h1>
				<p>Welcome to the next gen API server. Below are the available endpoints:</p>
				
				<div class="endpoint" onclick="makeRequest('GET', '/api/users')">
					<span class="method get">GET</span>
					<strong>/api/users</strong> - Get all users
					<div id="users-response" class="response"></div>
				</div>
				
				<div class="endpoint" onclick="showModal()">
					<span class="method post">POST</span>
					<strong>/api/users</strong> - Create a new user
					<div id="create-user-response" class="response"></div>
				</div>
				
				<div class="endpoint" onclick="getUserById()">
					<span class="method get">GET</span>
					<strong>/api/users/{id}</strong> - Get a specific user by ID
					<div id="user-by-id-response" class="response"></div>
				</div>
				
				<div class="endpoint" onclick="makeRequest('GET', '/api/health')">
					<span class="method get">GET</span>
					<strong>/api/health</strong> - Health check endpoint
					<div id="health-response" class="response"></div>
				</div>
				
				<div id="post-modal" class="modal">
					<div class="modal-content">
						<span class="close" onclick="closeModal()">&times;</span>
						<h3>Create New User</h3>
						<textarea id="post-data" placeholder="Enter JSON data for the new user">{"name": "New User", "email": "user@example.com"}</textarea>
						<button onclick="createUser()">Create User</button>
					</div>
				</div>

				<p>For more information, check out the <a href="https://github.com/aquaticcalf/next-gen/" class="repo-link">GitHub repository</a>.</p>

				<script>
					function makeRequest(method, url) {
						const responseId = url.replace(/\//g, '-').replace(/^-/, '') + '-response';
						const responseElem = document.getElementById(responseId);
						responseElem.style.display = 'block';
						responseElem.textContent = 'Loading...';
						
						fetch(url, {method})
							.then(response => response.json())
							.then(data => {
								responseElem.textContent = JSON.stringify(data, null, 2);
							})
							.catch(error => {
								responseElem.textContent = 'Error: ' + error.message;
							});
					}

					function getUserById() {
						const id = prompt('Enter user ID:');
						if (id) {
							makeRequest('GET', '/api/users/' + id);
						}
					}

					function showModal() {
						document.getElementById('post-modal').style.display = 'block';
					}

					function closeModal() {
						document.getElementById('post-modal').style.display = 'none';
					}

					function createUser() {
						const jsonData = document.getElementById('post-data').value;
						const responseElem = document.getElementById('create-user-response');
						responseElem.style.display = 'block';
						responseElem.textContent = 'Submitting...';
						
						fetch('/api/users', {
							method: 'POST',
							headers: {
								'Content-Type': 'application/json'
							},
							body: jsonData
						})
						.then(response => response.json())
						.then(data => {
							responseElem.textContent = JSON.stringify(data, null, 2);
							closeModal();
						})
						.catch(error => {
							responseElem.textContent = 'Error: ' + error.message;
							closeModal();
						});
					}

					// Close modal when clicking outside of it
					window.onclick = function(event) {
						const modal = document.getElementById('post-modal');
						if (event.target == modal) {
							modal.style.display = 'none';
						}
					}
				</script>
			</body>
		</html>
	`)
}

// getUsers returns all users
func getUsers(c *web.Context) {
	// You could add pagination, filtering, etc. here
	c.Success(MockUsers)
}

// getUserByID returns a specific user by ID
func getUserByID(c *web.Context) {
	id := c.GetParam("id")

	// Find the user with the specified ID
	for _, user := range MockUsers {
		if user.ID == id {
			c.Success(user)
			return
		}
	}

	c.Error(http.StatusNotFound, "User with ID %s not found", id)
}

// createUser creates a new user
func createUser(c *web.Context) {
	var req UserRequest

	if err := c.Bind(&req); err != nil {
		c.Error(http.StatusBadRequest, "Invalid request body: %v", err)
		return
	}

	// Validate request
	if req.Name == "" {
		c.Error(http.StatusBadRequest, "Name is required")
		return
	}

	if req.Email == "" {
		c.Error(http.StatusBadRequest, "Email is required")
		return
	}

	// Create a new user
	newUser := User{
		ID:        generateID(),
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
	}

	// Add to our mock database
	MockUsers = append(MockUsers, newUser)

	// Return the created user
	c.Success(newUser)
}

// healthCheck returns server health status
func healthCheck(c *web.Context) {
	c.Success(map[string]interface{}{
		"status": "ok",
		"time":   time.Now(),
	})
}

// generateID creates a new ID (simplified version)
func generateID() string {
	// In a real app, use a proper ID generation method
	return time.Now().Format("20060102150405")
}
