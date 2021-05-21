# cell

![logo](/docs/logo.png)

[日本語版](/docs/README_ja.md)

cell is a programming language for processing Excel files (.xlsx) for the command line.

For example, you can create an Excel file as follows.

```
$ cell -to greeting.xlsx '["A1"]="Hello, Excel"'
```

Alternatively, data from an Excel file can be fed into the command pipeline as shown below.

```
$ cell -from users.xlsx -n 'puts(["A".NER])' | grep twinbird
```

See the [tutorial](/docs/tutorial.md) for details.

Also, see the command help for options.

```
$ cell -h
```

## Install

You can download from the [release page](https://github.com/twinbird/cell/releases).

## Special Thanks

Thanks to the [Excelize project](https://github.com/360EntSecGroup-Skylar/excelize).

Without Excelize, this project would never have been created.

And many many thank you DeepL.
