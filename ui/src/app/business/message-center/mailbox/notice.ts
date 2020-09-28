import {BaseModel, BaseRequest} from '../../../shared/class/BaseModel';

export class Notice extends BaseModel {
    message: Message = new Message();
    readStatus: string;
    msgContent: {} = {};
    clusterName: string;
}

export class Message extends BaseModel {
    title: string;
    content: string;
    type: string;
    level: string;
    projectName: string;
}