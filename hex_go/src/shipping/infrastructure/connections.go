package infrastructure

import (
	"database/sql"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/streadway/amqp"
)

var (
	dbConnection    *sql.DB
	rabbitConn     *amqp.Connection
	rabbitChannel  *amqp.Channel
)

func InitMySQL() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	mysqlDSN := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName
	db, err := sql.Open("mysql", mysqlDSN)
	if err != nil {
		log.Printf("Error conectando a MySQL: %v", err)
		return nil, err
	}

	// Configurar el pool de conexiones desde variables de entorno
	maxOpenConns := getEnvAsInt("DB_MAX_OPEN_CONNS", 25)
	maxIdleConns := getEnvAsInt("DB_MAX_IDLE_CONNS", 25)
	connMaxLifetime := getEnvAsDuration("DB_CONN_MAX_LIFETIME", 0)

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(connMaxLifetime)

	err = db.Ping()
	if err != nil {
		log.Printf("Error verificando conexión a MySQL: %v", err)
		return nil, err
	}

	dbConnection = db
	log.Printf("✅ Conexión exitosa a MySQL")
	return db, nil
}

func InitRabbitMQ() (*amqp.Connection, *amqp.Channel, error) {
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	maxRetries := getEnvAsInt("RABBITMQ_MAX_RETRIES", 5)
	retryDelay := getEnvAsDuration("RABBITMQ_RETRY_DELAY", 5*time.Second)

	var err error
	var conn *amqp.Connection
	var ch *amqp.Channel

	// Implementar mecanismo de reintentos
	for i := 0; i < maxRetries; i++ {
		conn, err = amqp.Dial(rabbitmqURL)
		if err == nil {
			ch, err = conn.Channel()
			if err == nil {
				rabbitConn = conn
				rabbitChannel = ch
				log.Printf("✅ Conexión exitosa a RabbitMQ")
				return conn, ch, nil
			}
		}
		log.Printf("Intento %d de %d: Error conectando a RabbitMQ: %v", i+1, maxRetries, err)
		time.Sleep(retryDelay)
	}

	return nil, nil, err
}

func CloseConnections() {
	if dbConnection != nil {
		if err := dbConnection.Close(); err != nil {
			log.Printf("Error cerrando conexión MySQL: %v", err)
		} else {
			log.Println("Conexión MySQL cerrada correctamente")
		}
	}

	if rabbitChannel != nil {
		if err := rabbitChannel.Close(); err != nil {
			log.Printf("Error cerrando canal RabbitMQ: %v", err)
		} else {
			log.Println("Canal RabbitMQ cerrado correctamente")
		}
	}

	if rabbitConn != nil {
		if err := rabbitConn.Close(); err != nil {
			log.Printf("Error cerrando conexión RabbitMQ: %v", err)
		} else {
			log.Println("Conexión RabbitMQ cerrada correctamente")
		}
	}
}

// Funciones auxiliares para obtener variables de entorno
func getEnvAsInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultVal
}

func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultVal
}
