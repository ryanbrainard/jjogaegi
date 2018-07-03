# jjogaegi (쪼개기)

Command line utility to create, parse, and format Korean vocabulary to import into [Anki](http://ankisrs.net/), [Quizlet](https://quizlet.com/), and other flashcard apps. 

 - Accepts parsing from various import sources, including plain-text lists, TSV (e.g. Anki notes), Memrise lists, various dictionary sources, and even an interactive mode to create flashcards on the fly. 
- Output formats include TSV, CSV, and JSON. 
- Interacts with the National Institute of Korean Basic Dictionary (국립국어원 한국어기초사전) to lookup words and enhance flashcards with definitions, example sentences, pronunciation audio, and more.

## Installation

This is a command-line application and requires basic knowledge of using command line. If you are unfamiliar with the command line, see a tutorial like [this one](https://lifehacker.com/5633909/who-needs-a-mouse-learn-to-use-the-command-line-for-almost-anything) for your platform. 

Go to the [release page](https://github.com/ryanbrainard/jjogaegi/releases) and download the latest release for your platform. Unzip the release and contained the `jjogaegi` binary can be placed anywhere on your computer.

Before use, it is highly recommended to set the following [environment variables](https://www.schrodinger.com/kb/1842): 

- `KRDICT_API_KEY`: Dictionary API key to enable word lookups
- `MEDIA_DIR`: Directory to download images and audio. For use with Anki, see its manual entry on [File Locations](https://apps.ankiweb.net/docs/manual.html#files).
		

## Usage

`jjogaegi` is a small, sharp UNIX-like tool. As such, by default it reads from `stdin` and writes to `stdout`, so it can be used in a pipeline; however, it can be configured with the `-in` and `-out` options to read and write to and from files as well. The input and output formats can be configured with the `-parser` and `-formatter` options, respectively. If no parser is configured, you are prompted to input words interactively and the definitions are looked up automatically. With the other parsers, the same lookup functionality can be enabled with the `-lookup` option, and interactive mode can be enabled with the `-interactive` option to choose one of multiple definitions when a word is a homophone. 

The column order for parsing and formatting TSV and CSV files is fixed as follows:

- **Note ID**: unique id 
- **External ID**: id from external system (e.g. dictionary id)
- **Hangul**
- **Hanja**
- **Korean Definition**
- **English Definition**
- **Pronunciation**
- **Audio**: file name of an audio file in `MEDIA_DIR`. If formatted as a URL, the file will automatically be downloaded and the file name reformatted.
- **Image**: file name of an image file in `MEDIA_DIR`. If formatted as a URL, the file will automatically be downloaded and the file name reformatted.
- **Grade**: difficulty of word
- **Antonym**
- **Example 1 Korean**
- **Example 1 English**
- **Example 2 Korean**
- **Example 2 English**

To configure Anki to support these columns, import the [sample Anki card](https://github.com/ryanbrainard/jjogaegi/raw/master/assets/anki-sample-card.apkg) to create the `Korean++` note type.

### Prompt Mode

If no input file is provided, you are prompted for words interactively. The definitions, example sentences, etc. are automatically looked up in the National Institute of Korean Basic Dictionary (국립국어원 한국어기초사전). If there are multiple definitions or one cannot be found, you will be prompted. For example, entering 안경 (one definition), 안녕 (two definitions), and 아이폰 (no definition):

```
$ jjogaegi -out /tmp/test.tsv -formatter tsv
Enter a Korean word on each line: (press Ctrl+D to quit)
>>> 안경
안경 -> glasses; spectacles

>>> 안녕
안녕 -> Multiple results found:
 1) hello; hi; good-bye; bye
 2) peace; good health
Enter number: 1


>>> 아이폰
아이폰 -> Not found. Enter custom English definition: iPhone

<Ctrl+D>
```

The output file (`/tmp/test.tsv`) looks then likes like this:

```tsv
608a6497-1915-4d90-9d65-6f48d6add48c	krdict:kor:31484:단어	안경	眼鏡	눈을 보호하거나 시력이 좋지 않은 사람이 잘 볼 수 있도록 눈에 쓰는 물건.	glasses; spectacles := An instrument that one wears over the eyes to proctect them or to supplement his/her eyesight for better vision.	안ː경	[sound:6db1c685-37b7-40a2-b7f2-62df68c422f8.mp3]	"<img src=""65c2ce63-b68c-455f-8cbb-189f4307b36e.jpg"">"	초급		검은 테 안경.		그는 피곤한지 안경을 벗고 눈을 비볐다.
4a9dc630-a51d-49df-b687-5c715151d376	krdict:kor:17298:단어	안녕	安寧	친구 또는 아랫사람과 서로 만나거나 헤어질 때 하는 인사말.	hello; hi; good-bye; bye := A salutation uttered when the speaker meets or parts from his/her friend or junior.	안녕	[sound:0d65ec77-fc92-4807-916a-b10722f80632.mp3]		초급
5f34ffa0-652e-4969-9197-7b092f9e808b	-	아니폰			iPhone
```

This output file can then be imported into your favorite flashcard app.

### File Based

If there is a file named `input.txt` that looks like this:

```
•컴퓨터를 켜다 to turn on the computer
•브라우저를 열다 to open the web browser
•검색어를 입력하다 to type in the search word
```

It can be processed and written to `output.tsv` like this (UNIX pipes and redirects can also be used, but not showing for brevity):

```sh
$ jjogaegi -in input.txt -out output.tsv -parser list -formatter tsv
```

The output file (`output.tsv`) looks then likes like this:

```tsv
e1dd2d1d-eb85-4805-b9c8-536dad3b11e9		컴퓨터를 켜다			to turn on the computer
167187f9-51cb-4c0d-b5b7-538915f096e3		브라우저를 열다			to open the web browser
a66ca59c-c413-4b79-b902-e3b32bb6bd08		검색어를 입력하다			to type in the search word
```

This output file can then be imported into your favorite flashcard app.

# Options

The parser and formatter can be set for different inputs and outputs.

## Parsers

Options for `--parser` flag:

 - `prompt`: interatively prompt user for words. 
 - `list`: list of Korean terms followed by English definitions. Splits the line after the last 한글 character. Does not support 漢字.
 - `krdict-xml`: [한국어기초사전](https://krdict.korean.go.kr) Korean wordbook XML. Supports Id, 漢字, pronunciations, antonyms, and English definition fetching.
 - `naver-json`: [Naver Korean-English Dictionary](http://endic.naver.com/) wordbook JSON. Supports Id, 漢字. (experimental)
 - `naver-table`: [Naver Korean-English Dictionary](http://endic.naver.com/) wordbook printed PDF table. Supports 漢字. (experimental)

## Formatters

Options for `--formatter` flag:

 - `tsv`: (default) tab-separated values
 - `csv`: comma-separated values
 - `json`: JSON


## Lookup

Set `-lookup` to lookup words in dictionary and enhance details. Automatically set for prompt parser.

## Interactive

Set `-interactive` to enable interative mode to resolve homophone conflicts. Automatically set for prompt parser.

## Header

Options for `--header` flag:

If set, string will be prepended to output


# Web Interface

This application also has a web interface with a simplified set of options, which can easily be run locally or deployed to Heroku.

## Running Locally

Make sure you have [Go](http://golang.org/doc/install) and the [Heroku Toolbelt](https://toolbelt.heroku.com/) installed.

```sh
$ go get -u github.com/ryanbrainard/jjogaegi
$ cd $GOPATH/src/github.com/ryanbrainard/jjogaegi
$ heroku local
```

Your app should now be running on [localhost:5000](http://localhost:5000/).

## Deploying to Heroku

```sh
$ heroku create
$ git push heroku master
$ heroku open
```

or

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)

# Examples

The following are examples of using `jjogaegi` for common use cases:

## Importing 한국어기초사전 Word Books into Anki

The [National Institute of Korean Basic Dictionary (국립국어원 한국어기초사전)](https://krdict.korean.go.kr) allows word books to be downloaded in XML format, which are ideal for parsing with `jjogaegi` and creating flashcards.  The dictionary itself is designed for learners with simple definitions, audio samples, example sentences, and is translated into 10 languages. The word books can currently only be created in the Korean dictionary; however, `jjogaegi` will download the English definitions and Korean audio during processing. To keep flashcards simple, only the first definition and two examples are exported per entry.

Here's how to export XML from the dictionary:

1. [Create a login](https://krdict.korean.go.kr/login/privateForm).
1. [Sign in](https://krdict.korean.go.kr/login/login).
1. Search for a word and go to the entry. Be sure you are in the basic Korean dictionary, not the Korean-English dictionary.
1. Click **단어장에 추가** and follow the prompts to add the word to a word book.
1. Add more words.
1. Click **내 정보 관리** and then go to the **내 단어장** tab.
1. Select the words you want to export
1. Click **단어장 내려받기**, choose **XML**, and **내려받기**.
1. Go to the **사전 내려받기** tab.
1. Find the completed export and click **내려받기**.
1. Save the file to your computer (shown as `input.xml` below).

Now, you can use `jjogaegi` to convert the XML file into a TSV file for [Anki](https://apps.ankiweb.net/) . Make sure you have `jjogaegi` [installed](#installation) and run:

```sh
$ cat input.xml | jjogaegi --parser krdict-xml --header 'tags: example' --audiodir '/path/to/anki/media' > output.tsv
```

Let's break down this command a bit:

 - `cat input.xml |` prints the XML file and pipes it to the next command.
 - `jjogaegi` processes the XML with options:
     - `--parser krdict-xml` sets the parser for this dictionary.
     - `--header 'tags: example'` sets [tags for Anki](https://apps.ankiweb.net/docs/manual.html#adding-tags).
     - `--audiodir '/path/to/anki/media'` sets the the directory to download audio files. See [Anki file locations](https://apps.ankiweb.net/docs/manual.html#files) for details on what this should be set to on your computer.
 - `> output.tsv` save the output to a file.

Now, open Anki, go to **File**, **Import**, choose the TSV file and follow the wizard.

# Disclaimer

For personal use only. Do not use this tool for publishing copyrighted content. Respect copyright holders' rights of their content. Logos shown above are copyrighted by their respective owners.

# Thanks to JetBrains

Special thanks to [JetBrains](http://www.jetbrains.com) for supporting open source software and providing a license for developing this application.

[![jetbrains logo](assets/jetbrains.svg)](http://www.jetbrains.com)