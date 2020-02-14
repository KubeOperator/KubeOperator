import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable, throwError} from 'rxjs';
import {ItemResource, ItemResourceDTO} from './item-resource';
import {catchError} from 'rxjs/operators';

@Injectable({
  providedIn: 'root'
})
export class ItemResourceService {

  private baseURL = '/api/v1/resource/';


  constructor(private httpClient: HttpClient) {
  }

  getResources(itemName: string, resourceType: string): Observable<ItemResourceDTO[]> {
    return this.httpClient.get<ItemResourceDTO[]>(this.baseURL + itemName + '/' + resourceType + '/');
  }

  createItemResources(itemName: string, resourceType: string, itemResources: ItemResource[]): Observable<any> {
    return this.httpClient.post(this.baseURL + itemName + '/' + resourceType + '/', itemResources).pipe(
      catchError(error => throwError(error))
    );
  }

  getItemResources(itemName: string): Observable<ItemResourceDTO[]> {
    return this.httpClient.get<ItemResourceDTO[]>(this.baseURL + itemName + '/');
  }
}
