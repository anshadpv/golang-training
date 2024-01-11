package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"text/template"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// TemplateData represents the structure for template content.
type TemplateData struct {
	Name      string                 `json:"name"`
	Content   string                 `json:"content"`
	FieldData map[string]interface{} `json:"fieldData"`
}

// TemplateStorage represents a storage for templates.
type TemplateStorage map[string]TemplateData

// InMemoryDB is an in-memory storage for templates.
var InMemoryDB TemplateStorage

// MySQLDB is a MySQL database connection.
var MySQLDB *sql.DB

// RedisClient is a Redis client connection.
var RedisClient *redis.Client

var mu sync.Mutex

func init() {
	// Initialize in-memory storage.
	InMemoryDB = make(TemplateStorage)

	// Initialize MySQL database connection.
	var err error
	MySQLDB, err = sql.Open("mysql", "root:msf@12345@tcp(127.0.0.1:3306)/class")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MySQL !!")

	// Initialize Redis client connection.
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	fmt.Println("Connected to Redis !!")
}

func main() {
	r := mux.NewRouter()

	// Define API endpoints.
	r.HandleFunc("/create", createHandler).Methods("POST")
	r.HandleFunc("/delete/{template_name}", deleteHandler).Methods("DELETE")
	r.HandleFunc("/update/{template_name}", updateHandler).Methods("PUT")
	r.HandleFunc("/test", getAllTemplatesHandler).Methods("GET")
	r.HandleFunc("/execute/{template_name}", executeTemplateHandler).Methods("POST")

	// Start the server.
	fmt.Println("Listening on server : 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var templateData TemplateData

	// Decode the request body into a TemplateData struct.
	err := json.NewDecoder(r.Body).Decode(&templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the template name from the URL path parameter.
	//vars := mux.Vars(r)
	templateName := templateData.Name

	// Create a new template with the given data.
	newTemplate := TemplateData{
		Name:    templateData.Name,
		Content: templateData.Content,
	}

	// Store data in in-memory storage.
	mu.Lock()
	InMemoryDB[templateName] = newTemplate
	mu.Unlock()

	// Store data in MySQL database.
	_, err = MySQLDB.Exec("INSERT INTO templates (name, content) VALUES (?, ?)", templateName, newTemplate.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Store data in Redis.
	err = RedisClient.Set(r.Context(), templateName, newTemplate.Content, 0).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success message.
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Template %s created successfully", templateName)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// Extract template name from the URL path parameter.
	vars := mux.Vars(r)
	templateName := vars["template_name"]

	//checking the existance
	var existingTemplateContent string
	err := MySQLDB.QueryRow("SELECT content FROM templates WHERE name = ?", templateName).Scan(&existingTemplateContent)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, fmt.Sprintf("Template %s not found", templateName), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Remove data from in-memory storage.
	mu.Lock()
	delete(InMemoryDB, templateName)
	mu.Unlock()

	// Remove data from MySQL database.
	_, err = MySQLDB.Exec("DELETE FROM templates WHERE name = ?", templateName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Remove data from Redis.
	err = RedisClient.Del(r.Context(), templateName).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success message.
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Template %s deleted successfully", templateName)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	// Extract template name from the URL path parameter.
	vars := mux.Vars(r)
	templateName := vars["template_name"]

	// Decode the request body into a TemplateData struct.
	var updatedTemplate TemplateData
	err := json.NewDecoder(r.Body).Decode(&updatedTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the template exists.
	var existingTemplateContent string
	err = MySQLDB.QueryRow("SELECT content FROM templates WHERE name = ?", templateName).Scan(&existingTemplateContent)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, fmt.Sprintf("Template %s not found", templateName), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update in-memory storage.
	mu.Lock()
	InMemoryDB[templateName] = updatedTemplate
	mu.Unlock()

	// Update in MySQL database.
	_, err = MySQLDB.Exec("UPDATE templates SET content = ? WHERE name = ?", updatedTemplate.Content, templateName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update in Redis.
	err = RedisClient.Set(r.Context(), templateName, updatedTemplate.Content, 0).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Respond with success message.
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Template %s updated successfully", templateName)
}

func getAllTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve all templates from MySQL database.
	rows, err := MySQLDB.Query("SELECT name, content FROM templates")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var templates []TemplateData

	// Iterate over the result set and populate the templates slice.
	for rows.Next() {
		var template TemplateData
		err := rows.Scan(&template.Name, &template.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		templates = append(templates, template)
	}

	// Respond with the list of templates in JSON format.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

func executeTemplateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateName := vars["template_name"]

	// Retrieve the template from MySQL database.
	var templateContent string
	err := MySQLDB.QueryRow("SELECT content FROM templates WHERE name = ?", templateName).Scan(&templateContent)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, fmt.Sprintf("Template %s not found", templateName), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Retrieve the field data from the request body.
	var requestPayload struct {
		FieldData map[string]interface{} `json:"fieldData"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Replace placeholders in the template content with actual values.
	executedContent, err := executeTemplate(templateContent, requestPayload.FieldData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with the executed template content.
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(executedContent))
}

func executeTemplate(templateContent string, data map[string]interface{}) (string, error) {
	tmpl, err := template.New("template").Parse(templateContent)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}
