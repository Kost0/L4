#!/bin/bash

echo "=== Тест 1: разделение строк с помощью знака ',' в нескольких файлах==="

echo "Распределенный cut"
echo "Перед началом необходимо запустить воркеры"

time go run ./cmd/master test1.txt test2.tx test3.txt test4.txt test5.txt -f 1 -d , -s

echo -e "\nОригинальный cut"

time cut test1.txt test2.txt test3.txt test4.txt test5.txt -f 1 -d , -s


