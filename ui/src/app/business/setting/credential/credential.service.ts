import {Injectable} from '@angular/core';
import {BaseModelService} from '../../../shared/class/BaseModelService';
import {HttpClient} from '@angular/common/http';
import {Credential} from './credential';

@Injectable({
    providedIn: 'root'
})
export class CredentialService extends BaseModelService<Credential> {

    baseUrl = '/api/v1/credentials';

    constructor(http: HttpClient) {
        super(http);
    }
}
