import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class GrafanaService {
  constructor(private http: HttpClient) {
  }

  list_dashboard = '/api/search';

  list_grafana_url(baseUrl) {
    console.log(this.http.jsonp(baseUrl + this.list_dashboard, null));
  }


}
