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

  createItemResources(itemName: string, itemResources: ItemResource[]): Observable<any> {
    return this.httpClient.post(this.baseURL + itemName + '/', itemResources).pipe(
      catchError(error => throwError(error))
    );
  }

  getItemResources(itemName: string): Observable<ItemResourceDTO[]> {
    return this.httpClient.get<ItemResourceDTO[]>(this.baseURL + itemName + '/');
  }

  deleteItemResource(itemName: string, resourceType: string, resourceId: string): Observable<ItemResource> {
    return this.httpClient.delete<ItemResource>(this.baseURL + itemName + '/' + resourceType + '/' + resourceId + '/');
  }

  getClusters() {
    return this.httpClient.get<ItemResourceDTO[]>(this.baseURL + 'item/clusters/');
  }
}
