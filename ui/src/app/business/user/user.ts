import {BaseModel, BaseRequest} from '../../shared/class/BaseModel';

export class User extends BaseModel {
    name: string;
    id: string;
    password: string;
    language: string;
    isActive: boolean;
    confirmPassword: string;
    isAdmin: boolean;
    type: string;
}

export class UserCreateRequest extends BaseRequest {
    name: string;
    password: string;
    confirmPassword: string;
    isAdmin: boolean;
}

export class ChangePasswordRequest extends BaseRequest {
    name: string;
    password: string;
    original: string;
}
