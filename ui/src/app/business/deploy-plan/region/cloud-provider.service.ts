import { Injectable } from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {BaseModelService} from '../../../shared/class/BaseModelService';

@Injectable({
  providedIn: 'root'
})
export class CloudProviderService extends BaseModelService<any>{

  baseUrl = '/api/v1/cloud/providers/';

  constructor(http: HttpClient) {
    super(http);
  }
}
