package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	database "github.com/nightx1x/ecommerce/interval/db"
	models "github.com/nightx1x/ecommerce/interval/domain"
	repository "github.com/nightx1x/ecommerce/interval/repository/postgres"
	services "github.com/nightx1x/ecommerce/interval/service/product"
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// isRunningInDocker перевіряє, чи запущено застосунок у Docker контейнері
func isRunningInDocker() bool {
	_, err := os.Stat("/.dockerenv")
	return err == nil
}

// loadDBConfig формує конфігурацію підключення до БД
func loadDBConfig() database.Config {
	_ = godotenv.Load() // ігноруємо помилку — для Docker це нормально

	dbHost := getEnv("DB_HOST", "localhost")

	// Якщо Docker відсутній, а DB_HOST вказує на "db" — підміняємо на localhost
	if dbHost == "db" && !isRunningInDocker() {
		dbHost = "localhost"
	}

	return database.Config{
		Host:     dbHost,
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", ""),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "ecommerce_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}
func main() {
	log.Println("🚀 Запуск сервісу ecommerce-api...")

	dbConfig := loadDBConfig()
	log.Printf("🔌 Підключення до бази: %s@%s:%s/%s",
		dbConfig.User, dbConfig.Host, dbConfig.Port, dbConfig.DBName)

	// Ініціалізація підключення
	db, err := database.NewDB(dbConfig)
	if err != nil {
		log.Fatalf("❌ Помилка підключення до БД: %v", err)
	}
	defer db.Close()

	log.Println("✅ З'єднання з базою даних встановлено!")

	// Репозиторій продуктів
	productRepo := repository.NewProductRepository(db)

	productSrv := services.NewService(productRepo)

	ctx := context.Background()

	// Виведення списку товарів
	log.Println("📦 Отримання списку товарів...")

	filter := models.ListFilter{Limit: 100}
	products, err := productRepo.List(ctx, &filter)
	if err != nil {
		log.Fatalf("❌ Помилка отримання товарів: %v", err)
	}

	if len(products) == 0 {
		log.Println("⚠️  База даних порожня — товарів не знайдено.")
		return
	}

	fmt.Printf("\n✅ Знайдено товарів: %d\n", len(products))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i, p := range products {
		fmt.Printf("[%02d] 🆔 %s\n", i+1, p.ID)
		fmt.Printf("     📦 Назва: %s\n", p.Name)
		fmt.Printf("     💰 Ціна: %.2f грн\n", p.Price)
		fmt.Printf("     📊 На складі: %d шт.\n", p.Stock)
		if p.Description != nil && *p.Description != "" {
			fmt.Printf("     📝 Опис: %s\n", *p.Description)
		}
		if p.CategoryID != nil {
			fmt.Printf("     🏷  Категорія ID: %s\n", *p.CategoryID)
		}
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	}

	log.Println("✨ Програма завершена успішно!")
}
