build_containers:
				docker-compose up --build

clean_containers:
				docker rm task_api task_postgres

build_api		:
			go build -o server ./cmd/api/main.go  && ./server ${serverport}

clean_api		:
			rm server
