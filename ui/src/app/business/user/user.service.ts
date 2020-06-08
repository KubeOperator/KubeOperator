import {Injectable} from '@angular/core';
import {BaseModelService} from '../../shared/class/BaseModelService';
import {HttpClient} from '@angular/common/http';

@Injectable({
    providedIn: 'root'
})
export class UserService extends BaseModelService<any> {


    baseUrl = '/api/v1/users';

    constructor(http: HttpClient) {
        super(http);
    }

}
