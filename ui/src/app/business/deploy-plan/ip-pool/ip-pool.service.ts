import {Injectable} from '@angular/core';
import {BaseModelService} from '../../../shared/class/BaseModelService';
import {IpPool} from './ip-pool';
import {HttpClient} from '@angular/common/http';

@Injectable({
    providedIn: 'root'
})
export class IpPoolService extends BaseModelService<IpPool> {

    baseUrl = '/api/v1/ippools';

    constructor(http: HttpClient) {
        super(http);
    }
}
