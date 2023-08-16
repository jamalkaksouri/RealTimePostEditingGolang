package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

type Product struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	StockQuantity int       `json:"stockQuantity"`
	Version       uuid.UUID `json:"version"`
	AddToCart     bool      `json:"addToCart"`
}

type dbInstance struct {
	Db *sql.DB
}

type serverState struct {
	isUpdate bool
	msg      string
	quantity int
}

func main() {
	db, err := initDB()
	if err != nil {
		log.Fatal("Failed to initialize DB:", err)
	}

	app := &dbInstance{Db: db}
	serverState := &serverState{msg: "", quantity: 0, isUpdate: true}

	route := setupRouter(app, serverState)
	route.Run(":6835")
}

func setupRouter(app *dbInstance, serverState *serverState) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	route := gin.Default()

	route.GET("/", handleIndex)
	route.GET("/sse", handleSSE(app, serverState))
	route.POST("/checkout", func(c *gin.Context) {
		handleCheckout(c, app, serverState)
	})

	return route
}

func handleIndex(c *gin.Context) {
	http.ServeFile(c.Writer, c.Request, "index.html")
}

func handleSSE(app *dbInstance, st *serverState) gin.HandlerFunc {
	return func(c *gin.Context) {
		setupSSEHeaders(c)

		ticker := time.NewTicker(time.Second * 1)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				sendSSEMessage(c.Writer, formatTimeMessage())
				stockText, quantity := getStockStatus(app, st)
				sendSSEMessage(c.Writer, formatStockMessage(stockText, quantity))
			case <-c.Writer.CloseNotify():
				return
			}
		}
	}
}

func setupSSEHeaders(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
}

func sendSSEMessage(w http.ResponseWriter, message string) {
	_, _ = fmt.Fprintf(w, "data: %s\n\n", message)
	w.(http.Flusher).Flush()
}

func initDB() (*sql.DB, error) {
	loadConfig()

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("db.host"), viper.GetInt("db.port"), viper.GetString("db.user"), viper.GetString("db.password"), viper.GetString("db.name"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Connected to the database")
	return db, nil
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Error reading config file:", err)
	}
}

func getStockStatus(app *dbInstance, st *serverState) (string, int) {
	if st.isUpdate {
		state, quantity := isOutOfStock(app)
		if state {
			st.msg = "Product X is out of stock"
			st.quantity = quantity
			st.isUpdate = false
			return st.msg, st.quantity
		} else {
			st.msg = "In stock:"
			st.quantity = quantity
			st.isUpdate = false
		}
	}
	return st.msg, st.quantity
}

func formatStockMessage(stockStatus string, quantity int) string {
	if quantity == 0 {
		return fmt.Sprintf("{\"event\": \"isStock\", \"data\": \"%s\"}", stockStatus)
	}
	return fmt.Sprintf("{\"event\": \"isStock\", \"data\": \"%s %d\"}", stockStatus, quantity)
}

func formatTimeMessage() string {
	currentTime := time.Now().Format("15:04:05")
	return fmt.Sprintf("{\"event\": \"time\", \"data\": \"%s\"}", currentTime)
}

func isOutOfStock(app *dbInstance) (bool, int) {
	var quantity int
	err := app.Db.QueryRow(`
 			SELECT stock_quantity
 			FROM products
 			WHERE name = $1;
		`, "Product X").Scan(&quantity)
	if err != nil {
		log.Println(err)
	}
	if quantity < 1 {
		return true, quantity
	}
	return false, quantity
}

func handleCheckout(c *gin.Context, app *dbInstance, st *serverState) {
	if c.Request.Method != http.MethodPost {
		c.AbortWithStatus(http.StatusMethodNotAllowed)
		return
	}

	productID := 1

	state, _ := isOutOfStock(app)
	if !state {
		tx, err := app.Db.Begin()
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		defer tx.Rollback()

		_, err = tx.Exec(`
		UPDATE products
		SET stock_quantity = stock_quantity - 1,
			version = uuid_generate_v4()
		WHERE id = $1;
	`, productID)
		if err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		if err := tx.Commit(); err != nil {
			log.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		} else {
			st.isUpdate = true
		}
	}
}
