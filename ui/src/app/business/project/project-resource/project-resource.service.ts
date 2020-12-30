import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {BaseModelService} from '../../../shared/class/BaseModelService';
import {ProjectResource} from './project-resource';
import {Observable} from 'rxjs';
import {Page} from '../../../shared/class/Page';

@Injectable({
    providedIn: 'root'
})
export class ProjectResourceService extends BaseModelService<ProjectResource> {

    baseUrl = '/api/v1/project/resources';

    constructor(http: HttpClient) {
        super(http);
    }

    pageBy(page, size, projectName, resourceType): Observable<Page<ProjectResource>> {
        const pageUrl = `${this.baseUrl}/?pageNum=${page}&pageSize=${size}&resourceType=${resourceType}&project=${projectName}`;
        return this.http.get<Page<ProjectResource>>(pageUrl);
    }

    listResources(resourceType, projectName): Observable<any> {
        const resourceUrl = `${this.baseUrl}/list/?resourceType=${resourceType}&project=${projectName}`;
        return this.http.get<any>(resourceUrl);
    }
}
