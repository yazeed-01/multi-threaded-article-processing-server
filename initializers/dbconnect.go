package initializers

import (
	"log"
	"mta/models"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
	ES *elasticsearch.Client
)

func ConnectDB() {

	dsn := os.Getenv("DATABASE_URL")
	var err error

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v\n", err)
	}

	err = DB.AutoMigrate(
		&models.Article{},
		&models.Status{},
	)
	if err != nil {
		log.Fatalf("Unable to migrate the database schema: %v\n", err)
	}

	log.Println("Connected to the database and migrated schema successfully")
}

func InitElasticsearch() {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic("failed to create Elasticsearch client")
	}

	ES = es
}
