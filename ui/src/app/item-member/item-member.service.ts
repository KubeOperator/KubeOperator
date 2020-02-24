import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Profile} from '../shared/session-user';
import {ItemMember} from './item-member';

@Injectable({
  providedIn: 'root'
})
export class ItemMemberService {
  profileUrl = '/api/v1/profiles/';

  baseUrl = '/api/v1/item/profiles/{item_name}/';

  constructor(private http: HttpClient) {
  }

  getProfiles(): Observable<Profile[]> {
    return this.http.get<Profile[]>(this.profileUrl);
  }

  getItemProfiles(itemName: string): Observable<ItemMember> {
    return this.http.get<ItemMember>(this.baseUrl.replace('{item_name}', itemName));
  }

  setItemProfiles(obj: any, itemName): Observable<any> {
    return this.http.patch<any>(this.baseUrl.replace('{item_name}', itemName), obj);
  }

}
