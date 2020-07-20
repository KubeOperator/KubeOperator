import { Injectable } from '@angular/core';
import {BaseModelService} from '../../shared/class/BaseModelService';
import {Project} from './project';
import {HttpClient} from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class ProjectService extends BaseModelService<Project>{

  baseUrl = '/api/v1/projects';

  constructor(http: HttpClient) {
    super(http);
  }

}
