# cell

![logo](/docs/logo.png)

cellはコマンドライン向けのExcelファイル(.xlsx)処理用プログラミング言語です。

例えば、以下のようにExcelファイルを作成することができます。

```
$ cell -to greeting.xlsx '["A1"]="Hello, Excel"'
```

あるいは、以下のようにExcelファイルのデータをコマンドのパイプラインへ流し込む事もできます。

```
$ cell -from users.xlsx -n 'puts(["A".NER])' | grep twinbird
```

詳しくは[チュートリアル](/docs/tutorial_ja.md)を見てください。

オプションに関してはコマンドのヘルプを見てください。

```
$ cell -h
```

## インストール

[releaseページ](https://github.com/twinbird/cell/releases)からダウンロードして利用することができます。

## スペシャルサンクス

[Excelizeプロジェクト](https://github.com/360EntSecGroup-Skylar/excelize)に感謝します。

Excelizeなしには、このプロジェクトが作られることはありませんでした。
