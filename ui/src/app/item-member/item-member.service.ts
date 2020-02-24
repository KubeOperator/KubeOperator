import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {ItemMember, ItemMemberWrite} from './item-member';
import {Profile} from '../shared/session-user';

@Injectable({
  providedIn: 'root'
})
export class ItemMemberService {
  baseUrl = '/api/v1/profiles/';
  baseWriteUrl = '/api/v1/item/users/update/';

  constructor(private http: HttpClient) {
  }

  getProfiles(): Observable<Profile[]> {
    return this.http.get<Profile[]>(this.baseUrl);
  }

  setItemUsers(item: ItemMemberWrite, name: string) {
    return this.http.patch<ItemMemberWrite>(this.baseWriteUrl + name + '/', item);
  }
}
