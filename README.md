# lumberyard
Pipeline metadata persistence API

### _Routes_
GET / : Heartbeat, should return status 200 if alive


GET /api/v1/pipelines : Return all pipelines

POST /api/v1/pipelines : Create a new pipeline

GET /api/v1/pipeline/{pipeline_id} : Return a pipeline with id {pipeline_id}


GET /api/v1/pipeline/{pipeline_id}/stages : Return all stages for a given pipeline

POST /api/v1/pipeline/{pipeline_id}/stages : Add a new stage to a given pipeline

