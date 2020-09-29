import {BaseModel, BaseRequest} from '../../shared/class/BaseModel';

export class Project extends BaseModel {
    name: string;
    id: string;
    description: string;
}

export class ProjectCreateRequest extends BaseRequest {
    name: string;
    id: string;
    description: string;
    userName: string;
}
