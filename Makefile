IMAGE = ddd-timer-builder:latest
VOLUME_PATH_WIN = "$(shell pwd):/mnt"
SSH_HOST = 193.178.169.195

deploy:
	docker build -t $(IMAGE) .
