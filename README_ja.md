# go-todo-app

## 概要

Go でバックエンド専用の ToDo API を作りながら学ぶリポジトリです。API と Postgres をローカルで動かすために必要なものはすべて Docker に含めているので、Docker が入っていれば開発を開始できます。

## 必要環境

- Docker（24.x 以上を推奨）
- Docker Compose プラグイン
- ホスト側の空きポート: `8080`（API）, `5432`（Postgres）

## クイックスタート

```bash
docker compose up --build
```

- API: <http://localhost:8080>
- Postgres: `localhost:5432`（ユーザー `todo`, パスワード `todo`）
- 初回起動時は `db/init/` 内の SQL が実行され `todos` テーブルが作成されます。

`CTRL+C` で停止し、ローカル DB データも削除したい場合は `docker compose down -v` を実行してください。

## ディレクトリ構成

- `cmd/server/`: アプリのエントリーポイント。設定読み込み、DB 接続、HTTP サーバ起動を担当
- `internal/config`: 環境変数の読み取りとデフォルト設定
- `internal/database`: Postgres への接続ヘルパー
- `internal/httpserver`: HTTP サーバ、ミドルウェア、ルーティング
- `internal/httpx`: JSON レスポンスなどの HTTP ヘルパー
- `internal/todo`: Todo モデル、リポジトリ、HTTP ハンドラ
- `db/init/`: Postgres 起動時に実行される初期化 SQL

## Docker での開発

`docker-compose.yml` の `api` サービスは `Dockerfile` の `dev` ステージを利用します。ソースコードはバインドマウントされ、[`air`](https://github.com/air-verse/air) がファイル変更時に自動で再ビルド・再起動します。

便利なコマンド:

```bash
# API ログを追いかける
docker compose logs -f api

# Go のテストを実行
docker compose exec api go test ./...

# 依存を追加した後に API イメージを再ビルド
docker compose build api
```

> 💡 最初のビルドで Go モジュールがコンテナ内に取得され、ローカルに `go.sum` が生成されます。依存を追加した場合は `docker compose exec api go mod tidy` を実行して整備してください。

## 設定項目

| 変数名           | デフォルト                                                    | 説明                         |
| ---------------- | ------------------------------------------------------------- | ---------------------------- |
| `PORT`           | `8080`                                                         | API が待ち受けるポート       |
| `DATABASE_URL`   | `postgres://todo:todo@localhost:5432/todo?sslmode=disable`     | Postgres への接続文字列     |
| `POSTGRES_USER`  | `todo`（docker compose 内のみ）                                | DB ユーザー名                |
| `POSTGRES_DB`    | `todo`（docker compose 内のみ）                                | データベース名               |

`docker-compose.yml` で上書きするか、`docker compose --env-file`、もしくは Go バイナリを直接動かす際にシェルから設定してください。
