#!/bin/bash

echo "Создание патча..."

# Читаем версию из файла
file="./version"
version=$(head -1 "$file") # = [version](Get-Content $file | Select -First 1)
if [[ "$version" -eq "" ]]; then
  version="1.0.0"
fi
echo $version

# Увеличиваем версию
version=$(./increment_version.sh -p "$version")
echo "$version" > "$file"

# Отправляем новую версию в git
commit="version $version"
echo "git commit $commit"
git add -A
git commit -a -m "$commit"
git push origin

tag="v$version"
echo "git tag $tag"
git tag $tag
git push origin $tag