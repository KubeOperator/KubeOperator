import {Injectable} from '@angular/core';
import {BaseModelService} from '../../shared/class/BaseModelService';
import {HttpClient} from '@angular/common/http';
import {Host, HostCreateRequest} from '../host/host';
import {Observable} from 'rxjs';
import {ChangePasswordRequest} from './user';
import {ResetPassword} from '../../login/forgot-password/reset-password';

@Injectable({
    providedIn: 'root'
})
export class UserService extends BaseModelService<any> {


    baseUrl = '/api/v1/users';

    constructor(http: HttpClient) {
        super(http);
    }

    changePassword(item: ChangePasswordRequest): Observable<Host> {
        const itemUrl = `${this.baseUrl}/change/password/`;
        return this.http.post<Host>(itemUrl, item);
    }

    resetPassword(item: ResetPassword): Observable<any> {
        const itemUrl = `/api/v1/user/forgot/password/`;
        return this.http.post<Host>(itemUrl, item);
    }
}
