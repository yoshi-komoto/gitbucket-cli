# gb — GitBucket CLI

GitBucket の Pull Request 関連 API を呼ぶ Go 製 CLI。GitHub の `gh` ライクな操作感。

## Install

```sh
go install github.com/yoshi-komoto/gitbucket-cli@latest
```

## Configure

`~/.config/gitbucket/config.yaml` を作成:

```yaml
url: https://gitbucket.example.com
token: <Personal Access Token>
```

PAT は GitBucket の **Account Settings → Applications → Personal access tokens** で発行する。

環境変数 `GITBUCKET_URL` と `GITBUCKET_TOKEN` を両方指定すれば config ファイルを無視してそちらを使う (CI 用)。

## Usage

```sh
# 一覧 (デフォルト state=open, limit=30)
gb pr list

# state や limit を変える
gb pr list --state all --limit 50

# 詳細表示
gb pr view 12

# コメント一覧 / 追加
gb pr comment list 12
gb pr comment add 12 --body "LGTM"
gb pr comment add 12 --body-file ./review.md
echo "looks good" | gb pr comment add 12

# JSON で出して jq に食わせる
gb pr list --output json | jq '.[].title'
```

`--repo OWNER/REPO` を付けなければ、カレントディレクトリの `git remote get-url origin` から自動推測する。

## Exit codes

| code | 状況 |
| --- | --- |
| 0 | 成功 |
| 1 | その他のエラー |
| 2 | 設定不足 / リポジトリ未解決 |
| 4 | 認証失敗 (401/403) |
| 5 | Not Found (404) |
| 6 | その他の API エラー |
| 64 | usage error |

## Development

```sh
go test ./... -race
go build ./...
```
