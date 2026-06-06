# パッケージ - metadata

記事のメタデータに関するリポジトリパッケージ。

## ファイル

| ファイル名 | 内容 |
|-----------|-----|
| [dynamodb.go](./dynamodb.go) | DynamoDB のリポジトリ |
| [domain.go](./domain.go) | リポジトリで使う型とバリデーション |
| [converter.go](./converter.go) | ドメイン同士の変換関数 |

## メタデータの定義

メタデータは、記事の内容以外の部分を指す。次のプロパティを持つ。

| プロパティ名 | 型 | 詳細 |
|------------|----|------|
| `slug`       | `string`   | 記事固有の ID |
| `title`      | `string`   | 記事のタイトル |
| `source`     | `source`   | 記事が公開されている場所 |
| `contentKey` | `string`   | 記事の内容の場所 |
| `status`     | `published \| unpublished` | 記事の公開状況 |
| `tags`       | `string[]` | 記事のタグ |
| `updatedAt`  | `time.Time` | 更新日時 |
| `createdAt`  | `time.Time` | 作成日時 |
| `pv`         | `int`       | ページ閲覧数 |

## 操作

以下の操作を想定する。

- 記事の一覧を所得する
  - カーソルベースのページネーション (固定)
  - 更新日時順で降順ソート (固定)
  - タグ絞り込みが可能 (任意)

## DB 設計

DynamoDB を使用する。以下の 2 種類のアイテムを 1 テーブルに保存する。

### 記事ベース

タグ絞り込みをしない場合にヒットするアイテム。

| プロパティ名 | 型 | 詳細 |
|-------------|----|------|
| `slug`        | S  | slug を保存する |
| `item_type`   | S  | `ARTICLE` 固定 |
| `status`      | S  | status を保存する |
| `filter_tag`  | S  | `#ALL` 固定 |
| `title`       | S  | title を保存する |
| `source`      | S  | source を保存する |
| `content_key` | S  | contentKey を保存する |
| `tags`        | SS | tags を保存する |
| `updated_at`  | S  | updatedAt を ISO8601 / RFC3339 で保存する |
| `created_at`  | S  | createdAt を ISO8601 / RFC3339 で保存する |
| `pv`          | N  | pv を保存する |

### タグ

タグ絞り込みをした場合にヒットするアイテム。ある記事に対して、ついているタグの個数だけ作成される。`slug` を含む `item_type` と `filter_tag` 以外のプロパティは、対応する記事ベースと同じ値を設定する。

| プロパティ名 | 型 | 詳細 |
|-------------|----|------|
| `slug`        | S  | slug を保存する |
| `item_type`   | S  | `TAG#{tag}` |
| `status`      | S  | status を保存する |
| `filter_tag`  | S  | tag を保存する |
| `title`       | S  | title を保存する |
| `source`      | S  | source を保存する |
| `content_key` | S  | contentKey を保存する |
| `tags`        | SS | tags を保存する |
| `updated_at`  | S  | updatedAt を ISO8601 / RFC3339 で保存する |
| `created_at`  | S  | createdAt を ISO8601 / RFC3339 で保存する |

### GSI

記事一覧取得時のための GSI を次のように設計する。PK は複合キーを指定する。

- **インデックス名**: `GSI_StatusTag_UpdatedAt`
- **PKs**
  - `status` (S)
  - `filter_tag` (S)
- **SK**: `updated_at` (S)
- **Projection Type**: Include
  - `title`
  - `tags`
  - `source`
  - `content_key`
  - `created_at`

## テスト

テストの実行には dynamodb-local の起動が必要。`go tool task test` で自動で起動する。
