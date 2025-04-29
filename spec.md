## tamo CLI アプリケーション実装仕様書 (v1.0)

### 1. 概要

`tamo` は、開発者やAI Agentが使用することを想定した、チェックリスト、TODOリスト、メモ管理を統合したコマンドラインインターフェース (CLI) アプリケーションである。タスク（順序付き、完了状態あり）とメモ（非構造化テキスト）を管理し、JSONファイルにデータを永続化する。Markdownライクな記法でのタスクとメモの同時入力や、タスクに紐づくメモ情報を展開してフラットなテキストを生成する機能を特徴とする。

言語はgoで実装する。

### 2. コアコンセプト

*   **Task:** 実行すべき作業や指示を表す。以下の属性を持つ:
    *   一意なID (UUID推奨)
    *   タイトル (簡潔な名前)
    *   詳細な説明 (Markdown形式を許容)
    *   実行順序 (`order`): 浮動小数点数。小さい値が先。
    *   完了フラグ (`done`): boolean (true/false)
    *   参照メモIDリスト (`memo_refs`): 関連するメモのIDを配列で保持。
    *   作成日時 (`created_at`)
    *   更新日時 (`updated_at`)
*   **Memo:** タスク実行に必要な情報や、単なる覚え書きを保存する。以下の属性を持つ:
    *   一意なID (UUID推奨)
    *   タイトル (省略可能)
    *   内容 (`content`): 複数行のテキスト。
    *   作成日時 (`created_at`)
    *   更新日時 (`updated_at`)
*   **ID:** 各Task/Memoを一意に識別するための文字列 (UUID v4推奨)。
*   **順序 (Order):** Taskの実行順を示す。未完了タスクリストはこの `order` でソートされる。

### 3. データストレージ

*   **形式:** JSON
*   **デフォルトパス:** コマンド実行時のカレントディレクトリ配下の `.tamo/data.json`
    *   `tamo init` コマンドで `.tamo` ディレクトリと空の `data.json` を生成する。
*   **エンコーディング:** UTF-8

### 4. データ構造 (JSON)

```json
{
  "version": 1,
  "tasks": [
    {
      "id": "string (uuid)",
      "title": "string",
      "description": "string (multiline, markdown possible)",
      "order": "float",
      "done": "boolean",
      "memo_refs": ["string (uuid)", ...],
      "created_at": "string (ISO 8601)",
      "updated_at": "string (ISO 8601)"
    },
    // ... more tasks
  ],
  "memos": [
    {
      "id": "string (uuid)",
      "title": "string | null", // 省略可能
      "content": "string (multiline)",
      "created_at": "string (ISO 8601)",
      "updated_at": "string (ISO 8601)"
    },
    // ... more memos
  ]
}
```

### 5. コマンドラインインターフェース (CLI)

#### 5.1. 初期化

*   **コマンド:** `tamo init`
*   **動作:**
    *   カレントディレクトリに `.tamo` ディレクトリが存在しない場合は作成する。
    *   `.tamo` ディレクトリ内に、上記のデータ構造を持つ空の `data.json` を作成する。ファイルが既に存在する場合は上書きしない（エラーまたは警告メッセージを表示）。

#### 5.2. 追加 (Add)

*   **メモ追加:** `tamo add memo [<title>] [-c "<content>" | --from-stdin | --editor]`
    *   `[<title>]`: メモのタイトル（任意）。
    *   `-c "<content>"`: 内容をコマンドライン引数で指定。
    *   `--from-stdin`: 標準入力から内容を読み込む。
    *   `--editor`: 環境変数 `EDITOR` で指定されたエディタを起動し、内容を入力させる（デフォルト動作でも良い）。
    *   新しいMemoオブジェクトを生成し、`memos` 配列に追加。UUIDを自動生成。
