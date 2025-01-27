cd src
cd server
docker build . -t scooter-server:0.3
cd ..

cd etc
docker build -t spa:0.1 .

cd ../docker
docker compose up -d
