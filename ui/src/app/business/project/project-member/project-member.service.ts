import {Injectable} from '@angular/core';
import {BaseModelService} from '../../../shared/class/BaseModelService';
import {HttpClient} from '@angular/common/http';
import {ProjectMember, ProjectMemberResponse} from './project-member';
import {Observable} from 'rxjs';
import {Page} from '../../../shared/class/Page';
import {BaseRequest} from "../../../shared/class/BaseModel";

@Injectable({
    providedIn: 'root'
})
export class ProjectMemberService extends BaseModelService<any> {

    baseUrl = '/api/v1/project/members';

    constructor(http: HttpClient) {
        super(http);
    }

    pageBy(page, size, projectName): Observable<Page<ProjectMember>> {
        const pageUrl = `${this.baseUrl}?pageNum=${page}&pageSize=${size}&project=${projectName}`;
        return this.http.get<Page<ProjectMember>>(pageUrl);
    }

    getByUser(username: string, projectName: string): Observable<ProjectMember> {
        return this.http.get<ProjectMember>(`${this.baseUrl}/${username}?project=${projectName}`, {});
    }

    getUsers(name, projectName): Observable<ProjectMemberResponse> {
        const userUrl = `${this.baseUrl}/users?name=${name}`;
        return this.http.get<ProjectMemberResponse>(userUrl, {});
    }
}
