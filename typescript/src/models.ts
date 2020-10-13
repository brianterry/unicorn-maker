// This is a generated file. Modifications will be overwritten.
import { BaseModel, Dict, integer, Integer, Optional, transformValue } from 'cfn-rpdk';
import { Exclude, Expose, Type, Transform } from 'class-transformer';

export class ResourceModel extends BaseModel {
    ['constructor']: typeof ResourceModel;

    @Exclude()
    public static readonly TYPE_NAME: string = 'Brianterry::Unicorn::Maker';

    @Exclude()
    protected readonly IDENTIFIER_KEY_UID: string = '/properties/UID';

    @Expose({ name: 'UID' })
    @Transform(
        (value: any, obj: any) =>
            transformValue(String, 'UID', value, obj, []),
        {
            toClassOnly: true,
        }
    )
    UID?: Optional<string>;
    @Expose({ name: 'Name' })
    @Transform(
        (value: any, obj: any) =>
            transformValue(String, 'name', value, obj, []),
        {
            toClassOnly: true,
        }
    )
    Name?: Optional<string>;
    @Expose({ name: 'Color' })
    @Transform(
        (value: any, obj: any) =>
            transformValue(String, 'color', value, obj, []),
        {
            toClassOnly: true,
        }
    )
    Color?: Optional<string>;

    @Exclude()
    public getPrimaryIdentifier(): Dict {
        const identifier: Dict = {};
        if (this.UID != null) {
            identifier[this.IDENTIFIER_KEY_UID] = this.UID;
        }

        // only return the identifier if it can be used, i.e. if all components are present
        return Object.keys(identifier).length === 1 ? identifier : null;
    }

    @Exclude()
    public getAdditionalIdentifiers(): Array<Dict> {
        const identifiers: Array<Dict> = new Array<Dict>();
        // only return the identifiers if any can be used
        return identifiers.length === 0 ? null : identifiers;
    }
}
