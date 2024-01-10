#!/bin/bash

srcFile=coverage.out
if [ ! -e "$srcFile" ]
then
    echo "$srcFile not found."
    exit 1
fi

coverage="0.0%"
color=brightgreen
while read line
do
    if [[ "$line" =~ ^total.*%$ ]]
    then
        coverage=`echo $line | awk '{ print $3 }'`
        coverage=${coverage%\%}
    fi
done < $srcFile

echo
echo $coverage
echo "https://img.shields.io/badge/Coverage-${coverage}%25-${color}"
url="https://img.shields.io/badge/Coverage-${coverage}%25-${color}"
curl -kf -m 10 --create-dirs -o badges/badge-statements.svg $url