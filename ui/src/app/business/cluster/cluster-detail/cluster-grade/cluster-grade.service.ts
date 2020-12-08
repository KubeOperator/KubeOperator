import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Grade} from './grade';

@Injectable({
    providedIn: 'root'
})
export class ClusterGradeService {

    baseUrl = '/api/v1/clusters/grade';

    constructor(private http: HttpClient) {
    }

    getGrade(clusterName): Observable<Grade> {
        const gradeUrl = `${this.baseUrl}/${clusterName}`;
        return this.http.get<Grade>(gradeUrl);
    }
}
