# Tamo CLIアプリケーション

Tamoは、JSONを使用してデータを永続化するタスクとメモ管理のためのコマンドラインインターフェイス（CLI）アプリケーションです。開発者やAIエージェントがタスク、チェックリスト、および関連情報を管理するための、シンプルかつパワフルなツールとして設計されています。

## 概要

Tamoでできること：
- 順序と完了状態を持つタスクの管理
- タスクに関連するメモの作成と関連付け
- JSONフォーマットでのデータ保存
- タスクとメモ作成のためのMarkdown風構文の解析
- AIプロンプト用のメモ参照を展開したタスクのフラット化

## インストール

### ソースからビルド

```bash
# リポジトリをクローン
git clone https://github.com/zishida/tamo.git
cd tamo

# アプリケーションをビルド
go build -o tamo ./cmd/tamo

# バイナリをPATHの通ったディレクトリに移動（任意）
sudo mv tamo /usr/local/bin/
```

### Go Installを使用

```bash
go install github.com/zishida/tamo/cmd/tamo@latest
```

### GitHub Releasesから

複数のプラットフォーム向けのビルド済みバイナリが[GitHub Releases](https://github.com/zishida/tamo/releases)ページで入手できます。これらのバイナリは、新しいタグがリポジトリにプッシュされるたびにGitHub Actionsによって自動的にビルドおよび公開されます。

対応プラットフォーム：
- Linux AMD64
- Linux ARM64
- Darwin ARM64（macOS Apple Silicon）

ダウンロードとインストール方法：
1. [Releasesページ](https://github.com/zishida/tamo/releases)にアクセス
2. お使いのプラットフォームに適したバイナリをダウンロード
3. 実行可能にする：`chmod +x tamo-<platform>`
4. PATHの通ったディレクトリに移動：`sudo mv tamo-<platform> /usr/local/bin/tamo`

## 継続的インテグレーションとデプロイメント

Tamoは自動テストとリリースのためにGitHub Actionsを使用しています：

### 自動テスト

テストは次のタイミングで自動的に実行されます：
- mainブランチへの全てのプッシュ
- すべてのプルリクエスト

これによりコード品質が確保され、リグレッションを防止します。

### 自動リリース

リポジトリに新しいタグがプッシュされると、GitHub Actionsは自動的に：
1. 複数のプラットフォーム向けにバイナリをビルド（Linux AMD64、Linux ARM64、Darwin ARM64）
2. タグ名で新しいGitHubリリースを作成
3. ビルドしたバイナリをリリースのアセットとして添付

新しいリリースの作成方法：
```bash
# タグを作成してプッシュ
git tag v1.0.0
git push origin v1.0.0
```

## はじめに

### Tamoの初期化

Tamoを使用する前に、プロジェクトディレクトリで初期化する必要があります：

```bash
tamo init
```

これにより、カレントディレクトリに`.tamo`ディレクトリが作成され、空の`data.json`ファイルが生成されます。

## 基本的な使い方

### タスク管理

```bash
# リストの最後にタスクを追加
tamo add task "ドキュメントを完成させる" -d "プロジェクトの包括的なドキュメントを書く"

# リストの先頭にタスクを追加
tamo unshift task "優先度の高いタスク" -d "これを最初に行う必要がある"

# すべてのタスクを一覧表示
tamo list tasks

# タスクを完了としてマーク
tamo done <task_id>

# タスクの詳細を表示
tamo show <task_id>
```

### メモ管理

```bash
# メモを追加
tamo add memo "重要な情報" -c "これは覚えておくべき重要な内容です"

# すべてのメモを一覧表示
tamo list memos

# メモの詳細を表示
tamo show <memo_id>
```

### タスクワークフロー

```bash
# 次のタスク（最初の未完了タスク）を表示
tamo next

# 最初のタスクを完了としてマーク
tamo shift task --done

# 最後のタスクを削除
tamo pop task --rm
```

## 高度な機能

### Markdownの解析

Tamoは、Markdownファイルからタスクとメモを作成できます：

```bash
# Markdownファイルからタスクを作成
tamo add task -f task_description.md

# 標準入力からタスクを作成
cat task_description.md | tamo add task --from-stdin
```

Markdownファイルは次の形式に従う必要があります：

```markdown
# タスクのタイトル

タスクの説明をここに記述します。

```memo
これはタスクにリンクされる別のメモになります。
```

タスクの説明の続き。

```memo
追加情報を含む別のメモ。
```
```

### AIプロンプト用のタスクフラット化

タスクとそれに関連するすべてのメモ参照を展開してフラット化することができます：

```bash
tamo flattask <task_id>
```

これはAIツール用の包括的なプロンプトを作成したり、タスクとそれに関連するすべての情報の完全なビューを取得するのに役立ちます。

## プロジェクト構造

```
tamo/
├── cmd/
│   └── tamo/
│       └── main.go         # アプリケーションのエントリーポイント
├── internal/
│   ├── cli/
│   │   ├── cli.go          # CLIコマンド処理
│   │   └── markdown_parser.go # Markdown解析ロジック
│   ├── model/
│   │   └── model.go        # データモデル（Task、Memo、Store）
│   ├── storage/
│   │   └── storage.go      # JSON永続化
│   └── utils/
│       └── utils.go        # ユーティリティ関数
├── go.mod                  # Goモジュールファイル
└── README.md               # このファイル
```

## データモデル

- **Task**: ID、タイトル、説明、順序、完了状態、メモ参照などのプロパティを持つ、実行すべき作業を表します
- **Memo**: ID、タイトル、内容などのプロパティを持つ、タスクに関連する情報を保存します
- **Store**: すべてのタスクとメモを含むメインデータ構造

## ライセンス

[MITライセンス](LICENSE)