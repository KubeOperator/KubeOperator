import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {ItemMember, ItemMemberWrite} from './item-member';

@Injectable({
  providedIn: 'root'
})
export class ItemMemberService {
  baseUrl = '/api/v1/item/users/';
  baseWriteUrl = '/api/v1/item/users/update/';

  constructor(private http: HttpClient) {
  }

  getItemUsers(item: string): Observable<ItemMember> {
    return this.http.get<ItemMember>(this.baseUrl + item + '/');
  }

  setItemUsers(item: ItemMemberWrite, name: string) {
    return this.http.patch<ItemMemberWrite>(this.baseWriteUrl + name + '/', item);
  }
}
