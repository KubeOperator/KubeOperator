import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ItemResourceService {

  private baseURL = '/api/v1/resource/';


  constructor(private httpClient: HttpClient) {
  }

  getResources(itemName: string, resourceType: string): Observable<any[]> {
    return this.httpClient.get<any[]>(this.baseURL + itemName + '/' + resourceType + '/');
  }

}
