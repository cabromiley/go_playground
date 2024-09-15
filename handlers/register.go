package handlers

import (
	"database/sql"
	"log"
	"net/http"

	"cabromiley.classes/utils"
	"golang.org/x/crypto/bcrypt"
)

// Register handler - serves the registration form and handles form submission
func Register(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "GET" {
		// Serve the registration form
		err := tmpl.ExecuteTemplate(w, "layout.html", map[string]interface{}{
			"Page": "Register",
		})
		if err != nil {
			log.Println("Failed to render template for registration form:", err)
		}
	} else if r.Method == "POST" {
		// Handle form submission
		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Basic validation
		if name == "" || email == "" || password == "" {
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		if !utils.IsValidEmail(email) {
			http.Error(w, "Invalid email format", http.StatusBadRequest)
			return
		}

		// Hash the password using bcrypt
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Failed to hash password:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Insert the user into the database
		stmt, err := db.Prepare("INSERT INTO users (name, email, password, role) VALUES (?, ?, ?, ?)")
		if err != nil {
			log.Println("Failed to prepare statement for inserting user:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(name, email, hashedPassword, "unverified")
		if err != nil {
			log.Println("Failed to insert user into the database:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Redirect to the login page or show success message
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}