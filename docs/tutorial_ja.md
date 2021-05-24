# プログラミング言語cellチュートリアル

cellはExcelファイル(xlsx形式)を読み書きするためのコマンドと言語です。

awk、perlなどから影響を受けています。

## インストール

[GitHub](https://github.com/twinbird/cell/releases)からダウンロードしてください。

## Hello, world

以下のコマンドでA1セルへ"Hello, world"と設定されたgreeting.xlsxが作成されます。

```
$ cell -to greeting.xlsx '["A1"]="Hello, world";'
```

あるいは標準出力を経由して挨拶したい場合には以下を実行します。

```
$ cell 'puts("Hello, world");' # => "Hello, world"
```

cellは文末にセミコロンか改行を必要とします。
ただし、プログラムの末尾には改行が自動的に挿入されるので、上記のコマンドは以下と同等です。

```
$ cell 'puts("Hello, world")'  # => Hello, world
```

## 入力を得よう

Excelファイルから入力したい場合にはfromオプションが使えます。

以下はusers.xlsxのA1セルの内容をコンソールへ表示します。

```
$ cell -from users.xlsx 'puts(["A1"])'
```

標準入力の内容を取得したい場合にはgets()が使えます。

```
$ echo "Hello, world" | cell -to greeting.xlsx '["A1"]=gets()'
```

## 値・変数・式

cellのデータ型は文字列と64bit浮動小数点数のみです。

変数は宣言の必要はなく、プログラム内で現れた時点で空文字列で初期化されて用意されます。

変数のスコープはグローバルと関数ごとの2種類のみです。

配列などのデータ構造はありません。

代わりにExcelの構造を利用できます。

### Excelデータへのアクセス

["A1"]のように文字列を\[と\]で囲むと、開いているExcelブックのアクティブシートのA1セルへアクセスできます。

アクティブシートは@特殊変数で確認・設定することができます。

```
$ cell 'puts(@)' #=> "Sheet1"
```

```
$ cell '@="Sheet2";puts(@)' #=> "Sheet2"
```

@特殊変数へ設定した名前のシートがアクティブシートに設定されます。

作成されていないシート名を設定した場合にはそのシートが作成されます。

### 式

cellは演算子によって値の解釈を変えます。

```
one = 1
two = 2
puts(one + two) # => 3(add number)
puts(one . two)  # => 12(string concat)
```

ほとんどの演算子はどれも他の言語でよくあるものです。
演算子の詳細は[演算子一覧](#演算子)をご覧下さい。

ただしインクリメント演算子はかなり風変わりに感じるかもしれません。

これらの演算子は変数の値を数値として解釈し、一つ次の値にしますが、Excelの列番号として解釈できる場合はその次の文字列を設定します。

```
cell 'v=1;v++;puts(v);' # => 2
cell 'v="a";v++;puts(v);' #=> B
cell 'v="Z";v++;puts(v);' #=> AA
cell 'v="string";v++;puts(v);' #=> 1
```

## 分岐を使う

cellでは分岐にif文を使います。
他の言語でよく採用されているのと同様、elseを使うこともできます。

```
$ cell 'if (gets()){ puts("true"); }else{ puts("false");}'
```

Cを祖先とする言語によくある通り、分岐の内容が単文の場合にはブロックの{}は不要です。

```
$ cell 'if (gets()) puts("true"); else puts("false");'
```

if文は()内の値が空文字列か数値の0の場合は偽に、その他は真となります。

## ループ

cellにはwhileとdo-while、forの3種類のループ構造があります。

if文と同様にブロック本体の内容が単文の場合にはブロックの{}は不要です。

### while

```
$ cell 'i=0;while(i<3){puts(i);i++;}'
# => 0
# => 1
# => 2
```

### do-while

```
$ cell 'do { puts("output this text"); } while(0);'
# => output this text
```

### for

```
$ cell 'for(i=0; i<3;i++) { puts(i); }'
# => 0
# => 1
# => 2
```

## 関数の定義

関数は呼び出す前に定義されていなければなりません。

関数内では変数のスコープは別のものになります。

```
# return x + y
function add(x, y) {
  return x + y;
}
puts(add(1, 2));
# => 3
```

## コメント

\#から行末まではコメントです。

## 実践的な例

百聞は一見に如かずといいますよね。

### ユーザの一覧をExcel表にする

システムのユーザとホームディレクトリを一覧にします。

```
cat /etc/passwd | cell -n -F ":" -to users.xlsx '["A".NR]=$1;["B".NR]=$6'
```

### 商品情報の入ったExcelをpsql経由でDBへ登録します

1行目にヘッダがあるExcelの表からSQLを作ります。

```
cell -from items.xlsx -N -s 2 'puts("INSERT INTO items(name, value) (" . ["A".NER] . "," . ["B".NER] .");")' | psql mydb
```

### 見積書が複数入ったブックから見積の件名の一覧表を作ります

日本においてはExcelで帳票を作成するケースが非常に多いです。

見積の件名がC4セルに入っていて、1シート1見積書となっているExcelブックから見積一覧のExcelファイルを作りましょう。

```
 cell -from estimates.xlsx 'head(); for(i=1;i<=count();i++) {puts(["C4"]);@++;}' | cell -F "\t" -to estimate_list.xlsx -n '["A".NR]=$1'
```

### 部署のメンバーの名刺を作る

日本の会社ではExcelを使って名刺を作ることもあります。

ひな形のシート(template)を含むExcelブックtemplate.xlsxを使って名前を書くシートのD5セルへ入れていきましょう。

```
cat member.txt | cell -from template.xlsx -to business_cards.xlsx -n 'copy("template", $1);@=$1;["D5"]=$1;'
```

## 簡易リファレンス

### 演算子

#### 代入/参照演算子

| 演算子 | 意味 |
| --------|------|
| \= | 変数へ代入 |
| \[文字列\] | セルの値を参照します |
| \[文字列\] \= | セルへ値を設定します |

#### 数値演算子

| 演算子 | 意味 |
| --------|------|
| + | 数値として解釈し、加算 | 
| - | 数値として解釈し、減算 |
| * | 数値として解釈し、乗算 | 
| / | 数値として解釈し、除算 |
| % | 数値として解釈し、剰余 |
| ** | 数値として解釈し、べき乗 |
| +(単項) | 文字列を数値として解釈 | 
| -(単項) | 文字列を数値として解釈し、符号を反転 |
| += | 加算して代入 |
| -= | 減算して代入 |
| /= | 除算して代入 |
| *= | 乗算して代入 |
| %= | 剰余を代入 |
| **= | べき乗して代入 |

#### 文字列演算子

| 演算子 | 意味 |
| --------|------|
| . | 文字列を結合 |

#### 数値比較演算子

| 演算子 | 意味 |
| --------|------|
| < | 数値として解釈し、比較(より小さい) |
| > | 数値として解釈し、比較(より大きい) |
| <= | 数値として解釈し、比較(以下) |
| >= | 数値として解釈し、比較(以上) |
| == | 数値として解釈し、比較(等しい) |
| != | 数値として解釈し、比較(等しくない) |

#### 文字列比較演算子

| 演算子 | 意味 |
| --------|------|
| eq | 文字列として解釈し、等しい |
| ne | 文字列として解釈し、等しくない |
| ~ | 正規表現文字列にマッチしている |
| !~ | 正規表現文字列にマッチしていない |

#### セル番号比較演算子

| 演算子 | 意味 |
| --------|------|
| lt | Excelの列番号として解釈し、より小さい |
| le | Excelの列番号として解釈し、以下 |
| gt | Excelの列番号として解釈し、より大きい |
| ge | Excelの列番号として解釈し、以上 |

#### 論理演算子

| 演算子 | 意味 |
| --------|------|
| && | 論理積 |
| \|\| | 論理和 |
| ! | 論理否定 |

#### インクリメント/デクリメント演算子

インクリメント/デクリメント演算子は変数に対して適用する演算子です。

前置と後置があります。

これらの演算子は変数や変数の値によって動作が異なります。

 * 数値の場合にはインクリメント/デクリメントを行います
 * 文字列の場合には数値として解釈し、インクリメント/デクリメントを行います
 * @変数の場合には次のシート/前のシートへ変更します
 * 列番号文字列の場合には次の列番号/前の列番号へ変更します

### 特殊変数

| 変数 | 意味 |
| -----|-----|
| @ | アクティブシート名 |
| FS | 標準入力のフィールドセパレータ(デフォルトは単一のスペースかタブ) |
| OFS | 標準出力のフィールドセパレータ(デフォルトは単一のスペース) |
| RS | 標準入力のレコードセパレータ(デフォルトは改行) |
| ORS | 標準出力のレコードセパレータ(デフォルトは改行) |
| NR | 標準入力から取り込んだ行数 |
| NER | ループ処理でエクセルの何行目を示しているか(オプションNで利用) |
| SER | エクセルの何行目から処理を行うか(オプションNで利用) |
| LR | アクティブシートの最終行番号 |
| LC | アクティブシートの最終列番号 |
| LCC | アクティブシートの最終列番号(アルファベット名) |
| $0 | gets()で得た直前の標準入力 |
| $1 | gets()で得た標準入力をフィールドセパレータで分割したものの1フィールド目 |
| $n | gets()で得た標準入力をフィールドセパレータで分割したもののnフィールド目 |
| $_0 | ~(マッチ演算子)でマッチした文字列 |
| $_1 | ~(マッチ演算子)でマッチした際にキャプチャした1つめの文字列。キャプチャは()で行います。 |
| $_n | ~(マッチ演算子)でマッチした際にキャプチャしたn番めの文字列 |

### コマンドラインオプション

| オプション | 意味 |
| --------------|------|
| -to | 処理結果のExcelファイルを保存するパスを指定します |
| -from | 処理のために読み込むExcelファイルパスを指定します |
| -f | cellプログラムの書かれたファイルを指定します。このオプションが指定された場合は第一引数のプログラムは実行されません。 |
| -F | フィールドセパレータ(FS変数)を指定します |
| -n | 実行するプログラム全体を[while(gets()){... ;}]で囲みます |
| -N | 実行するプログラム全体を[for(NER = SER; NER <= LR; NER++){... ;}]で囲みます |
| -s | SER変数の値を設定します |
| -S | @変数の値を設定します |
| -V | バージョン情報を表示します |
| -h | ヘルプを表示します |


### 組み込み関数

#### exit(n)

exitは終了コードnでプログラムを終了します。
toオプションが指定されている場合には終了時に処理中のExcelファイルを出力します。

#### abort(n)

abortは終了コードnでプログラムを終了します。
exitと異なり、toオプションが指定されていても終了時にExcelファイルを出力しません。

#### gets()

getsは標準入力から改行までの1行分の文字列を読み取って返します。

入力末尾に達した場合には空文字列を返します。

getsは特殊変数$0, $1...$nに値を設定します。

$0には読み取った文字列全体、$1はFSで区切られた１つ目のフィールドが入ります。$2以降も同様です。

また、読み取るたびに特殊変数NRをインクリメントします。

#### puts(s...)

putsは文字列sを標準出力へ出力します。

引数なしで呼び出す場合には$0の内容を出力します。

複数の引数を指定する場合にはOFSで結合した文字列を出力します。

出力する文字列の末尾にはORSが付与されます。

#### head()

特殊変数@へ現在開いているExcelブックの先頭のシートを設定します。

#### tail()

特殊変数@へ現在開いているExcelブックの末尾のシートを設定します。

#### rename(old, new)

シート名oldをnewに変更します。

#### exist(sheetname)

現在開いているExcelブックにシート名のシートがある場合には1を、ない場合には0を返します。

#### count()

現在開いているExcelブックのシート数を返します。

#### delete(sheetname)

sheetnameシートを削除します。

#### copy(from, to)

fromシートをtoという名前でコピーします。

#### srand(n)

疑似乱数の種を設定します。

乱数生成器の新しい種としてnを使用します。 nが指定されない場合は、現在の時刻を使用します。

#### rand()

0以上1以下の疑似乱数を返します。

#### floor(n)

nの小数点以下を切り捨てにした値を返します。

#### ceil(n)

nの小数点以下を切り上げにした値を返します。

#### round(n)

nの小数点以下で四捨五入した値を返します。