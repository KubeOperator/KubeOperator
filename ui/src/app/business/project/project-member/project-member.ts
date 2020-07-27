import {BaseModel, BaseRequest} from '../../../shared/class/BaseModel';

export class ProjectMember extends BaseModel {
    userName: string;
    role: string;
    email: string;
}

export class ProjectMemberRequest extends BaseRequest {
    name: string;
    role: string;
    projectId: string;
}

export class ProjectMemberResponse {
    items: string[];
}

export class ProjectMemberCreate extends BaseRequest {
    projectName: string;
    userName: string;
    role: string;
}
