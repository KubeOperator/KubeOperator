import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Settings} from './setting';

@Injectable({
  providedIn: 'root'
})
export class SettingService {
  baseUrl = '/api/v1/settings/';

  constructor(private http: HttpClient) {
  }

  getSettings(): Observable<Settings> {
    return this.http.get<Settings>(this.baseUrl);
  }
  updateSettings(settings: Settings): Observable<Settings> {
    return this.http.post<Settings>(this.baseUrl, settings);
  }
}
