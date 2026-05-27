package httpapi

import (
	"crypto/subtle"
	"net/http"
)

const swaggerHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <title>Olario API Docs</title>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.ui = SwaggerUIBundle({
      url: "/swagger/openapi.json",
      dom_id: "#swagger-ui"
    });
  </script>
</body>
</html>`

const openAPISpec = `{
  "openapi": "3.0.3",
  "info": {
    "title": "Olario Platform Backend API",
    "version": "0.1.0",
    "description": "Development API documentation for the Olario modular monolith."
  },
  "servers": [
    {
      "url": "/",
      "description": "Same origin as Swagger UI"
    }
  ],
  "paths": {
    "/healthz": {
      "get": {
        "summary": "Health check",
        "responses": {
          "200": {
            "description": "API process is running"
          }
        }
      }
    },
    "/api/v1/auth/register": {
      "post": {
        "summary": "Register tenant admin by invitation",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/RegisterRequest"
              }
            }
          }
        },
        "responses": {
          "201": {
            "description": "Registered and authenticated"
          },
          "400": {
            "description": "Invalid invitation or request"
          }
        }
      }
    },
    "/api/v1/auth/login": {
      "post": {
        "summary": "Login tenant user",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/LoginRequest"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "Authenticated token pair"
          },
          "401": {
            "description": "Invalid credentials"
          }
        }
      }
    },
    "/api/v1/auth/refresh": {
      "post": {
        "summary": "Rotate refresh token",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/RefreshRequest"
              }
            }
          }
        },
        "responses": {
          "200": {
            "description": "New token pair"
          },
          "401": {
            "description": "Invalid refresh token"
          }
        }
      }
    },
    "/api/v1/auth/logout": {
      "post": {
        "summary": "Delete refresh token session",
        "requestBody": {
          "required": true,
          "content": {
            "application/json": {
              "schema": {
                "$ref": "#/components/schemas/RefreshRequest"
              }
            }
          }
        },
        "responses": {
          "204": {
            "description": "Logged out"
          }
        }
      }
    },
    "/api/v1/dev/full-circle": {
      "post": {
        "summary": "Local-only full-circle grocery demo",
        "responses": {
          "201": {
            "description": "Demo tenant/product/order/audit/cache flow created"
          }
        }
      }
    }
  },
  "components": {
    "schemas": {
      "RegisterRequest": {
        "type": "object",
        "required": ["invitation_token", "tenant_name", "tenant_slug", "admin_name", "email", "password"],
        "properties": {
          "invitation_token": { "type": "string" },
          "tenant_name": { "type": "string" },
          "tenant_slug": { "type": "string" },
          "admin_name": { "type": "string" },
          "email": { "type": "string", "format": "email" },
          "password": { "type": "string", "format": "password" }
        }
      },
      "LoginRequest": {
        "type": "object",
        "required": ["tenant_slug", "email", "password"],
        "properties": {
          "tenant_slug": { "type": "string" },
          "email": { "type": "string", "format": "email" },
          "password": { "type": "string", "format": "password" }
        }
      },
      "RefreshRequest": {
        "type": "object",
        "required": ["refresh_token"],
        "properties": {
          "refresh_token": { "type": "string" }
        }
      }
    }
  }
}`

// SwaggerBasicAuth protects local API docs until real superadmin JWT
// authorization exists. It solves the immediate problem of avoiding open docs
// while still keeping Swagger easy to use in a browser during development.
func SwaggerBasicAuth(username, password string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotUser, gotPassword, ok := r.BasicAuth()
			if !ok || !constantTimeEqual(gotUser, username) || !constantTimeEqual(gotPassword, password) || password == "" {
				w.Header().Set("WWW-Authenticate", `Basic realm="Olario Swagger"`)
				writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "swagger access requires superadmin credentials"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func swaggerUIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(swaggerHTML))
}

func openAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(openAPISpec))
}

func constantTimeEqual(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
