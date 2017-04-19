tools :
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...
	go get github.com/laher/goxc

build-fe :
	cd frontend; npm run build

build : tools
	go-bindata-assetfs -pkg web assets/... ; mv bindata_assetfs.go web

release : build-fe build
	cd cmd/
	goxc -pv="$(v)" -d="$(dest)"
