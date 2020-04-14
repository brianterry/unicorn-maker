import {
    Action,
    BaseResource,
    exceptions,
    handlerEvent,
    HandlerErrorCode,
    OperationStatus,
    Optional,
    ProgressEvent,
    ResourceHandlerRequest,
    SessionProxy,
} from 'cfn-rpdk';
import fetch, { Response } from 'node-fetch';

import { ResourceModel } from './models';


// Use this logger to forward log messages to CloudWatch Logs.
const LOGGER = console;
const CRUD_CRUD_ID = '<CRUDCRUD_API_ID_GOES_HERE>';
const API_ENDPOINT = `https://crudcrud.com/api/${CRUD_CRUD_ID}/unicorns`;
const DEFAULT_HEADERS = {
    'Accept': 'application/json; charset=utf-8',
    'Content-Type': 'application/json; charset=utf-8'
};

const checkedResponse = async (response: Response, uid?: string): Promise<any> => {
    if (response.status === 404) {
        throw new exceptions.NotFound(ResourceModel.TYPE_NAME, uid);
    } else if (![200, 201].includes(response.status)) {
        throw new exceptions.InternalFailure(
            `crudcrud.com error ${response.status} ${response.statusText}`,
            HandlerErrorCode.InternalFailure,
        );
    }
    const data = await response.text() || '{}';
    LOGGER.debug(`HTTP response ${data}`);
    return JSON.parse(data);
}

class Resource extends BaseResource<ResourceModel> {

    @handlerEvent(Action.Create)
    public async create(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: Map<string, any>,
    ): Promise<ProgressEvent> {
        LOGGER.debug(`CREATE request ${JSON.stringify(request)}`);
        const model: ResourceModel = request.desiredResourceState;
        const progress: ProgressEvent<ResourceModel> = ProgressEvent.builder()
            .status(OperationStatus.InProgress)
            .resourceModel(model)
            .build() as ProgressEvent<ResourceModel>;
        const body: Object = model.toObject();
        LOGGER.debug(`CREATE body ${JSON.stringify(body)}`);
        const response: Response = await fetch(API_ENDPOINT, {
            method: 'POST',
            headers: DEFAULT_HEADERS,
            body: JSON.stringify(body),
        });
        const jsonData: any = await checkedResponse(response);
        progress.resourceModel.UID = jsonData['_id'];
        progress.status = OperationStatus.Success;
        return progress;
    }

    @handlerEvent(Action.Update)
    public async update(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: Map<string, any>,
    ): Promise<ProgressEvent> {
        LOGGER.debug(`UPDATE request ${JSON.stringify(request)}`);
        const model: ResourceModel = request.desiredResourceState;
        const progress: ProgressEvent<ResourceModel> = ProgressEvent.builder()
            .status(OperationStatus.InProgress)
            .resourceModel(model)
            .build() as ProgressEvent<ResourceModel>;
        const body: any = model.toObject();
        delete body['UID'];
        LOGGER.debug(`UPDATE body ${JSON.stringify(body)}`);
        const response: Response = await fetch(`${API_ENDPOINT}/${model.UID}`, {
            method: 'PUT',
            headers: DEFAULT_HEADERS,
            body: JSON.stringify(body),
        });
        await checkedResponse(response, model.UID);
        progress.status = OperationStatus.Success;
        return progress;
    }

    @handlerEvent(Action.Delete)
    public async delete(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: Map<string, any>,
    ): Promise<ProgressEvent> {
        LOGGER.debug(`DELETE request ${JSON.stringify(request)}`);
        const model: ResourceModel = request.desiredResourceState;
        const progress: ProgressEvent<ResourceModel> = ProgressEvent.builder()
            .status(OperationStatus.InProgress)
            .resourceModel(model)
            .build() as ProgressEvent<ResourceModel>;
        const response: Response = await fetch(`${API_ENDPOINT}/${model.UID}`, {
            method: 'DELETE',
            headers: DEFAULT_HEADERS,
        });
        await checkedResponse(response, model.UID);
        progress.status = OperationStatus.Success;
        return progress;
    }

    @handlerEvent(Action.Read)
    public async read(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: Map<string, any>,
    ): Promise<ProgressEvent> {
        LOGGER.debug(`READ request ${JSON.stringify(request)}`);
        const model: ResourceModel = request.desiredResourceState;
        const response: Response = await fetch(`${API_ENDPOINT}/${model.UID}`, {
            method: 'GET',
            headers: DEFAULT_HEADERS,
        });
        const jsonData: any = await checkedResponse(response, model.UID);
        model.Name = jsonData['Name'];
        model.Color = jsonData['Color'];
        const progress: ProgressEvent = ProgressEvent.builder()
            .status(OperationStatus.Success)
            .resourceModel(model)
            .build();
        return progress;
    }

    @handlerEvent(Action.List)
    public async list(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: Map<string, any>,
    ): Promise<ProgressEvent> {
        LOGGER.debug(`LIST request ${JSON.stringify(request)}`);
        const response: Response = await fetch(API_ENDPOINT, {
            method: 'GET',
            headers: DEFAULT_HEADERS,
        });
        const jsonData: any[] = await checkedResponse(response);
        const models: Array<ResourceModel> = jsonData.map((unicorn: any) => {
            return new ResourceModel(new Map(Object.entries({
                UID: unicorn['_id'],
                Name: unicorn['Name'],
                Color: unicorn['Color'],
            })));
        });
        const progress: ProgressEvent = ProgressEvent.builder()
            .status(OperationStatus.Success)
            .resourceModels(models)
            .build();
        return progress;
    }
}

const resource = new Resource(ResourceModel.TYPE_NAME, ResourceModel);

export const entrypoint = resource.entrypoint;

export const testEntrypoint = resource.testEntrypoint;
