#!/bin/bash

awk  -F, '
1-NR {
  gsub(/\$/, "", $5)
  printf("%d,\"%s%s%s\",%s\n", $1, $2 ? $2 "-" : "", $3 ? $3 " " : "", $4, $5)
}' data.csv > simplified.csv
