
```bash
# билд 
go build  -o . ./...
# запуск 
./download2json

# билд в докере
docker build -t godownloader2json .
# запуск в докере
docker run godownloader2json
```


Реализовать парсер с сайта https://jsonplaceholder.typicode.com/ в несколько потоков объекты albums/ и photos/.
https://jsonplaceholder.typicode.com/albums/
https://jsonplaceholder.typicode.com/photos/
Скачиваем все альбомы и фотографии, кладем их по папкам /альбом/название_фотографии
