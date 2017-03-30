build :
	# cd frontend; npm run build; cd -
	go-bindata-assetfs -pkg web assets/... ; mv bindata_assetfs.go web
	go build main.go
