#!/bin/bash

# Параметры теста
DURATION=5          # Продолжительность
RATE=3000             # Запросов в секунду
TARGETS="targets_multi.json" # Файл со списком сайтов
OUTPUT_DIR="results"    # Папка для результатов

# Создаем папку для результатов, если ее нет
mkdir -p "$OUTPUT_DIR"

# Определяем путь к vegeta (локальный файл или из PATH)
if [ -f "./vegeta.exe" ]; then
    VEGETA="./vegeta.exe"
elif [ -f "./vegeta" ]; then
    VEGETA="./vegeta"
else
    VEGETA="vegeta"
fi

echo -e "\e[36mЗапуск нагрузочного теста...\e[0m"
echo "Интенсивность: $RATE запросов/сек"
echo "Продолжительность: $DURATION"

# Запуск Vegeta
cat $TARGETS | $VEGETA attack -format=json -duration=5s -rate=$RATE | tee "$OUTPUT_DIR/results.bin" | $VEGETA report

# Генерация HTML отчета
$VEGETA plot "$OUTPUT_DIR/results.bin" > "$OUTPUT_DIR/report.html"

echo -e "\n\e[32mТест завершен. Отчет сохранен в $OUTPUT_DIR/report.html\e[0m"
