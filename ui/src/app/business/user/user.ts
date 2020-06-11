import {BaseModel, BaseRequest} from '../../shared/class/BaseModel';

export class User extends BaseModel {
    name: string;
    id: string;
    password: string;
    email: string;
    language: string;
    isActive: string;
    confirmPassword: string;
}

export class UserCreateRequest extends BaseRequest {
    name: string;
    password: string;
    email: string;
    confirmPassword: string;
}

export class ChangePasswordRequest extends BaseRequest {
    name: string;
    password: string;
    original: string;
}
