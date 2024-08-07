#!/bin/bash

if ! command -v swag &> /dev/null
then
    echo "swag не установлен. Установите swag с помощью команды 'go install github.com/swaggo/swag/cmd/swag@latest'"
    exit 1
fi

# Путь к папке с API
API_PATH="schedule/api/domain"

# Путь к папке для документации
DOCS_PATH="schedule/docs/"

# Удаляем старую документацию
rm -rf $DOCS_PATH/*

# Создаем папку для документации, если она не существует
mkdir -p $DOCS_PATH

# Проходимся по каждой папке в api/domain/
for version_path in $API_PATH/*; do
    if [ -d "$version_path" ]; then
        # Извлекаем версию из пути (например, v1, v2)
        version=$(basename $version_path)

        # Путь к папке для текущей версии документации
        version_docs_path="$DOCS_PATH/$version"

        # Создаем папку для текущей версии документации
        mkdir -p $version_docs_path

        # Генерируем документацию для текущей версии
        swag init --ot json --dir schedule/api -g domain/$version/routes/routes.go --output $version_docs_path
        swag fmt

        echo "Документация для $version создана в $version_docs_path"
    fi
done

echo "Генерация документации завершена."
