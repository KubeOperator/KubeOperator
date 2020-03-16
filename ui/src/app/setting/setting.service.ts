import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Settings} from './setting';

@Injectable({
  providedIn: 'root'
})
export class SettingService {
  baseUrl = '/api/v1/settings';

  constructor(private http: HttpClient) {
  }

  getSettings(): Observable<Settings> {
    return this.http.get<Settings>(this.baseUrl);
  }

  getSettingsByTab(t: string): Observable<Settings> {
    const url = this.baseUrl.concat(`?tab=${t}`);
    return this.http.get<Settings>(url);
  }

  updateSettings(settings: Settings, t: string): Observable<Settings> {
    const url = this.baseUrl.concat(`?tab=${t}`);
    return this.http.post<Settings>(url, settings);
  }

}
