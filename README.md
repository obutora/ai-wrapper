# Go GraphQL Template

Go言語とGraphQLを使用したバックエンドサーバーのテンプレートプロジェクトです。このテンプレートは、GraphQLサーバーを迅速に構築するための基盤として使用できます。

## 主な機能

- **GraphQLサーバー**: [gqlgen](https://github.com/99designs/gqlgen)を使用したGraphQLサーバー実装
- **エンティティフレームワーク**: [ent](https://entgo.io/)を使用したORM機能
- **データベース接続**: TiDB/MySQLへの接続サポート
- **GraphQLプレイグラウンド**: 開発時に便利なGraphQLプレイグラウンド
- **コード自動生成**: スキーマからのコード自動生成
- **Dockerサポート**: コンテナ化対応

## 技術スタック

- [Go](https://golang.org/) - プログラミング言語
- [gqlgen](https://github.com/99designs/gqlgen) - GraphQLサーバー実装
- [ent](https://entgo.io/) - エンティティフレームワーク
- [chi](https://github.com/go-chi/chi) - HTTPルーターとミドルウェア
- [MySQL/TiDB](https://www.mysql.com/) - データベース

## セットアップ

### 前提条件

- Go 1.23以上
- MySQL/TiDBデータベース

### インストール

1. リポジトリをクローン

```bash
git clone https://github.com/yourusername/go_graphql_template.git
cd go_graphql_template
```

2. 依存関係のインストール

```bash
go mod download
```

3. 環境変数の設定

`.env`ファイルを作成し、以下の環境変数を設定します：

```
TIDB_CONNECT=user:password@tcp(host:port)/database?tls=tidb
```

4. コードの生成

```bash
go generate ./...
```

5. サーバーの起動

```bash
go run cmd/server/main.go
```

サーバーが起動したら、[http://localhost:8080](http://localhost:8080)でGraphQLプレイグラウンドにアクセスできます。

## Docker

Dockerを使用してサーバーを実行することもできます：

```bash
docker build -t go-graphql-template .
docker run -p 8080:8080 --env-file .env go-graphql-template
```

## 既存プロジェクトへの統合

このテンプレートを既存のGitリポジトリに組み込むには、git subtreeを使用することができます。これにより、独立したリポジトリとして管理されているこのテンプレートを、既存プロジェクトのサブディレクトリとして追加できます。

### git subtreeを使った統合方法

1. 既存のリポジトリのルートディレクトリに移動します：

```bash
cd your-existing-repository
```

2. git subtreeコマンドを使用してこのリポジトリを追加します：

```bash
git subtree add --prefix=<サブディレクトリ名> <リポジトリURL> <ブランチ名> --squash
```

例えば：

```bash
git subtree add --prefix=graphql-server https://github.com/obutora/go_graphql_template.git main --squash
```

このコマンドは、go_graphql_templateリポジトリのmainブランチを、既存リポジトリの`graphql-server`ディレクトリに追加します。`--squash`オプションは、追加されるリポジトリの履歴を1つのコミットに圧縮します。

### 更新方法

統合したリポジトリを最新バージョンに更新するには：

```bash
git subtree pull --prefix=<サブディレクトリ名> <リポジトリURL> <ブランチ名> --squash
```

例：

```bash
git subtree pull --prefix=graphql-server https://github.com/obutora/go_graphql_template.git main --squash
```

### 注意点

- git subtreeを使用すると、追加したリポジトリのコードは既存リポジトリの一部となります
- 統合後は必要に応じて設定ファイルやパスを調整する必要があるかもしれません
- 大規模なリポジトリを統合する場合、初回のsubtree addコマンドの実行に時間がかかる場合があります

## プロジェクト構造

```
.
├── cmd/                  # コマンドラインアプリケーション
│   └── server/           # サーバーエントリーポイント
├── ent/                  # entの設定
├── generated/            # 自動生成されたコード
│   └── ent/              # entによって生成されたエンティティ
├── graph/                # GraphQL関連のコード
│   ├── generated/        # gqlgenによって生成されたコード
│   ├── model/            # GraphQLモデル
│   └── resolver/         # GraphQLリゾルバー
├── infra/                # インフラストラクチャコード
├── schema/               # スキーマ定義
│   ├── *.graphql         # GraphQLスキーマ
│   └── *.go              # entスキーマ
├── .env                  # 環境変数（gitignore対象）
├── Dockerfile            # Dockerビルド設定
├── generate.go           # コード生成設定
├── go.mod                # Goモジュール定義
├── go.sum                # Goモジュールのチェックサム
└── gqlgen.yml            # gqlgen設定
```

## カスタマイズ

### 新しいエンティティの追加

1. `schema/`ディレクトリに新しいentスキーマを作成します：

```go
// schema/user.go
package schema

import (
    "entgo.io/contrib/entgql"
    "entgo.io/ent"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/field"
)

type User struct {
    ent.Schema
}

func (User) Fields() []ent.Field {
    return []ent.Field{
        field.String("name"),
        field.String("email").Unique(),
    }
}

func (User) Annotations() []schema.Annotation {
    return []schema.Annotation{
        entgql.QueryField(),
        entgql.RelayConnection(),
    }
}
```

2. `schema/`ディレクトリに新しいGraphQLスキーマを作成します：

```graphql
// schema/users.graphql
extend type Mutation {
    createUser(name: String!, email: String!): User!
}
```

3. コードを生成します：

```bash
go generate ./...
```

4. リゾルバーを実装します：

```go
// graph/resolver/users.resolvers.go
func (r *mutationResolver) CreateUser(ctx context.Context, name string, email string) (*ent.User, error) {
    return r.Client.User.Create().
        SetName(name).
        SetEmail(email).
        Save(ctx)
}
```

### データベース接続の変更

`cmd/server/main.go`ファイルでデータベース接続文字列を変更します：

```go
connectionString := "user:password@tcp(localhost:3306)/database?parseTime=true"
client, err := ent.Open("mysql", connectionString)
```

## コントリビューション

1. このリポジトリをフォーク
2. 機能ブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチをプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

## ライセンス

MIT License
