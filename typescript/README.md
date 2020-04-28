# Brianterry::Unicorn::Maker

Congratulations on starting development! Next steps:

1. Check the JSON schema describing your resource, [brianterry-unicorn-maker.json](./brianterry-unicorn-maker.json)
2. Run `npm install --optional`
3. Implement/modify your resource handlers in [handlers.ts](./src/handlers.ts)
4. Build the project with `npm run build`

> Don't modify [models.ts](./src/models.ts) by hand, any modifications will be overwritten when the `generate` or `package` commands are run.

Implement CloudFormation resource here. Each function must always return a ProgressEvent.

While importing the [cfn-rdpk](https://github.com/eduardomourar/cloudformation-cli-typescript-plugin) library, failures can be passed back to CloudFormation by either raising an exception from `exceptions`, or setting the ProgressEvent's `status` to `OperationStatus.Failed` and `errorCode` to one of `HandlerErrorCode`. There is a static helper function, `ProgressEvent.failed`, for this common case.
