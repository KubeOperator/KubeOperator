import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Host} from './host';
import {ModelService} from '../shared/class/model-service';
import {Observable} from 'rxjs';


@Injectable({
  providedIn: 'root'
})
export class HostService extends ModelService<Host> {

  baseUrl = '/api/v1/host/';
  importUrl = 'import/';

  constructor(http: HttpClient) {
    super(http);
  }

  import(file_names: string[]): Observable<any> {
    return this.http.post<any>(`${this.baseUrl}${this.importUrl}`, {'source': file_names});
  }

  byItem(itemName: string): Observable<Host[]> {
    return this.http.get<Host[]>(`${this.baseUrl}?item=${itemName}`);
  }

  // listHosts(): Observable<Host[]> {
  //   return this.http.get<Host[]>(baseUrl);
  // }
  //
  // listItemHosts(itemName: string): Observable<Host[]> {
  //   return this.http.get<Host[]>(baseUrl + '?itemName=' + itemName);
  // }
  //
  // getHost(hostId: string): Observable<Host> {
  //   return this.http.get<Host>(baseUrl + hostId + '/');
  // }
  //
  // createHost(host: Host): Observable<Host> {
  //   return this.http.post<Host>(baseUrl, host);
  // }
  //
  // deleteHost(hostId: string): Observable<any> {
  //   return this.http.delete<any>(baseUrl + hostId + '/');
  // }
  //
  // importHost(source: string[]): Observable<any> {
  //   return this.http.post<any>(baseUrl + 'import/', {'source': source});
  // }


}
