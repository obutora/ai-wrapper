package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/obutora/go_graphql_template/generated/ent"
	"github.com/obutora/go_graphql_template/graph/generated"
	"github.com/obutora/go_graphql_template/graph/resolver"
	slogchi "github.com/samber/slog-chi"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("[.env] not found | If you showing this message in CloudRun, env is already set by SecretManager because you don't need .env file.")
	}
}

func main() {
	ctx := context.Background()

	mysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: "gateway01.ap-northeast-1.prod.aws.tidbcloud.com",
	})

	// 接続文字列にタイムアウト関連のパラメータを追加
	connectionString := os.Getenv("TIDB_CONNECT") + "&timeout=30s&readTimeout=30s&writeTimeout=30s"

	// 接続の再試行ロジックを実装
	var client *ent.Client
	var err error
	maxRetries := 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		client, err = ent.Open("mysql", connectionString)
		if err == nil {
			break
		}
		lastErr = err
		log.Printf("データベース接続の試行 %d/%d に失敗しました: %v", i+1, maxRetries, err)
		// 再試行前に少し待機
		if i < maxRetries-1 {
			time.Sleep(time.Second * 2)
		}
	}

	if client == nil {
		panic(fmt.Errorf("データベース接続に失敗しました: %w", lastErr))
	}

	// 注意: client.DB()メソッドは存在しないため、接続プールの設定は行わない
	// 代わりに、接続文字列のパラメータで制御する

	defer client.Close()

	if err := client.Schema.Create(ctx); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &resolver.Resolver{
		Client: client,
	}}))

	// authClient, err := initFirebase(ctx)
	// if err != nil {
	// 	panic(fmt.Errorf("Firebase Authクライアントの初期化に失敗しました: %w", err))
	// }

	// Create a slog logger, which:
	//   - Logs to stdout.
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	r := chi.NewRouter()
	// Middleware
	r.Use(slogchi.New(logger))
	r.Use(middleware.Recoverer)
	// r.Use(FirebaseAuthMiddleware(authClient))
	r.Get("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", srv)

	fmt.Println("connect to http://localhost:8080/ for GraphQL playground")

	// NOTE: localhostを付けないとWindows環境で8080portがFirewallを通らない
	log.Fatal(http.ListenAndServe(":"+"8080", r))
}

// // FirebaseAuthMiddleware Firebase認証ミドルウェア
// func FirebaseAuthMiddleware(authClient *auth.Client) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			authHeader := r.Header.Get("Authorization")
// 			if authHeader == "" {
// 				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
// 				return
// 			}

// 			idToken := strings.TrimPrefix(authHeader, "Bearer ")
// 			if idToken == "" {
// 				http.Error(w, "Bearer token is required", http.StatusUnauthorized)
// 				return
// 			}

// 			token, err := authClient.VerifyIDToken(r.Context(), idToken)
// 			if err != nil {
// 				http.Error(w, "Invalid ID token", http.StatusUnauthorized)
// 				return
// 			}

// 			ctx := context.WithValue(r.Context(), "uid", token.UID)
// 			next.ServeHTTP(w, r.WithContext(ctx))
// 		})
// 	}
// }

// func initFirebase(ctx context.Context) (*auth.Client, error) {
// 	// まずBase64エンコードされた認証情報を確認
// 	credBase64 := os.Getenv("FIREBASE_ADMIN_CREDENTIALS_BASE64")
// 	var opt option.ClientOption

// 	if credBase64 != "" {
// 		// Base64エンコードされた認証情報を使用
// 		credJSON, err := base64.StdEncoding.DecodeString(credBase64)
// 		if err != nil {
// 			return nil, fmt.Errorf("Firebase認証情報のBase64デコードに失敗しました: %w", err)
// 		}
// 		opt = option.WithCredentialsJSON(credJSON)
// 	} else {
// 		// 通常のJSON文字列の認証情報を使用
// 		credJSON := os.Getenv("FIREBASE_ADMIN_CREDENTIALS")
// 		if credJSON == "" {
// 			return nil, errors.New("Firebase認証情報が設定されていません")
// 		}

// 		// JSONとして解析できるか確認
// 		var jsonCheck map[string]interface{}
// 		if err := json.Unmarshal([]byte(credJSON), &jsonCheck); err != nil {
// 			return nil, fmt.Errorf("Firebase認証情報がJSON形式ではありません: %w", err)
// 		}

// 		opt = option.WithCredentialsJSON([]byte(credJSON))
// 	}

// 	// Firebaseアプリの初期化
// 	app, err := firebase.NewApp(ctx, nil, opt)
// 	if err != nil {
// 		return nil, fmt.Errorf("Firebaseアプリの初期化に失敗しました: %w", err)
// 	}

// 	// Auth clientの取得
// 	authClient, err := app.Auth(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("Firebase Authクライアントの初期化に失敗しました: %w", err)
// 	}

// 	return authClient, nil
// }
