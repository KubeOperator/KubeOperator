import {BaseModel, BaseRequest} from '../../../../shared/class/BaseModel';

export class F5 extends BaseModel {
    id: string;
    clusterName: string;
    url: string;
    user: string;
    password: string;
    partition: string;
    publicIP: string;
    status: string;
}
