import {BaseModel, BaseRequest} from '../../../../shared/class/BaseModel';

export class F5 extends BaseModel {
    id: string;
    clusterID: string;
    url: string;
    user: string;
    password: string;
    partition: string;
    publicIP: string;
    status: boolean;
}

export class F5CreateRequest extends BaseModel {
    url: string;
    user: string;
    password: string;
    partition: string;
    publicIP: string;
}
