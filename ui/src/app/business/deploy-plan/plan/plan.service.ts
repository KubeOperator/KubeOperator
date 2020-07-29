import {Injectable} from '@angular/core';
import {BaseModelService} from '../../../shared/class/BaseModelService';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Plan} from './plan';

@Injectable({
    providedIn: 'root'
})
export class PlanService extends BaseModelService<Plan> {

    baseUrl = '/api/v1/plans';


    constructor(http: HttpClient) {
        super(http);
    }

    listVmConfigs(regionName: string): Observable<any> {
        const itemUrl = `${this.baseUrl}/configs/` + regionName + '/';
        return this.http.get<any>(itemUrl);
    }

    listByProjectName(projectName: string): Observable<any> {
        const itemUrl = `${this.baseUrl}/?projectName=${projectName}`;
        return this.http.get<any>(itemUrl);
    }
}
