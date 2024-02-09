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

	"github.com/bluele/gcache"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type TemplateData struct {
	Name      string                 `json:"name"`
	Content   string                 `json:"content"`
	FieldData map[string]interface{} `json:"fieldData"`
}

// in-memory storage for templates.
var InMemoryDB gcache.Cache

// MySQL database connection.
var MySQLDB *sql.DB

// Redis client connection.
var RedisClient *redis.Client

var mu sync.Mutex

func init() {
	// Initializing in-memory storage.
	InMemoryDB = gcache.New(20).Build()

	// Initializing MySQL database connection.
	var err error
	MySQLDB, err = sql.Open("mysql", "root:msf@12345@tcp(127.0.0.1:3306)/class")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MySQL !!")

	// Initializing Redis client connection.
	RedisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	fmt.Println("Connected to Redis !!")
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/create", createHandler).Methods("POST")
	r.HandleFunc("/delete/{template_name}", deleteHandler).Methods("DELETE")
	r.HandleFunc("/update/{template_name}", updateHandler).Methods("PUT")
	r.HandleFunc("/test", getAllTemplatesHandler).Methods("GET")
	r.HandleFunc("/execute/{template_name}", executeTemplateHandler).Methods("POST")
	r.HandleFunc("/refresh", refreshHandler).Methods("POST")

	// Starting the server.
	fmt.Println("Listening on server : 8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}

func createHandler(w http.ResponseWriter, r *http.Request) {
	var templateData TemplateData

	// Decoding the request body into a TemplateData struct.
	err := json.NewDecoder(r.Body).Decode(&templateData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//checking if template name is empty
	if strings.TrimSpace(templateData.Name) == "" {
		http.Error(w, "Template name cannot be empty.", http.StatusBadRequest)
		return
	}

	//checking if template data is empty
	if strings.TrimSpace(templateData.Content) == "" {
		http.Error(w, "Template Content cannot be empty.", http.StatusBadRequest)
		return
	}

	var existing string
	err = MySQLDB.QueryRow("SELECT content FROM templates WHERE name = ?", templateData.Name).Scan(&existing)
	switch {
	case err == sql.ErrNoRows:

	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error on selecting from db.")
		return

	default:
		log.Println("Template already exists.")
		http.Error(w, "Template already exists.", http.StatusConflict)
		return
	}

	// Store data in in-memory storage.
	mu.Lock()
	InMemoryDB.Set(templateData.Name, templateData.Content)
	mu.Unlock()

	// Store data in MySQL database.
	_, err = MySQLDB.Exec("INSERT INTO templates (name, content) VALUES (?, ?)", templateData.Name, templateData.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error on inserting to mySQL.")
		return
	}

	// Store data in Redis.
	err = RedisClient.Set(r.Context(), templateData.Name, templateData.Content, 0).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error on inserting to redis.")
		return
	}

	// Respond with success message.
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Template %s created successfully", templateData.Name)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	// Extracting template name from the URL path parameter.
	vars := mux.Vars(r)
	templateName := vars["template_name"]

	//checking the existance
	var existing string
	err := MySQLDB.QueryRow("SELECT content FROM templates WHERE name = ?", templateName).Scan(&existing)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, fmt.Sprintf("Template %s not found", templateName), http.StatusNotFound)
			log.Println("No such template to delete.")
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error on selecting from db.")
		return
	}

	// Removing data from in-memory storage.
	mu.Lock()
	InMemoryDB.Remove(templateName)
	mu.Unlock()

	// Removing data from MySQL database.
	_, err = MySQLDB.Exec("DELETE FROM templates WHERE name = ?", templateName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error on deleting from mySQL.")
		return
	}

	// Removing data from Redis.
	err = RedisClient.Del(r.Context(), templateName).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error on deleting from redis.")
		return
	}

	// Responding with success message.
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Template %s deleted successfully", templateName)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	// Extracting template name from the URL path parameter.
	vars := mux.Vars(r)
	templateName := vars["template_name"]

	// Decoding the request body into a TemplateData struct.
	var updatedTemplate TemplateData
	err := json.NewDecoder(r.Body).Decode(&updatedTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Error on decoding the request body.")
		return
	}

	// Checking if the template exists.
	err = MySQLDB.QueryRow("SELECT content FROM templates WHERE name = ?", templateName).Err()
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, fmt.Sprintf("Template %s not found", templateName), http.StatusNotFound)
			log.Println("No such template to update.")
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error on selecting from db.")
		return
	}

	// Updating in-memory storage.
	mu.Lock()
	InMemoryDB.Set(templateName, updatedTemplate.Content)
	mu.Unlock()

	// Updating in MySQL database.
	_, err = MySQLDB.Exec("UPDATE templates SET content = ? WHERE name = ?", updatedTemplate.Content, templateName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error on updating in mySQL.")
		return
	}

	// Updating in Redis.
	err = RedisClient.Set(r.Context(), templateName, updatedTemplate.Content, 0).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error on updating in redis.")
		return
	}
	// Respond with success message.
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Template %s updated successfully", templateName)
}

func getAllTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieving all templates from MySQL database.
	rows, err := MySQLDB.Query("SELECT name, content FROM templates")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error on selecting from db.")
		return
	}
	defer rows.Close()

	var templates []TemplateData

	// Iterating over the result set and populate the templates slice.
	for rows.Next() {
		var template TemplateData
		err := rows.Scan(&template.Name, &template.Content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error on scanning from db.")
			return
		}
		templates = append(templates, template)
	}

	// Responding with the list of templates in JSON format.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(templates)
}

func executeTemplateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	templateName := vars["template_name"]

	// Retrieving the template from MySQL database.
	var templateContent string
	err := MySQLDB.QueryRow("SELECT content FROM templates WHERE name = ?", templateName).Scan(&templateContent)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, fmt.Sprintf("Template %s not found", templateName), http.StatusNotFound)
			log.Println("No such template to execute.")
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error on selecting from db.")
		return
	}

	// Retrieving the field data from the request body.
	var requestPayload struct {
		FieldData map[string]interface{} `json:"fieldData"`
	}
	err = json.NewDecoder(r.Body).Decode(&requestPayload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Error on decoding from request body.")
		return
	}

	// Replacing placeholders in the template content with actual values.
	executedContent, err := executeTemplate(templateContent, requestPayload.FieldData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error on executing the template.")
		return
	}

	// Responding with the executed template content.
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(executedContent))
}

func executeTemplate(templateContent string, data map[string]interface{}) (string, error) {

	//creating new template
	tmpl, err := template.New("template").Parse(templateContent)
	if err != nil {
		log.Println("Error on creating template.")
		return "", err
	}

	//updating placeholder
	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		log.Println("Error on replacing placeholders.")
		return "", err
	}

	return result.String(), nil
}

func refreshHandler(w http.ResponseWriter, r *http.Request) {
	// Fetching the latest data from the MySQL database.
	rows, err := MySQLDB.Query("SELECT name, content FROM templates")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error on selecting from db.")
		return
	}
	defer rows.Close()

	// Iterating over the result set and populate the refreshedData map.
	for rows.Next() {
		var template TemplateData
		if err := rows.Scan(&template.Name, &template.Content); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println("Error on refreshing the db.")
			return
		}
		//updating the in memory db
		mu.Lock()
		InMemoryDB.Set(template.Name, template.Content)
		mu.Unlock()
	}

	// Responding with a success message.
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Data refreshed successfully")
}
