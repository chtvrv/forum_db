PROJECT_DIR := ${CURDIR}
DOCUMENTATION_CONTAINER_NAME=documentation_db
DOCUMENTATION_FILE=swagger.yml

# документация
doc-prepare:
	npm install speccy -g
	docker pull swaggerapi/swagger-ui

doc-host:	
	docker run --name=${DOCUMENTATION_CONTAINER_NAME} -d -p 82:8080 -e SWAGGER_JSON=/${DOCUMENTATION_FILE} -v $(PROJECT_DIR)/${DOCUMENTATION_FILE}:/${DOCUMENTATION_FILE} swaggerapi/swagger-ui

doc-stop:
	docker stop ${DOCUMENTATION_CONTAINER_NAME}
	docker rm ${DOCUMENTATION_CONTAINER_NAME}

.PHONY:
	start stop
