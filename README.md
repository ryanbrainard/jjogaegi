# 쪼개기

Put an unstructured list of Korean vocabulary in, get [Quizlet](https://quizlet.com/)-ready tab-separated values out. 

## Setup

    $ go install

## Usage

Pipe a list of 한글 terms with English definitions into `jjogaegi` and they will be parsed into tab-separated values:

    $ cat <<EOF | jjogaegi
    •입력하다 to input
    •올리다 to post, to upload
    •읽다 to read
    EOF
    입력하다	to input
    올리다	to post, to upload
    읽다	to read

The output can be used with Quizlet's [import](https://quizlet.com/help/2444107/convert-a-word-doc-into-a-quizlet-set) feature to create flashcards.

On MacOS, this works great with `pbpaste` and `ppbcopy` to avoid manually copying and pasting the input and output:

    $  pbpaste | jjogaegi | pbcopy 
