import {Injectable} from '@angular/core';
import {Task} from './task';
import {HttpClient} from '@angular/common/http';
import {BaseModelService} from '../../../../shared/class/BaseModelService';

@Injectable({
    providedIn: 'root'
})
export class TaskService extends BaseModelService<Task> {
    baseUrl = '/api/v1/cluster/{cluster_name}/tasks';

    constructor(http: HttpClient) {
        super(http);
    }
}
