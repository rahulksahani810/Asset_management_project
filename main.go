package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Admin struct {
	Id         uuid.UUID  `json:"id,omitempty" db:"id"`
	Name       string     `json:"name,omitempty" db:"name"`
	Email      string     `json:"email,omitempty" db:"email"`
	Password   string     `json:"password,omitempty" db:"password"`
	CreatedAt  time.Time  `json:"createdAt,omitempty" db:"created_at"`
	ArchivedAt *time.Time `json:"archivedAt,omitempty" db:"archived_at"`
}

type Admin_session struct {
	Id         uuid.UUID  `json:"id,omitempty" db:"id"`
	AdminId    uuid.UUID  `json:"Adminid,omitempty" db:"Admin_id"`
	CreatedAt  time.Time  `json:"createdAt,omitempty" db:"created_at"`
	ArchivedAt *time.Time `json:"archivedAt,omitempty" db:"archived_at"`
}

type Asset struct {
	Id         uuid.UUID  `json:"id,omitempty" db:"id"`
	Model      string     `json:"Model,omitempty" db:"Model"`
	Company    string     `json:"company,omitempty" db:"company"`
	CreatedAt  time.Time  `json:"createdAt,omitempty" db:"created_at"`
	ArchivedAt *time.Time `json:"archivedAt,omitempty" db:"archived_at"`
}

type Emplyoee struct {
	Id         uuid.UUID  `json:"id,omitempty" db:"id"`
	Name       string     `json:"name,omitempty" db:"name"`
	Email      string     `json:"email,omitempty" db:"email"`
	Role       string     `json:"role,omitempty" db:"role"`
	CreatedAt  time.Time  `json:"createdAt,omitempty" db:"created_at"`
	ArchivedAt *time.Time `json:"archivedAt,omitempty" db:"archived_at"`
}

type EmployeeAssetMapping struct {
	Id         uuid.UUID  `json:"id,omitempty" db:"id"`
	AssetId    string     `json:"assetid,omitempty" db:"assetid"`
	EmplyoeeID string     `json:"emplyoeeid,omitempty" db:"emplyoeeid"`
	CreatedAt  time.Time  `json:"createdAt,omitempty" db:"created_at"`
	ArchivedAt *time.Time `json:"archivedAt,omitempty" db:"archived_at"`
}

///******************  ADMIN ************************************************************//////

func getAdmin(w http.ResponseWriter, r *http.Request) {
	var admins []Admin
	SQL := "select id, name, email,password, created_at, archived_at from admin"
	err := DB.Select(&admins, SQL)
	if err != nil {
		fmt.Println("error reading users", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	RespondJSON(w, http.StatusOK, admins)
}

func addAdmin(w http.ResponseWriter, r *http.Request) {
	var body Admin
	if err := ParseBody(r.Body, &body); err != nil {
		RespondJSON(w, http.StatusBadRequest, err)
		return
	}

	SQL := `insert into admin (id, name, email, password ,created_at) values ($1, $2, $3, $4,$5)`
	_, err := DB.Queryx(SQL, uuid.New(), body.Name, body.Email, body.Password, time.Now())
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, err)
		return
	}

	RespondJSON(w, http.StatusCreated, nil)
}

func deleteAdmin(w http.ResponseWriter, r *http.Request) {
	var body Admin
	if err := ParseBody(r.Body, &body); err != nil {
		RespondJSON(w, http.StatusBadRequest, err)
		return
	}
	SQL := `update admin set archived_at = now() where id = $1`
	_, err := DB.Queryx(SQL, body.Id)
	if err != nil {
		fmt.Println("error archiving user", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

	RespondJSON(w, http.StatusOK, nil)
}

// *************************** Assets ************************************************************//////
func getAsset(w http.ResponseWriter, r *http.Request) {
	var assets []Asset
	SQL := "select id, model, company, created_at, archived_at from assets"
	err := DB.Select(&assets, SQL)
	if err != nil {
		fmt.Println("error reading users", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	RespondJSON(w, http.StatusOK, assets)
}
func addAsset(w http.ResponseWriter, r *http.Request) {
	var body Asset
	if err := ParseBody(r.Body, &body); err != nil {
		RespondJSON(w, http.StatusBadRequest, err)
		return
	}

	SQL := `insert into assets (id, model, company, created_at) values ($1, $2, $3, $4)`
	_, err := DB.Queryx(SQL, uuid.New(), body.Model, body.Company, time.Now())
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, err)
		return
	}

	RespondJSON(w, http.StatusCreated, nil)
}

func deleteAsset(w http.ResponseWriter, r *http.Request) {
	var body Asset
	if err := ParseBody(r.Body, &body); err != nil {
		fmt.Println("error parsing body", err)
		RespondJSON(w, http.StatusBadRequest, err)
		return
	}
	fmt.Println("body", body)
	SQL := `update assets set archived_at = $1 where id = $2`
	_, err := DB.Queryx(SQL, time.Now(), body.Id)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, err)
		return
	}
	RespondJSON(w, http.StatusCreated, nil)
}

