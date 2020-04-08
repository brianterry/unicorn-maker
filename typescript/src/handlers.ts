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
import { ResourceModel } from './models';


// Use this logger to forward log messages to CloudWatch Logs.
const LOGGER = console;

class Resource extends BaseResource<ResourceModel> {

    @handlerEvent(Action.Create)
    public create(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: Map<string, any>,
    ): ProgressEvent {
        const model: ResourceModel = request.desiredResourceState;
        const progress: ProgressEvent = ProgressEvent.builder()
            .status(OperationStatus.InProgress)
            .resourceModel(model)
            .build();
        // TODO: put code here

        // Example:
        try {
            if (session instanceof SessionProxy) {
                const client = session.client('S3');
            }
            // Setting Status to success will signal to cfn that the operation is complete
            progress.status = OperationStatus.Success;
        } catch(err) {
            LOGGER.log(err);
            // exceptions module lets CloudFormation know the type of failure that occurred
            throw new exceptions.InternalFailure(err.message);
            // this can also be done by returning a failed progress event
            // return ProgressEvent.failed(HandlerErrorCode.InternalFailure, err.message);
        }
        return progress;
    }

    @handlerEvent(Action.Update)
    public update(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: Map<string, any>,
    ): ProgressEvent {
        const model: ResourceModel = request.desiredResourceState;
        const progress: ProgressEvent = ProgressEvent.builder()
            .status(OperationStatus.InProgress)
            .resourceModel(model)
            .build();
        // TODO: put code here
        return progress;
    }

    @handlerEvent(Action.Delete)
    public delete(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: Map<string, any>,
    ): ProgressEvent {
        const model: ResourceModel = request.desiredResourceState;
        const progress: ProgressEvent = ProgressEvent.builder({
            status: OperationStatus.InProgress,
            resourceModel: model,
        }).build();
        // TODO: put code here
        return progress;
    }

    @handlerEvent(Action.Read)
    public read(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: Map<string, any>,
    ): ProgressEvent {
        const model: ResourceModel = request.desiredResourceState;
        // TODO: put code here
        const progress: ProgressEvent = ProgressEvent.builder()
            .status(OperationStatus.Success)
            .resourceModel(model)
            .build();
        return progress;
    }

    @handlerEvent(Action.List)
    public list(
        session: Optional<SessionProxy>,
        request: ResourceHandlerRequest<ResourceModel>,
        callbackContext: Map<string, any>,
    ): ProgressEvent {
        // TODO: put code here
        const progress: ProgressEvent = ProgressEvent.builder()
            .status(OperationStatus.Success)
            .resourceModels([])
            .build();
        return progress;
    }
}

const resource = new Resource(ResourceModel.TYPE_NAME, ResourceModel);

export const entrypoint = resource.entrypoint;

export const testEntrypoint = resource.testEntrypoint;
