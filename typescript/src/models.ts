// This is a generated file. Modifications will be overwritten.
import { BaseResourceModel, Optional } from 'cfn-rpdk';
import { allArgsConstructor, builder } from 'tombok';

@builder
@allArgsConstructor
export class ResourceModel extends BaseResourceModel {
    ['constructor']: typeof ResourceModel;
    public static readonly TYPE_NAME: string = 'Brianterry::Unicorn::Maker';

    uID: Optional<string>;
    name: Optional<string>;
    color: Optional<string>;

    constructor(...args: any[]) {super()}
    public static builder() {}
}

