package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	"github.com/redis/go-redis/v9"
	"google.golang.org/api/option"

	"trip-history-api/internal/adapters/handlers/http"
	"trip-history-api/internal/adapters/handlers/ws"
	"trip-history-api/internal/adapters/repositories/firestore"
	cache "trip-history-api/internal/adapters/repositories/redis"
	"trip-history-api/internal/core/services"
)

var ctx = context.Background()

func main() {
	// 1. ตั้งค่าเชื่อมต่อ Firebase (Firestore & Auth)
	opt := option.WithCredentialsFile("serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase app: %v\n", err)
	}

	firestoreClient, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error initializing firestore: %v\n", err)
	}
	defer firestoreClient.Close()
	fmt.Println("✅ Firebase Firestore connected successfully!")

	authClient, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error initializing auth: %v\n", err)
	}
	fmt.Println("✅ Firebase Auth connected successfully!")

	messagingClient, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error initializing messaging: %v\n", err)
	}
	fmt.Println("✅ Firebase Messaging connected successfully!")

	// 2. ตั้งค่าเชื่อมต่อ Redis
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	var redisOpt *redis.Options
	if strings.HasPrefix(redisAddr, "redis://") || strings.HasPrefix(redisAddr, "rediss://") {
		var err error
		redisOpt, err = redis.ParseURL(redisAddr)
		if err != nil {
			log.Fatalf("error parsing redis url: %v\n", err)
		}
	} else {
		redisOpt = &redis.Options{
			Addr:     redisAddr,
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       0,
		}
	}

	redisClient := redis.NewClient(redisOpt)

	_, err = redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("error connecting to redis: %v\n", err)
	}
	fmt.Println("✅ Redis connected successfully!")

	// 3. Initialize WebSocket Hub
	trackingHub := ws.NewTrackingHub()
	go trackingHub.Run()

	// 4. Initialize Repositories (Adapters)
	tripRepo := firestore.NewTripRepository(firestoreClient)
	userRepo := firestore.NewUserRepository(firestoreClient, authClient)
	logRepo := firestore.NewLogRepository(firestoreClient)
	jobPlanRepo := firestore.NewJobPlanRepository(firestoreClient)
	locationCache := cache.NewLocationCache(redisClient)

	// 5. Initialize Services (Core)
	tripService := services.NewTripService(tripRepo, locationCache, userRepo)
	userService := services.NewUserService(userRepo)
	logService := services.NewLogService(logRepo)
	notificationService := services.NewNotificationService(messagingClient)
	jobPlanService := services.NewJobPlanService(jobPlanRepo, userRepo, notificationService)
	systemService := services.NewSystemService(userRepo, tripRepo, logRepo, locationCache, jobPlanRepo)

	// 6. Initialize Handlers and Router (Inbound Adapters)
	tripHandler := http.NewTripHandler(tripService, trackingHub)
	userHandler := http.NewUserHandler(userService)
	logHandler := http.NewLogHandler(logService)
	systemHandler := http.NewSystemHandler(systemService)
	jobPlanHandler := http.NewJobPlanHandler(jobPlanService)

	r := http.SetupRouter(tripHandler, userHandler, logHandler, systemHandler, jobPlanHandler, trackingHub)

	// 7. รันเซิร์ฟเวอร์
	fmt.Println("🚀 Server is running on port 8080...")
	ports := os.Getenv("PORT")
	if ports == "" {
		ports = "8080"
	}
	r.Run(":" + ports)
}
