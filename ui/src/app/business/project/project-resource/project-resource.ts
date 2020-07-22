import {BaseModel, BaseRequest} from '../../../shared/class/BaseModel';

export class ProjectResource extends BaseModel {
    name: string;
}

export class ProjectResourceCreateRequest extends BaseRequest {
    projectId: string;
    resourceType: string;
    resourceId: string;
}

export class ProjectResourceCheck {
    checked: boolean;
    data: {} = {};
}