*   **タスク追加 (Markdown記法):** `tamo add task -f <filepath> | --from-stdin`
    *   `-f <filepath>`: 指定されたファイルからMarkdownテキストを読み込む。
    *   `--from-stdin`: 標準入力からMarkdownテキストを読み込む。
    *   **解析ロジック:**
        1.  最初のH1見出し (`# title`) をTaskの `title` とする。見出しがない場合はエラーまたはファイル名を仮タイトルとする。
        2.  ` ```memo ... ``` ` ブロックを検出する。
        3.  各 ` ```memo ... ``` ` ブロックについて:
            *   ブロック内の内容で新しい `Memo` オブジェクトを生成（UUID付与）。
            *   `memos` 配列に追加。
        4.  元のテキストからH1行と ` ```memo ... ``` ` ブロックを除去する。
        5.  元の ` ```memo ... ``` ` があった箇所を、対応するMemoへの参照 (`[memo](<memo_id>)` 形式) に置き換えたものを `description` とする。
        6.  生成された `title`, `description`, Memo IDリスト (`memo_refs`) で新しい `Task` オブジェクトを生成（UUID付与）。`order` は既存タスクの最大値+1.0とする。`done`は`false`。
        7.  `tasks` 配列に追加。
*   **タスク追加 (標準):** `tamo add task "<title>" [-d "<description>"] [-m <memo_id>,...]`
    *   シンプルなTaskを追加する。これは `tamo push task ...` のエイリアスとして実装しても良い。
    *   `-d "<description>"`: 詳細説明。
    *   `-m <memo_id>,...`: 参照するメモID (カンマ区切り)。
    *   `order` は既存タスクの最大値+1.0、`done` は `false`。
*   **タスク追加 (末尾):** `tamo push task "<title>" [-d "<description>"] [-m <memo_id>,...]`
    *   未完了タスクリストの **末尾** にタスクを追加する。
    *   内部的に、既存の全タスクの `order` の最大値 + 1.0 を新しい `order` として設定する。
*   **タスク追加 (先頭):** `tamo unshift task "<title>" [-d "<description>"] [-m <memo_id>,...]`
    *   未完了タスクリストの **先頭** にタスクを追加する。
    *   内部的に、既存の全タスクの `order` の最小値 - 1.0 (または最小値が1以上なら最小値/2.0など、衝突しない値) を新しい `order` として設定する。

#### 5.3. 表示 (List/Show)

*   **一覧表示:** `tamo list [tasks|memos|all] [--done|--undone] [--refs <memo_id>]`
    *   `tasks` (デフォルト): 未完了タスク (`done: false`) を `order` 昇順で表示。
    *   `memos`: 全てのメモを表示。
    *   `all`: 全てのタスクとメモを表示。
    *   `--done`: 完了済みタスクのみ表示。
    *   `--undone`: 未完了タスクのみ表示。
    *   `--refs <memo_id>`: 指定したメモを参照しているタスクを表示。
    *   **出力形式 (Task):** ID(短縮形でも可), Order, Done状態([x] or [ ]), Title
    *   **出力形式 (Memo):** ID(短縮形でも可), Title (あれば), Content (最初の1行程度)
*   **詳細表示:** `tamo show <id>`
    *   指定したIDのTaskまたはMemoの詳細情報を表示する。
    *   Taskの場合: Title, ID, Order, Done状態, Description (Markdownレンダリングされると望ましい), 参照しているMemoのIDとタイトル(あれば)。
    *   Memoの場合: Title(あれば), ID, Content。

#### 5.4. 編集 (Edit/Modify)

*   **編集:** `tamo edit <id> [--editor]`
    *   指定したIDのTaskまたはMemoを編集する。
    *   環境変数 `EDITOR` で指定されたエディタを起動し、現在の内容を一時ファイルに書き出して編集させる。保存後に `data.json` を更新する。
    *   Task編集対象: `title`, `description`, `memo_refs` (IDリストを編集できるようにする)。
    *   Memo編集対象: `title`, `content`。
*   **完了化:** `tamo done <task_id>`
    *   指定したTaskの `done` フラグを `true` に設定する。
*   **未完了化:** `tamo undone <task_id>`
    *   指定したTaskの `done` フラグを `false` に設定する。
*   **移動 (絶対):** `tamo mv <task_id> <target_order>`
    *   指定したTaskの `order` を `<target_order>` (float) に直接設定する。
