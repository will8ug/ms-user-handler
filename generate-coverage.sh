#!/bin/bash

srcFile=coverage.out
if [ ! -e "$srcFile" ]
then
    echo "$srcFile not found."
    exit 1
fi

coverage="0.0%"
color=green
while read line
do
    if [[ "$line" =~ ^total.*%$ ]]
    then
        coverage=`echo $line | awk '{ print $3 }'`
    fi
done < $srcFile
echo Coverage: $coverage

coverage=${coverage%\%}
if (( $(echo "$coverage <= 30" | bc -l) )) ; then
    color=red
elif (( $(echo "$coverage > 80" | bc -l) )); then
    color=green
else
    color=orange
fi

echo "Downloading https://img.shields.io/badge/Coverage-${coverage}%25-${color}"
url="https://img.shields.io/badge/Coverage-${coverage}%25-${color}"
curl -kf -m 10 --create-dirs -o badges/badge-statements.svg $url