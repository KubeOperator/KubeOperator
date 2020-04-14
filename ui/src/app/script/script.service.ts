import {Injectable} from '@angular/core';
import {ModelService} from '../shared/class/model-service';
import {Script} from './script';
import {HttpClient} from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class ScriptService extends ModelService<Script> {

  baseUrl = '/api/v1/scripts/';

  constructor(http: HttpClient) {
    super(http);
  }
}
