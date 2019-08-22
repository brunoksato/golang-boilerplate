test:
	go test ./model ./api ./util

test-mid:
	go test ./middleware -v

test-core:
	go test ./model -v

test-util:
	go test ./util -v

test-api:
	go test ./api -v	

deploy:
	git push heroku master