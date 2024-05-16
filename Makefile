hello:
	echo Hello World

clean:
	@echo "Cleaning up..."
	rm -f myapp

run:
	@echo "Starting Docker Compose..."
	docker-compose up -d