// *************************** Employee ************************************************************//////
func getEmployee(w http.ResponseWriter, r *http.Request) {
	var employees []Emplyoee
	SQL := "select id, name, email,role, created_at, archived_at from employee"
	err := DB.Select(&employees, SQL)
	if err != nil {
		fmt.Println("error reading users", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	RespondJSON(w, http.StatusOK, employees)
}
func addEmployee(w http.ResponseWriter, r *http.Request) {
	var body Emplyoee
	if err := ParseBody(r.Body, &body); err != nil {
		RespondJSON(w, http.StatusBadRequest, err)
		return
	}

	SQL := `insert into employee (id, name, email, role, created_at) values ($1, $2, $3, $4,$5)`
	_, err := DB.Queryx(SQL, uuid.New(), body.Name, body.Email, body.Role, time.Now())
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, err)
		return
	}

	RespondJSON(w, http.StatusCreated, nil)
}
func deleteEmployee(w http.ResponseWriter, r *http.Request) {
	var body Emplyoee
	if err := ParseBody(r.Body, &body); err != nil {
		fmt.Println("error parsing body", err)
		RespondJSON(w, http.StatusBadRequest, err)
		return
	}
	fmt.Println("body", body)
	SQL := `update employee set archived_at = $1 where id = $2`
	_, err := DB.Queryx(SQL, time.Now(), body.Id)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, err)
		return
	}
	RespondJSON(w, http.StatusCreated, nil)
}

// *************************** Employee Asset Mapping ************************************************************//////
func getEmployeeAssetMapping(w http.ResponseWriter, r *http.Request) {
	var employeeassetmappings []EmployeeAssetMapping
	SQL := "select id, assetid, emplyoeeid, created_at, archived_at from employeeassetmapping"
	err := DB.Select(&employeeassetmappings, SQL)
	if err != nil {
		fmt.Println("error reading users", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	RespondJSON(w, http.StatusOK, employeeassetmappings)
}
func AssignEmployeeAsset(w http.ResponseWriter, r *http.Request) {
	var body EmployeeAssetMapping
	if err := ParseBody(r.Body, &body); err != nil {
		RespondJSON(w, http.StatusBadRequest, err)
		return
	}

	SQL := `insert into employeeassetmapping (id, assetid, emplyoeeid, created_at) values ($1, $2, $3, $4)`
	_, err := DB.Queryx(SQL, uuid.New(), body.AssetId, body.EmplyoeeID, time.Now())
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, err)
		return
	}

	RespondJSON(w, http.StatusCreated, nil)
}
func returnedAseet(w http.ResponseWriter, r *http.Request) {
	var body EmployeeAssetMapping
	if err := ParseBody(r.Body, &body); err != nil {
		fmt.Println("error parsing body", err)
		RespondJSON(w, http.StatusBadRequest, err)
		return
	}
	fmt.Println("body", body)
	SQL := `update employeeassetmapping set archived_at = $1 where id = $2`
	_, err := DB.Queryx(SQL, time.Now(), body.Id)
	if err != nil {
		RespondJSON(w, http.StatusInternalServerError, err)
		return
	}
	RespondJSON(w, http.StatusCreated, nil)
}

//***************************  MAIN ************************************************************//////

var DB *sqlx.DB

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", "5433", "postgres", "local", "assetmanagement")

	var err error
	DB, err = sqlx.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println("Unable to Connect to the Database, ", err)
		return
	}
	err = DB.Ping()
	if err != nil {
		fmt.Println("Ping Panic", err)
		return
	} else {
		fmt.Println("Connected to the Database")
	}
	router := chi.NewRouter()
	//**************************  ADMIN ************************************************************//////

	router.Route("/admin", func(r chi.Router) {
		r.Get("/", getAdmin)
		r.Post("/", addAdmin)
		r.Route("/{id}", func(AdminIdRouter chi.Router) {
			AdminIdRouter.Put("/", deleteAdmin)
		})
	})
	//**************************  ASSETS ************************************************************//////

	router.Route("/asset", func(r chi.Router) {
		r.Get("/", getAsset)
		r.Post("/", addAsset)
		r.Route("/{id}", func(AssetIdRouter chi.Router) {
			AssetIdRouter.Delete("/", deleteAsset)
		})
	})
	//**************************  EMPLOYEE ************************************************************//////
	router.Route("/employee", func(r chi.Router) {
		r.Get("/", getEmployee)
		r.Post("/", addEmployee)
		r.Route("/{id}", func(EmployeeIdRouter chi.Router) {
			EmployeeIdRouter.Delete("/", deleteEmployee)
		})
	})
	//**************************  EMPLOYEE ASSET MAPPING ************************************************************//////
	router.Route("/employeeassetmapping", func(r chi.Router) {
		r.Get("/", getEmployeeAssetMapping)
		r.Post("/", AssignEmployeeAsset)
		r.Route("/{id}", func(EmployeeAssetMappingIdRouter chi.Router) {
			EmployeeAssetMappingIdRouter.Put("/returnedAseet", returnedAseet)
		})

	})

	fmt.Println("starting server at port 8080")
	http.ListenAndServe(":8080", router)

}

func RespondJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.WriteHeader(statusCode)
	if body != nil {
		if err := EncodeJSONBody(w, body); err != nil {
			fmt.Println(fmt.Errorf("failed to respond JSON with error: %+v", err))
		}
	}
}

func EncodeJSONBody(resp http.ResponseWriter, data interface{}) error {
	return json.NewEncoder(resp).Encode(data)
}

func ParseBody(body io.Reader, out interface{}) error {
	err := json.NewDecoder(body).Decode(out)
	if err != nil {
		return err
	}

	return nil
}