*   **移動 (相対):** `tamo mv <task_id> before|after <other_task_id>`
    *   `before`: `<other_task_id>` の `order` よりわずかに小さい値に設定する (例: `(prev_order + other_order) / 2.0`)。
    *   `after`: `<other_task_id>` の `order` よりわずかに大きい値に設定する (例: `(other_order + next_order) / 2.0`)。

#### 5.5. 削除 (Remove)

*   **削除:** `tamo rm <id> [-f|--force]`
    *   指定したIDのTaskまたはMemoを `data.json` から削除する。
    *   Memoを削除する際、そのMemoを参照しているTaskが存在する場合は警告を表示し、`-f` または `--force` オプションがない限り削除を中止する。

#### 5.6. 実行補助 / ワークフロー

*   **末尾操作:** `tamo pop task [--done | --rm]`
    *   未完了タスクリストの **末尾** (最大の`order`を持つ) のタスクを対象とする。
    *   オプションなし: 末尾タスクの詳細を `tamo show` 相当で表示。
    *   `--done`: 末尾タスクの `done` を `true` にし、そのタスク情報を表示。
    *   `--rm`: 末尾タスクを削除し、そのタスク情報を表示（要確認プロンプト or `-f`）。
    *   未完了タスクがない場合はエラーメッセージを表示。
*   **先頭操作:** `tamo shift task [--done | --rm]`
    *   未完了タスクリストの **先頭** (最小の`order`を持つ) のタスクを対象とする。
    *   オプションなし: 先頭タスクの詳細を `tamo show` 相当で表示。
    *   `--done`: 先頭タスクの `done` を `true` にし、そのタスク情報を表示。
    *   `--rm`: 先頭タスクを削除し、そのタスク情報を表示（要確認プロンプト or `-f`）。
    *   未完了タスクがない場合はエラーメッセージを表示。
*   **次タスク表示:** `tamo next`
    *   `tamo shift task` (オプションなし) のエイリアス。最も `order` が小さい未完了タスクを表示する。
*   **タスク平坦化:** `tamo flattask <task_id>` (旧: `tamo prompt`, `tamo expand`)
    *   指定した `task_id` の情報を元に、参照している全てのMemoの内容を展開（フラット化）して一つのテキストを生成し、標準出力に表示する。
    *   **出力形式:**
        ```markdown
        # <Task Title>

        <Task Descriptionのテキスト>
        (Description中の [memo](<memo_id>) 参照を...)

        <対応するMemoの内容を埋め込む>
        (例: ``` で囲んで表示)

        <Descriptionの続き...>
        ```
    *   AIへのプロンプト生成や、タスクの全体像把握に利用。

### 6. エラーハンドリング

*   `data.json` が見つからない、または不正な形式の場合、適切なエラーメッセージを表示する。
*   存在しないIDを指定された場合、エラーメッセージを表示する。
*   コマンドやオプションの使い方が間違っている場合、ヘルプメッセージやエラーメッセージを表示する。
*   リストが空の状態で `pop`/`shift` を実行した場合、エラーメッセージを表示する。
*   破壊的な操作 (rm, pop --rm, shift --rm) では、`-f` オプションがない限り確認を求めるか、警告を出す。

### 7. 実装上の考慮点

*   **言語:** Python, Go, Rust, Node.jsなどが候補。JSON操作、CLIライブラリ、UUID生成、日時扱いの容易さを考慮する。
*   **依存:** 外部ライブラリは最小限に留める。
*   **ID生成:** 標準ライブラリのUUID v4を使用する。
*   **ファイルI/O:** JSONの読み書きはアトミックに行うことが望ましい（一時ファイル経由での書き込みなど）。ファイルロックも検討。
*   **日時:** ISO 8601形式 (UTC推奨) で記録する。
*   **CLIフレームワーク:** `argparse` (Python), `cobra` (Go) などの利用を推奨。
*   **コード構造:** データ管理ロジックとCLI表示ロジックを分離し、将来的なライブラリ化/API化を容易にする。

