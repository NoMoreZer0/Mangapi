# Mangapi
Api to interract with Mangas

| Reques Type        | Path           | Description | 
| ------------- |:-------------:| -----:| 
| GET | /v1/healthcheck | check if application is working | 
| POST | /v1/users | register user |
| PUT | /v1/users/activated | to activate user | 
| POST | /v1/tokens/authentication | authenitcate a user and get token |
| POST | /v1/mangas | create manga | 
| GET | /v1/mangas | list mangas (with several options) | 
| GET | /v1/mangas/:id | list one manga | 
| PATCH | /v1/mangas/:id | update one manga information | 
| DELETE | /v1/mangas/:id | delete one manga | 
