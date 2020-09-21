import {BaseModel} from '../../shared/class/BaseModel';

export class Message {

}

export class UserReceiver extends BaseModel {
    id: string;
    userId: string;
    vars: {} = {};
}

export class UserNotificationConfig extends BaseModel {
    id: string;
    userId: string;
    vars: {} = {};
    type: string;
}
