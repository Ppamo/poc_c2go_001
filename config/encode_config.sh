#!/bin/bash

CONFIG='server.json'

printf "=> Encoding config with base64\n"
ENCODED=$(cat $CONFIG | base64)
ENCODED=${ENCODED//$'\n'/}
ENCODED=${ENCODED:: -2}
printf "> LEN: %d\nTEXT: %s\n" "${#ENCODED}" "$ENCODED"
echo


printf "=> Generating random data\n"
RANDTEXT=$(dd if=/dev/urandom bs=256 count=1 2>/dev/null | base64 )
RANDTEXT=${RANDTEXT//$'\n'/}
RANDTEXT=${RANDTEXT//=/}
RANDTEXT_LENGTH=${#RANDTEXT}
printf "> LEN: %d\nTEXT: %s\n" "$RANDTEXT_LENGTH" "$RANDTEXT"
echo

printf "=> Generating prefix\n"
CHAR=${RANDTEXT:0:1}
CHARVALUE=$(printf "%d" \'$CHAR)
printf "> CHAR: %c - VALUE: %d\n" "$CHAR" "$CHARVALUE"
PREFIX=${RANDTEXT:0:$CHARVALUE}
printf "> LEN: %d\nTEXT: %s\n" "${#PREFIX}" "$PREFIX"
echo

printf "=> Generating suffix\n"
CHAR=${RANDTEXT: -1}
CHARVALUE=$(printf "%d" \'$CHAR)
printf "> CHAR: %c - VALUE: %d\n" "$CHAR" "$CHARVALUE"
SUFFIX=${RANDTEXT: -$CHARVALUE}
printf "> LEN: %d\nTEXT: %s\n" "${#SUFFIX}" "$SUFFIX"
echo

printf "=> Generating encoded text\n"
OUTPUT=$(printf "%s%s%s" "$PREFIX" "$ENCODED" "$SUFFIX")
OUTPUT=$(echo -n "$OUTPUT" | base64)
OUTPUT=${OUTPUT//$'\n'/}
echo "$OUTPUT" > $CONFIG.encoded
printf "> OUTPUT: %s\n" "$OUTPUT"
