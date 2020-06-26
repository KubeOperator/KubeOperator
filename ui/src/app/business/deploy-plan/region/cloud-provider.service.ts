import { Injectable } from '@angular/core';
import {BaseModelService} from "../../../shared/class/BaseModelService";
import {CloudProvider} from "./cloud-provider";
import {HttpClient} from "@angular/common/http";

@Injectable({
  providedIn: 'root'
})
export class CloudProviderService extends BaseModelService<CloudProvider>{

  baseUrl = '/api/v1/cloud/providers/'

  constructor(http: HttpClient) {
    super(http);
  }
}
