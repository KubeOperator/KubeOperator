import {Injectable} from '@angular/core';
import {BaseModelService} from '../../../shared/class/BaseModelService';
import {VmConfig} from './vm-config';
import {HttpClient} from '@angular/common/http';

@Injectable({
    providedIn: 'root'
})
export class VmConfigService extends BaseModelService<VmConfig> {

    baseUrl = '/api/v1/vm/configs';

    constructor(httpClient: HttpClient) {
        super(httpClient);
    }
}
