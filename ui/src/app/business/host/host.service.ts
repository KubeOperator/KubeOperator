import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {BaseModelService} from '../../shared/class/BaseModelService';

@Injectable({
    providedIn: 'root'
})
export class HostService extends BaseModelService<any> {

    baseUrl = '/api/v1/hosts';

    constructor(http: HttpClient) {
        super(http);
    }
}
