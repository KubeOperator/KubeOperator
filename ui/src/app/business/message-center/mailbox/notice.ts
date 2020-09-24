import {BaseModel, BaseRequest} from '../../../shared/class/BaseModel';

export class Notice extends BaseModel {
    content: string;
    type: string;
    level: string;
    isRead: boolean;
}

export class NoticeCreateRequest extends BaseRequest {
    content: string;
    type: string;
    level: string;
    isRead: boolean;
}
