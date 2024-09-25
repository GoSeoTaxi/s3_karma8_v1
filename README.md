ReadME.md

git clone https://github.com/GoSeoTaxi/s3_karma8_v1.git
cd s3_karma8_v1
make docker-up



Use:

cd /tmp
dd if=/dev/urandom bs=1G count=10 of=tmpfile.ext 
md5sum tmpfile.ext
curl -X POST "http://localhost:8080/files/upload" -F "file=@tmpfile.ext"
rm tmpfile.ext
curl -OJ "http://localhost:8080/files/tmpfile.ext/download"
md5sum tmpfile.ext
rm tmpfile.ext 


// additional commands
curl -OJ "http://localhost:8080/files/tmpfile.ext/download" - Полная загрузка файла
curl -X GET http://localhost:8080/files/tmpfile.ext/parts - Получить список частей файла
curl -OJ "http://localhost:8080/files/tmpfile.ext/parts/10" - Скачать часть файла

+++ log
dd if=/dev/urandom bs=1G count=10 of=tmpfile.ext
10+0 records in
10+0 records out
10737418240 bytes transferred in 37.928976 secs (283092753 bytes/sec)

md5sum tmpfile.ext
8cac33d02c96da7026088c8f560dcc16  tmpfile.ext

curl -X POST "http://localhost:8080/files/upload" -F "file=@tmpfile.ext"
Файл успешно загружен!

rm tmpfile.ext

curl -OJ "http://localhost:8080/files/tmpfile.ext/download"
% Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
Dload  Upload   Total   Spent    Left  Speed
100 10.0G    0 10.0G    0     0  63.0M      0 --:--:--  0:02:42 --:--:-- 26.5M

md5sum tmpfile.ext
8cac33d02c96da7026088c8f560dcc16  tmpfile.ext