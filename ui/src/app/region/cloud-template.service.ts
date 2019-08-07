import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {CloudTemplate} from './region';

@Injectable({
  providedIn: 'root'
})
export class CloudTemplateService {
  baseUrl = '/api/v1/provider/template/';

  constructor(private http: HttpClient) {
  }

  listCloudTemplate(): Observable<CloudTemplate[]> {
    return this.http.get<CloudTemplate[]>(this.baseUrl);
  }

  getCloudTemplate(name: string): Observable<CloudTemplate> {
    return this.http.get<CloudTemplate>(this.baseUrl + name + '/');
  }
}
