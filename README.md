# Saikyo UI

## 概要
Saikyo UIは、Go製のWebアプリケーションとMySQLを組み合わせて、ユーザー登録やかき氷メニュー選択・注文フローを提供するプロジェクトです。Echoフレームワークを用いたAPIと、複数ステップのメニューUI、外部注文API連携を備えています。

## 主要技術
- Go (Echo, GORM)
- MySQL 8
- Docker / Docker Compose
- HTML + CSS + JavaScript (静的アセットは`public/`配下、テンプレートは`views/`配下)

## ディレクトリ構成（一部）
```text
.
├─configuration/   # DB接続や環境変数読み込み
├─db/              # 初期化用SQL
├─functions/       # APIロジック（ユーザー操作など）
├─handlers/        # Web/APIハンドラー
├─public/          # 静的ファイル（CSS, JS, JSON）
├─routes/          # ルーティング定義
├─views/           # Echoテンプレート
├─Dockerfile
└─docker-compose.yml
```

## 前提条件
- Docker および Docker Compose
- （オプション）Go 1.25 以上：ローカルで直接実行したい場合
- MySQLクライアント（動作確認や接続テストを行いたい場合）

## セットアップ
1. リポジトリを取得します。
   ```bash
   git clone https://github.com/U-T-kuroitigo/Saikyo_UI.git
   cd Saikyo_UI
   ```
2. `.env.example` をコピーして `.env` を作成し、必要な値を入力します。
   ```bash
   cp .env.example .env
   ```
3. `.env` 内の値を環境に合わせて更新します。特に下記は必須です。

   | 変数名 | 用途 | デフォルト/例 |
   | --- | --- | --- |
   | `APP_PORT` | アプリケーションの公開ポート | `8080` |
   | `MYSQL_HOST_PORT` | ホスト側のMySQL待受ポート | `33306` |
   | `MYSQL_ROOT_PASSWORD` | MySQL root パスワード | 任意の安全な値 |
   | `MYSQL_DATABASE` | 作成するデータベース名 | 例: `saikyo_db` |
   | `MYSQL_USER` / `MYSQL_PASSWORD` | アプリケーション用MySQLユーザー | 例: `saikyo_user` / `saikyo_pass` |
   | `Server` / `Port` / `User` / `Password` / `Database` | アプリからDBへ接続するための値。Docker Compose利用時は例のままで利用可 | `Server=db`, `Port=3306`, `User=${MYSQL_USER}` など |
   | `STORE_ID` | 外部注文API (https://kakigori-api.fly.dev) にリクエストする際の店舗ID | 発行されたIDを設定 |

   `.env` はアプリコンテナとDBコンテナ双方で読み込まれます。

## 起動方法
### Docker Compose での起動
1. イメージをビルドしてコンテナを立ち上げます。
   ```bash
   docker compose up -d --build
   ```
2. コンテナの状態を確認します。
   ```bash
   docker compose ps
   ```
3. ヘルスチェックエンドポイントで起動を確認します。
   ```bash
   curl http://localhost:8080/health
   ```
   `HTTP/1.1 200 OK` が返ればアプリは正常です。`APP_PORT` を変更している場合はURLのポートも変更してください。
4. ログを確認したい場合は以下を利用します。
   ```bash
   docker compose logs app
   docker compose logs db
   ```
5. MySQL初期化は `db/init.sql` がコンテナ起動時に自動実行され、DB・ユーザー・権限が作成されます。

### ローカル実行（開発者向け）
1. `.env` を読み込めるようローカル環境に配置し、MySQLを手動で起動しておきます。
2. 依存パッケージを取得します。
   ```bash
   go mod tidy
   ```
3. アプリケーションを起動します。
   ```bash
   go run main.go
   ```
4. ブラウザで `http://localhost:8080`（`APP_PORT` を変更している場合はそのポート）へアクセスします。

## 提供ポート
- アプリケーション: `APP_PORT`（デフォルト `8080`） → `http://localhost:8080`
- MySQL: `MYSQL_HOST_PORT`（デフォルト `33306`） → `mysql -h 127.0.0.1 -P 33306`

## 主なエンドポイント
### Webページ
- `/menu` : ステップ形式のメニュー診断ページ
- `/order` : 注文内容の確認ページ
- `/terms` : 利用規約ページ
- `/health` : ヘルスチェック（200を返却）

### API
- `POST /api/orders` : 外部注文APIへリクエストを転送。`STORE_ID` が必要です。
- `GET /api/orders/id` : （現在はモック）注文状況確認。将来的に外部APIへ委譲予定。


## 開発用コマンド
- Goコード整形: `go fmt ./...`
- フロントエンド資産の整形: `npx prettier --write .`
- ヘルスチェック: `curl http://localhost:8080/health`

## トラブルシュート
- アプリが起動しない場合: `docker compose logs app` を確認。
- DB接続に失敗する場合: `.env` の `Server` / `Port` / `User` / `Password` / `Database` の値を再確認。
- 外部注文APIからエラーが返る場合: `STORE_ID` が正しいか、タイムアウトが発生していないかを確認。
