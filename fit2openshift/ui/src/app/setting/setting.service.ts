import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Setting} from './setting';

@Injectable({
  providedIn: 'root'
})
export class SettingService {
  baseUrl = '/api/v1/setting/';

  constructor(private http: HttpClient) {
  }

  listSettings(): Observable<Setting[]> {
    return this.http.get<Setting[]>(this.baseUrl);
  }

  updateSetting(key: string, setting: Setting): Observable<Setting> {
    return this.http.patch<Setting>(this.baseUrl + key + '/', setting);
  }
}
