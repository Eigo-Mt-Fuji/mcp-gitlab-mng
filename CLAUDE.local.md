# プロジェクト基本情報

このプロジェクトは Go で書かれた Web アプリケーションです。

# 共通コマンド

- `go build`: プロジェクトのビルド実行
- `go test ./...`: テストスイート実行
- `gofmt -w .`: コードフォーマット適用
- `go mod tidy`: 依存関係の整理

# コードスタイル

- パッケージ名は小文字のみを使用
- 関数名は PascalCase（エクスポートする場合）または camelCase（内部使用）
- 定数名は 大文字のみを使用　
- インデント: タブ
- 最大行長: 100文字
- インターフェース名は er で終わることが推奨（例：Reader, Writer）
- エラー処理は明示的に行い、必ずエラーをチェック

# コミットメッセージ　

## フォーマット

- コミットする際、メッセージのフォーマットは次の形式に従って下さい( [参照](https://github.com/conventional-changelog/commitlint/tree/master/@commitlint/config-conventional) )

```plaintext
<type>[optional scope]: <subject>

[optional body]

[optional footer(s)]
```

- `type`, `subject` は必須です。`body`, `footer` を入れる場合はそれぞれ空行を間に挿みます。

## typeについて

`type` は lowerCase で表記し、次のいずれかの指定が必須です。  

| name | description |
| --- | ---|
| build | ビルド |
| ci | CI |
| chore | 雑事（カテゴライズする必要ないようなもの） |
| docs | ドキュメント |
| feat | 新機能 |
| fix | バグフィックス |
| perf | パフォーマンス |
| refactor | リファクタリング |
| revert | コミット取り消し（git revert） |
| style | コードスタイル修正 |
| test | テスト |

- 具体例

```
: some message     # fails
foo: some message  # fails
FIX: some message  # fails
fix: some message  # passes
```

### scope
`type` には追加のコンテキスト情報を表す `scope` を含むことが出来ます。  
`scope` は `type` の後に括弧付きで表し、lowerCase で表記します。

```
fix(SCOPE): some message  # fails
fix(scope): some message  # passes
```

### subject
`subject` は必須であり、命令形・現在系の動詞から始めます（例：'changed' や 'changes' ではなく 'change' から始まります）。  
コミットメッセージは「何をしたか」を記録するというよりも、「このコミットを適用するとどうなるか」を示す方が望ましいためです。

`subject` は lowerCase で表記します。

```
fix:               # fails
fix: Some Message  # fails
fix: SomeMessage   # fails
fix: SOMEMESSAGE   # fails
fix: some message  # passes
```

`subject` の末尾を `.` で終わらせてはなりません。

```
fix: some message. # fails
fix: some message  # passes
```

### body
`subject` の詳細情報が必要な場合は `body` セクションに記述してください。

```
fix: correct minor typos in code

see the issue for details on the typos fixed

closes issue #12
```

### footer
`footer` には、Breaking Changes についての情報や、このコミットがクローズした GitHub の課題を参照する場所でもあります。

Breaking Changes は、最初に `BREAKING CHANGE:` という単語で始まり、スペースか改行で始まります。

破壊的な変更は全て `footer` の `BREAKING CHANGE` ブロックとして記載しなければなりません。  
`BREAKING CHANGE` ブロックには、変更の説明、変更理由、移行の注意事項などを記載します。

```
BREAKING CHANGE: isolate scope bindings definition has changed and
  the inject option for the directive controller injection was removed.
  
  To migrate the code follow the example below:
  
  Before:
  
  scope: {
    myAttr: 'attribute',
    myBind: 'bind',
    myExpression: 'expression',
    myEval: 'evaluate',
    myAccessor: 'accessor'
  }
  
  After:
  
  scope: {
    myAttr: '@',
    myBind: '@',
    myExpression: '&',
    // myEval - usually not useful, but in cases where the expression is assignable, you can use '='
    myAccessor: '=' // in directive's template change myAccessor() to myAccessor
  }
  
  The removed `inject` wasn't generaly useful for directives so there should be no code using it.
```


クローズしたGitHub Issueへの参照を追加する場合、次のように `Closes` キーワードを先頭にして記述してください。

```
Closes #234
```

複数のIssueへの参照を追加するには、カンマ区切りで記述します。

```
Closes #123, #245, #992
```

# ワークフロー

- `go vet` でコード静的解析を実行
- テスト実行時は `-v` フラグで詳細な出力を確認
- ベンチマークテストは `go test -bench=.` で実行
