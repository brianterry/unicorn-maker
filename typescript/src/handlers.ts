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

interface CallbackContext extends Record<string, any> {}

class Resource extends BaseResource<ResourceModel> {

    /**
     * CloudFormation invokes this handler when the resource is initially created
     * during stack create operations.
     *
     * @param session Current AWS session passed through from caller
     * @param request The request object for the provisioning request passed to the implementor
     * @param callbackContext Custom context object to allow the passing through of additional
     * state or metadata between subsequent retries
     */
    @handlerEvent(Action.Create)
    public async create(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: CallbackContext,
    ): Promise<ProgressEvent> {
        LOGGER.debug('CREATE request', request);
        const model: ResourceModel = request.desiredResourceState;
        if (model.UID) throw new exceptions.InvalidRequest("Create unicorn with readOnly property");

        const progress = ProgressEvent.progress<ProgressEvent<ResourceModel, CallbackContext>>(model);
        const body: Object = { ...model };
        LOGGER.debug('CREATE body', body);
        const response: Response = await fetch(API_ENDPOINT, {
            method: 'POST',
            headers: DEFAULT_HEADERS,
            body: JSON.stringify(body),
        });
        const jsonData: any = await checkedResponse(response);
        progress.resourceModel.UID = jsonData['_id'];
        progress.status = OperationStatus.Success;
        LOGGER.log('CREATE progress', { ...progress });
        return progress;
    }

    /**
     * CloudFormation invokes this handler when the resource is updated
     * as part of a stack update operation.
     *
     * @param session Current AWS session passed through from caller
     * @param request The request object for the provisioning request passed to the implementor
     * @param callbackContext Custom context object to allow the passing through of additional
     * state or metadata between subsequent retries
     */
    @handlerEvent(Action.Update)
    public async update(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: CallbackContext,
    ): Promise<ProgressEvent> {
        LOGGER.debug('UPDATE request', request);
        const model: ResourceModel = request.desiredResourceState;
        const progress = ProgressEvent.progress<ProgressEvent<ResourceModel, CallbackContext>>(model);
        const body: any = { ...model };
        delete body['UID'];
        LOGGER.debug('UPDATE body', body);
        const response: Response = await fetch(`${API_ENDPOINT}/${model.UID}`, {
            method: 'PUT',
            headers: DEFAULT_HEADERS,
            body: JSON.stringify(body),
        });
        await checkedResponse(response, model.UID);
        progress.status = OperationStatus.Success;
        LOGGER.log('UPDATE progress', { ...progress });
        return progress;
    }

    /**
     * CloudFormation invokes this handler when the resource is deleted, either when
     * the resource is deleted from the stack as part of a stack update operation,
     * or the stack itself is deleted.
     *
     * @param session Current AWS session passed through from caller
     * @param request The request object for the provisioning request passed to the implementor
     * @param callbackContext Custom context object to allow the passing through of additional
     * state or metadata between subsequent retries
     */
    @handlerEvent(Action.Delete)
    public async delete(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: CallbackContext,
    ): Promise<ProgressEvent> {
        LOGGER.debug('DELETE request', request);
        const model: ResourceModel = request.desiredResourceState;
        const progress = ProgressEvent.progress<ProgressEvent<ResourceModel, CallbackContext>>();
        const response: Response = await fetch(`${API_ENDPOINT}/${model.UID}`, {
            method: 'DELETE',
            headers: DEFAULT_HEADERS,
        });
        await checkedResponse(response, model.UID);
        progress.status = OperationStatus.Success;
        LOGGER.log('DELETE progress', { ...progress });
        return progress;
    }

    /**
     * CloudFormation invokes this handler as part of a stack update operation when
     * detailed information about the resource's current state is required.
     *
     * @param session Current AWS session passed through from caller
     * @param request The request object for the provisioning request passed to the implementor
     * @param callbackContext Custom context object to allow the passing through of additional
     * state or metadata between subsequent retries
     */
    @handlerEvent(Action.Read)
    public async read(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: CallbackContext,
    ): Promise<ProgressEvent> {
        LOGGER.debug('READ request', request);
        const model: ResourceModel = request.desiredResourceState;
        const response: Response = await fetch(`${API_ENDPOINT}/${model.UID}`, {
            method: 'GET',
            headers: DEFAULT_HEADERS,
        });
        const jsonData: any = await checkedResponse(response, model.UID);
        model.Name = jsonData['Name'];
        model.Color = jsonData['Color'];
        const progress: ProgressEvent<ResourceModel> = ProgressEvent.builder()
            .status(OperationStatus.Success)
            .resourceModel(model)
            .build() as ProgressEvent<ResourceModel>;
        LOGGER.log('READ progress', { ...progress });
        return progress;
    }

    /**
     * CloudFormation invokes this handler when summary information about multiple
     * resources of this resource provider is required.
     *
     * @param session Current AWS session passed through from caller
     * @param request The request object for the provisioning request passed to the implementor
     * @param callbackContext Custom context object to allow the passing through of additional
     * state or metadata between subsequent retries
     */
    @handlerEvent(Action.List)
    public async list(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: CallbackContext,
    ): Promise<ProgressEvent> {
        LOGGER.debug('LIST request', request);
        const response: Response = await fetch(API_ENDPOINT, {
            method: 'GET',
            headers: DEFAULT_HEADERS,
        });
        const jsonData: any[] = await checkedResponse(response);
        const models: Array<ResourceModel> = jsonData.map((unicorn: any) => {
            LOGGER.log("\n\nHERE", unicorn, "\n\n");
            return new ResourceModel({
                UID: unicorn['_id'],
                Name: unicorn['Name'],
                Color: unicorn['Color'],
            });
        });
        const progress: ProgressEvent<ResourceModel> = ProgressEvent.builder()
            .status(OperationStatus.Success)
            .resourceModels(models)
            .build() as ProgressEvent<ResourceModel>;
        LOGGER.log('LIST progress test logger', { ...progress });
        return progress;
    }
}

const resource = new Resource(ResourceModel.TYPE_NAME, ResourceModel);

export const entrypoint = resource.entrypoint;

export const testEntrypoint = resource.testEntrypoint;
