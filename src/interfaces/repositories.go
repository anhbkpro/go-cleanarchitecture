package interfaces

import (
	"fmt"

	"github.com/anhbkpro/go-cleanarchitecture/src/domain"
	"github.com/anhbkpro/go-cleanarchitecture/src/usercases"
)

// define Interfaces here
type DbHandler interface {
	Execute(statement string)
	Query(statement string) Row
}

type Row interface {
	Scan(dest ...interface{})
	Next() bool
}

type DbRepo struct {
	dbHandlers map[string]DbHandler
	dbHandler  DbHandler
}

type DbUserRepo DbRepo
type DbCustomerRepo DbRepo
type DbOrderRepo DbRepo
type DbItemRepo DbRepo

func NewDbUserRepo(dbHandlers map[string]DbHandler) *DbUserRepo {
	dbUserRepo := new(DbUserRepo)
	dbUserRepo.dbHandlers = dbHandlers
	dbUserRepo.dbHandler = dbHandlers["DbUserRepo"]
	return dbUserRepo
}

// UserRepository implementation
func (repo *DbUserRepo) Store(user usercases.User) {
	isAdmin := "no"
	if user.IsAdmin {
		isAdmin = "yes"
	}
	repo.dbHandler.Execute(fmt.Sprintf("INSERT INTO users (id, customer_id, is_admin) VALUES ('%d', '%d', '%v')", user.Id, user.Customer.Id, isAdmin))
	customerRepo := NewDbCustomerRepo(repo.dbHandlers)
	customerRepo.Store(user.Customer)
}

// DbUserRepo.FindById().
// It’s a good example to illustrate that in our architecture,
// interfaces really are all about transforming data from one layer to the next.
// FindById reads database rows and produces domain and usescases entities.
func (repo *DbUserRepo) FindById(id int) usercases.User {
	row := repo.dbHandler.Query(fmt.Sprintf("SELECT is_admin, customer_id FROM users WHERE id = '%d' LIMIT 1", id))
	var isAdmin string
	var customerId int
	row.Next()
	row.Scan(&isAdmin, customerId)
	customerRepo := NewDbCustomerRepo(repo.dbHandlers)
	u := usercases.User{Id: id, Customer: customerRepo.FindById(customerId)}

	// I have deliberately made the database representation of the
	// User.IsAdmin attribute more complicated than neccessary,
	// by storing it as “yes” and “no” varchars in the database.
	// In the usecases entity User, it’s represented as a boolean value of course.
	u.IsAdmin = false
	if isAdmin == "yes" {
		u.IsAdmin = true // Bridging the gap of these very different representations is the job of the repository.
	}
	return u
}

func NewDbCustomerRepo(dbHandlers map[string]DbHandler) *DbCustomerRepo {
	dbCustomerRepo := new(DbCustomerRepo)
	dbCustomerRepo.dbHandlers = dbHandlers
	dbCustomerRepo.dbHandler = dbHandlers["DbCustomerRepo"]
	return dbCustomerRepo
}

func (repo *DbCustomerRepo) Store(customer domain.Customer) {
	repo.dbHandler.Execute(fmt.Sprintf("INSERT INTO customers (id, name) VALUES ('%d', '%v')", customer.Id, customer.Name))
}

func (repo *DbCustomerRepo) FindById(id int) domain.Customer {
	row := repo.dbHandler.Query(fmt.Sprintf("SELECT name FROM customers WHERE id = '%d' LIMIT 1", id))
	var name string
	row.Next()
	row.Scan(&name)
	return domain.Customer{Id: id, Name: name}
}

func NewDbOrderRepo(dbHandlers map[string]DbHandler) *DbOrderRepo {
	dbOrderRepo := new(DbOrderRepo)
	dbOrderRepo.dbHandlers = dbHandlers
	dbOrderRepo.dbHandler = dbHandlers["DbOrderRepo"]
	return dbOrderRepo
}

func (repo *DbOrderRepo) Store(order domain.Order) {
	repo.dbHandler.Execute(fmt.Sprintf("INSERT INTO orders (id, customer_id) VALUES ('%d', '%v')", order.Id, order.Customer.Id))
	for _, item := range order.Items {
		repo.dbHandler.Execute(fmt.Sprintf("INSERT INTO items2orders (item_id, order_id) VALUES ('%d', '%d')", item.Id, order.Id))
	}
}

func (repo *DbOrderRepo) FindById(id int) domain.Order {
	row := repo.dbHandler.Query(fmt.Sprintf("SELECT customer_id FROM orders WHERE id = '%d' LIMIT 1", id))
	var customerId int
	row.Next()
	row.Scan(&customerId)
	customerRepo := NewDbCustomerRepo(repo.dbHandlers)
	order := domain.Order{Id: id, Customer: customerRepo.FindById(customerId)}
	var itemId int
	itemRepo := NewDbItemRepo(repo.dbHandlers)
	row = repo.dbHandler.Query(fmt.Sprintf("SELECT item_id FROM items2orders WHERE order_id = '%d'", order.Id))
	for row.Next() {
		row.Scan(&itemId)
		order.Add(itemRepo.FindById(itemId))
	}
	return order
}

func NewDbItemRepo(dbHandlers map[string]DbHandler) *DbItemRepo {
	dbItemRepo := new(DbItemRepo)
	dbItemRepo.dbHandlers = dbHandlers
	dbItemRepo.dbHandler = dbHandlers["DbItemRepo"]
	return dbItemRepo
}

func (repo *DbItemRepo) Store(item domain.Item) {
	available := "no"
	if item.Available {
		available = "yes"
	}
	repo.dbHandler.Execute(fmt.Sprintf("INSERT INTO items (id, name, value, available) VALUES ('%d', '%v', '%f', '%v')", item.Id, item.Name, item.Value, available))
}

func (repo *DbItemRepo) FindById(id int) domain.Item {
	row := repo.dbHandler.Query(fmt.Sprintf("SELECT name, value, available FROM items WHERE id = '%d' LIMIT 1", id))
	var name string
	var value float64
	var available string
	row.Next()
	row.Scan(&name, &value, &available)
	item := domain.Item{Id: id, Name: name, Value: value}
	item.Available = false
	if available == "yes" {
		item.Available = true
	}
	return item
}
