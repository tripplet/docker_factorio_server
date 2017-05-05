sudo docker run --rm -it \
	-v /factorio/server:/factorio \
	--name factorio \
	factorio "$@"
find /factorio/server -type f
