import {Injectable} from '@angular/core';
import {BaseModelService} from '../../shared/class/BaseModelService';
import {System} from './system/system';
import {HttpClient} from '@angular/common/http';

@Injectable({
    providedIn: 'root'
})
export class SystemService extends BaseModelService<System> {

    baseUrl = '/api/v1/systemSettings';

    constructor(http: HttpClient) {
        super(http);
    }
}
