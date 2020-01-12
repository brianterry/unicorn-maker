import logging


import requests
from cloudformation_cli_python_lib import (
    Action,
    OperationStatus,
    ProgressEvent,
    Resource,
    exceptions,
)

from .models import ResourceHandlerRequest, ResourceModel

# Use this logger to forward log messages to CloudWatch Logs.
LOG = logging.getLogger(__name__)
TYPE_NAME = "Brianterry::Unicorn::Maker"

resource = Resource(TYPE_NAME, ResourceModel)
test_entrypoint = resource.test_entrypoint

crud_crud_id = "<CRUDCRUD_API_ID_GOES_HERE>"
api_endpoint = f"https://crudcrud.com/api/{crud_crud_id}/unicorns"


def check_response(response, uid=None):
    if response.status_code == 404:
        raise exceptions.NotFound(TYPE_NAME, uid)
    if response.status_code not in [200, 201]:
        raise exceptions.InternalFailure(
            f"crudcrud.com error {response.status_code} {response.reason}"
        )



@resource.handler(Action.CREATE)
def create_handler(_, request: ResourceHandlerRequest, __) -> ProgressEvent:
    model = request.desiredResourceState
    progress: ProgressEvent = ProgressEvent(
        status=OperationStatus.SUCCESS, resourceModel=model,
    )
    response = requests.post(api_endpoint, json=model._serialize())
    check_response(response)
    progress.resourceModel.UID = response.json()["_id"]
    return progress


@resource.handler(Action.UPDATE)
def update_handler(_, request: ResourceHandlerRequest, __) -> ProgressEvent:
    model = request.desiredResourceState
    progress: ProgressEvent = ProgressEvent(
        status=OperationStatus.SUCCESS, resourceModel=model,
    )
    response = requests.put(f"{api_endpoint}/{model.UID}", json=model._serialize())
    check_response(response, model.UID)
    return progress


@resource.handler(Action.DELETE)
def delete_handler(_, request: ResourceHandlerRequest, __) -> ProgressEvent:
    model = request.desiredResourceState
    progress: ProgressEvent = ProgressEvent(
        status=OperationStatus.SUCCESS, resourceModel=model,
    )
    response = requests.delete(f"{api_endpoint}/{model.UID}")
    check_response(response, model.UID)
    return progress


@resource.handler(Action.READ)
def read_handler(_, request: ResourceHandlerRequest, __,) -> ProgressEvent:
    model = request.desiredResourceState
    response = requests.get(f"{api_endpoint}/{model.UID}")
    check_response(response, model.UID)
    model.Name = response.json()["Name"]
    model.Color = response.json()["Color"]
    return ProgressEvent(
        status=OperationStatus.SUCCESS,
        resourceModel=model,
    )


@resource.handler(Action.LIST)
def list_handler(_, __, ___) -> ProgressEvent:
    models = []
    response = requests.get(api_endpoint)
    check_response(response)
    for unicorn in response.json():
        models.append(ResourceModel(unicorn["_id"], unicorn["Name"], unicorn["Color"]))
    return ProgressEvent(
        status=OperationStatus.SUCCESS,
        resourceModels=models,
    )
