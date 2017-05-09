# 쪼개기

Parses Korean vocabulary and formats items for easy importing from lists or dictionaries into [Quizlet](https://quizlet.com/), [Anki](http://ankisrs.net/), and other flashcard apps.

## Installation

First install [Go](https://golang.org/doc/install) and then install `jjogaegi`:

    $ go get github.com/ryanbrainard/jjogaegi

## Usage

`jjogaegi` is a small, sharp UNIX-like tool. As such, it only reads from stdin and writes to stdout, so it can be used in a pipeline. Use a command like `cat` or `pbpaste` to input a file or the clipboard, and then use redirection (i.e. `>`) or `pbpaste` to output back to a file or the clipboard.

For example, if there is a file named `input.txt` that looks like this:

```
•컴퓨터를 켜다 to turn on the computer
•브라우저를 열다 to open the web browser
•검색어를 입력하다 to type in the search word
```

It can be processed and written to `output.tsv` like this:

```sh
$ cat input.txt | jjogaegi > output.tsv
```

The resulting file will be in tab-separated format:

```tsv
컴퓨터를 켜다	to turn on the computer
브라우저를 열다	to open the web browser
검색어를 입력하다	to type in the search word
```

The output can then be imported into your favorite flashcard app.

Alternatively, if you'd rather work with just data from/to the clipboard, use `pbpaste` and `pbcopy` on MacOS:

```sh
$ pbpaste | jjogaegi | pbcopy
```

The examples above shows processing without any options. This processes a simple list and outputs a TSV by default; however, the parser for the input and formatter for the output can be customized. See options below for details.

For example, convert a [한국어기초사전](https://krdict.korean.go.kr) Korean wordbook XML file for import into [Anki](http://ankisrs.net/) complete with tags, English definitions, and audio files:

```sh
$ cat input.xml | jjogaegi --parser krdict-xml --header 'tags: example' --audiodir /path/to/anki/media' > ouput.tsv
```

# Options

The parser and formatter can be set for different inputs and outputs.

## Parsers

Options for `--parser` flag:

 - `list`: (default) list of Korean terms followed by English definitions. Splits the line after the last 한글 character. Does not support 漢字.
 - `krdict-xml`: [한국어기초사전](https://krdict.korean.go.kr) Korean wordbook XML. Supports Id, 漢字, pronunciations, antonyms, and English definition fetching.
 - `naver-json`: [Naver Korean-English Dictionary](http://endic.naver.com/) wordbook JSON. Supports Id, 漢字.
 - `naver-table`: [Naver Korean-English Dictionary](http://endic.naver.com/) wordbook printed PDF table. Supports 漢字. (not recommended)

## Formatters

Options for `--formatter` flag:

 - `tsv`: (default) tab-separated values
 - `csv`: comma-separated values
 
## Hanja

Options for `--hanja` flag:

 - `none`: (default) do not include 漢字 with 한글
 - `parens`: output 漢字 in parenthesis next to 한글

## Header

Options for `--header` flag:

If set, string will be prepended to output

## Audio Dir

Options for `--audiodir` flag:

If set, audio will be downloaded into the directory specified

# Disclaimer

For personal use only. Do not use this tool for publishing copyrighted content. Respect copyright holders' rights of their content.
