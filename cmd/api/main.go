package main

import (
	"context"
	"time"

	"github.com/backend-production-go-1/internal/auth"
	"github.com/backend-production-go-1/internal/db"
	"github.com/backend-production-go-1/internal/env"
	"github.com/backend-production-go-1/internal/mailer"
	"github.com/backend-production-go-1/internal/ratelimiter"
	"github.com/backend-production-go-1/internal/store"
	"github.com/backend-production-go-1/internal/store/cache"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

const version = "1.1.0"

//	@title			GopherSocial API
//	@description	API for GopherSocial, a social network for gophers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	cfg := config{
		addr:        env.GetString("ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:5173"),
		db: dbConfig{
			addr:         env.GetString("DB_MIGRATOR_ADDR", "postgres://admin:adminpassword@localhost/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		redisCfg: redisConfig{
			addr:    env.GetString("REDIS_ADDR", "localhost:6379"),
			pw:      env.GetString("REDIS_PW", ""),
			db:      env.GetInt("REDIS_DB", 0),
			enabled: env.GetBool("REDIS_ENABLED", false), // if we want to use redis turn to true --> off false
		},
		env: env.GetString("ENV", "development"),
		mail: mailConfig{
			exp:       time.Hour * 24 * 3, // 3 days
			fromEmail: env.GetString("FROM_EMAIL", ""),
			// sendGrid: sendGridConfig{
			//  apiKey: env.GetString("SENDGRID_API_KEY", ""),
			// },
			mailTrap: mailTrapConfig{
				apiKey: env.GetString("MAILTRAP_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
				exp:    time.Hour * 24 * 3, // 3 days
				iss:    "gophersocial",
			},
		},
		rateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.GetInt("RATELIMITER_REQUESTS_COUNT", 20),
			TimeFrame:            time.Second * 5,
			Enabled:              env.GetBool("RATE_LIMITER_ENABLED", true),
		},
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Main Database
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		logger.Fatal(err)
	}
	// liq/pq
	// db, err := sql.Open("postgres", cfg.db.addr)
	// if err != nil {
	// 	panic(err)
	// }
	defer db.Close()
	logger.Info("database connection pool established")
	//pgx
	// db, err := pgx.Connect(context.Background(), cfg.db.addr)
	// if err != nil {
	// 	panic(err)
	// }
	// defer db.Close(context.Background())
	// logger.Info("database connection pool established")

	// Cache
	//check param
	// log.Println(cfg.redisCfg.enabled)
	// log.Println(cfg.redisCfg.addr)
	// log.Println(cfg.redisCfg.db)
	// log.Println(cfg.redisCfg.pw)

	var rdb *redis.Client

	if cfg.redisCfg.enabled {
		rdb = redis.NewClient(&redis.Options{
			Addr:     cfg.redisCfg.addr,
			Password: cfg.redisCfg.pw,
			DB:       cfg.redisCfg.db,
		})
		_, err := rdb.Ping(context.Background()).Result()
		if err != nil {
			logger.Fatal("Failed to initialize Redis client:", err)
		}
		logger.Info("redis cache connection established")
		defer rdb.Close()
	}

	// Rate limiter
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	// Mailer
	// mailer := mailer.NewSendgrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail)
	mailtrap, err := mailer.NewMailTrapClient(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	// Authenticator
	jwtAuthenticator := auth.NewJWTAuthenticator(
		cfg.auth.token.secret,
		cfg.auth.token.iss,
		cfg.auth.token.iss,
	)

	store := store.NewPostgresStorage(db)
	cacheStorage := cache.NewRedisStorage(rdb)

	app := &application{
		config:        cfg,
		store:         store,
		cacheStorage:  cacheStorage,
		logger:        logger,
		mailer:        mailtrap,
		authenticator: jwtAuthenticator,
		rateLimiter:   rateLimiter,
	}

	// Metrics collected
	// expvar.NewString("version").Set(version)
	// expvar.Publish("database", expvar.Func(func() any {
	// 	return db.Stats()
	// }))
	// expvar.Publish("goroutines", expvar.Func(func() any {
	// 	return runtime.NumGoroutine()
	// }))

	mux := app.mount()

	logger.Fatal(app.run(mux))
}
