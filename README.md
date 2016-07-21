# 쪼개기

Put an unstructured list of Korean vocabulary in, get [Quizlet](https://quizlet.com/)-ready tab-separated values out. 

## Setup

    $ go install

## Usage

Pipe a list of 한글 terms with English definitions into `jjogaegi` and they will be parsed into tab-separated values. The tab character is inserted after the last 한글 character on the line, with spaces preserved and extraneous characters and headings removed. For example:

    $ cat <<EOF | jjogaegi
    My Word List
    •컴퓨터를 켜다 to turn on the computer
    •브라우저를 열다 to open the web browser
    •검색어를 입력하다 to type in the search word
    EOF
    컴퓨터를 켜다	to turn on the computer
    브라우저를 열다	to open the web browser
    검색어를 입력하다	to type in the search word

The output can then be used with Quizlet's [import](https://quizlet.com/help/2444107/convert-a-word-doc-into-a-quizlet-set) feature to create flashcards.

On MacOS, this works great with `pbpaste` and `ppbcopy` to avoid manually copying and pasting the input and output:

    $  pbpaste | jjogaegi | pbcopy 
