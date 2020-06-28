import {Injectable} from '@angular/core';
import {BaseModelService} from '../../../shared/class/BaseModelService';
import {HttpClient} from '@angular/common/http';
import {Zone} from './zone';

@Injectable({
    providedIn: 'root'
})
export class ZoneService extends BaseModelService<Zone> {

    baseUrl = '/api/v1/zones';

    constructor(http: HttpClient) {
        super(http);
    }

}
