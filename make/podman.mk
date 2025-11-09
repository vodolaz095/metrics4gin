podman/up:
	podman-compose up -d
	podman ps

podman/resource:
	podman-compose up -d victoria
	podman ps

podman/down:
	podman-compose down

podman/prune:
	podman system prune -a --volumes
