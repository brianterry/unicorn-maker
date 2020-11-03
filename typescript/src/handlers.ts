import {
    Action,
    BaseResource,
    exceptions,
    handlerEvent,
    HandlerErrorCode,
    LoggerProxy,
    OperationStatus,
    Optional,
    ProgressEvent,
    ResourceHandlerRequest,
    SessionProxy,
} from 'cfn-rpdk';
import fetch, { Response } from 'node-fetch';

import { ResourceModel } from './models';

interface Unicorn {
    readonly _id?: string;
    readonly Name?: string;
    readonly Color?: string;
}

class Resource extends BaseResource<ResourceModel> {
    // APIEndpoint is the crudcrud endpoint.
    // You can  obtain an endpoint by going to https://crudcrud.com.
    static readonly API_ENDPOINT = `https://crudcrud.com/api/<Your API ID>/unicorns`;
    static readonly DEFAULT_HEADERS = {
        'Accept': 'application/json; charset=utf-8',
        'Content-Type': 'application/json; charset=utf-8'
    };

    private async checkResponse(response: Response, logger: LoggerProxy, uid?: string): Promise<any> {
        if (response.status === 404) {
            throw new exceptions.NotFound(this.typeName, uid);
        } else if (![200, 201].includes(response.status)) {
            throw new exceptions.InternalFailure(
                `crudcrud.com error ${response.status} ${response.statusText}`,
                HandlerErrorCode.InternalFailure,
            );
        }
        const data = await response.text() || '{}';
        logger.log(`HTTP response ${data}`);
        return JSON.parse(data);
    }

    @handlerEvent(Action.Create)
    public async create(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: Map<string, any>,
        logger: LoggerProxy
    ): Promise<ProgressEvent<ResourceModel>> {
        logger.log('CREATE request', request);
        const model = new ResourceModel(request.desiredResourceState);
        if (model.id) {
            throw new exceptions.InvalidRequest('Create unicorn with readOnly property');
        }
        const progress = ProgressEvent.progress<ProgressEvent<ResourceModel>>(model);
        const body = model.toJSON();
        logger.log('UPDATE body', body);
        const response: Response = await fetch(Resource.API_ENDPOINT, {
            method: 'POST',
            headers: Resource.DEFAULT_HEADERS,
            body: JSON.stringify(body),
        });
        const jsonData: Unicorn = await this.checkResponse(response, logger);
        progress.resourceModel.id = jsonData._id;
        progress.status = OperationStatus.Success;
        logger.log('CREATE progress', progress);
        return progress;
    }

    @handlerEvent(Action.Update)
    public async update(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: Map<string, any>,
        logger: LoggerProxy
    ): Promise<ProgressEvent<ResourceModel>> {
        logger.log('UPDATE request', request);
        const model = new ResourceModel(request.desiredResourceState);
        const progress = ProgressEvent.progress<ProgressEvent<ResourceModel>>(model);
        const body = model.toJSON();
        delete body['Id'];
        logger.log('UPDATE body', body);
        const response: Response = await fetch(`${Resource.API_ENDPOINT}/${model.id}`, {
            method: 'PUT',
            headers: Resource.DEFAULT_HEADERS,
            body: JSON.stringify(body),
        });
        await this.checkResponse(response, logger, model.id);
        progress.status = OperationStatus.Success;
        logger.log('UPDATE progress', progress);
        return progress;
    }

    @handlerEvent(Action.Delete)
    public async delete(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: Map<string, any>,
        logger: LoggerProxy
    ): Promise<ProgressEvent<ResourceModel>> {
        logger.log('DELETE request', request);
        const model = new ResourceModel(request.desiredResourceState);
        const progress = ProgressEvent.progress<ProgressEvent<ResourceModel>>();
        const response: Response = await fetch(`${Resource.API_ENDPOINT}/${model.id}`, {
            method: 'DELETE',
            headers: Resource.DEFAULT_HEADERS,
        });
        await this.checkResponse(response, logger, model.id);
        progress.status = OperationStatus.Success;
        logger.log('DELETE progress', progress);
        return progress;
    }

    @handlerEvent(Action.Read)
    public async read(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: Map<string, any>,
        logger: LoggerProxy
    ): Promise<ProgressEvent<ResourceModel>> {
        logger.log('READ request', request);
        const model = new ResourceModel(request.desiredResourceState);
        const response: Response = await fetch(`${Resource.API_ENDPOINT}/${model.id}`, {
            method: 'GET',
            headers: Resource.DEFAULT_HEADERS,
        });
        const jsonData: Unicorn = await this.checkResponse(response, logger, model.id);
        model.name = jsonData.Name;
        model.color = jsonData.Color;
        const progress = ProgressEvent.success<ProgressEvent<ResourceModel>>(model);
        logger.log('READ progress', progress);
        return progress;
    }

    @handlerEvent(Action.List)
    public async list(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: Map<string, any>,
        logger: LoggerProxy
    ): Promise<ProgressEvent<ResourceModel>> {
        logger.log('LIST request', request);
        const response: Response = await fetch(Resource.API_ENDPOINT, {
            method: 'GET',
            headers: Resource.DEFAULT_HEADERS,
        });
        const jsonData: Unicorn[] = await this.checkResponse(response, logger);
        const progress = ProgressEvent.success<ProgressEvent<ResourceModel>>();
        progress.resourceModels = jsonData.map((unicorn: Unicorn) => {
            return new ResourceModel({
                id: unicorn._id,
                name: unicorn.Name,
                color: unicorn.Color,
            });
        });
        logger.log('LIST progress', progress);
        return progress;
    }
}

export const resource = new Resource(ResourceModel.TYPE_NAME, ResourceModel);

export const entrypoint = resource.entrypoint;

export const testEntrypoint = resource.testEntrypoint;
