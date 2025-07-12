#!/bin/bash

# Модифицированный тест с 3 попытками на 20 одновременных играх
seq 1 20 | xargs -P 10 -I {} bash -c '
    for i in {1..3}; do
        OUTPUT=$(python3 checker.py check 127.0.0.1 2>&1);
        EXIT_CODE=$?;
        if [ $EXIT_CODE -eq 101 ]; then
            echo "✅ Проверка {} успешна (попытка $i)";
            break;
        elif [ $i -eq 3 ]; then
            echo "❌ Ошибка в процессе {} (код $EXIT_CODE)";
            echo "$OUTPUT";
        fi
    done
'