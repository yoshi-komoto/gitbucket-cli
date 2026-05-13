# gb — GitBucket CLI

GitBucket の Pull Request 関連 API を呼ぶ CLI ツール。

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

PAT は GitBucket の Account Settings → Applications → "Personal access tokens" で発行する。

## Usage

```sh
gb pr list
gb pr view 12
gb pr comment list 12
gb pr comment add 12 --body "LGTM"
```
