import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable, throwError} from 'rxjs';
import {User} from './user';
import {catchError} from 'rxjs/operators';
import {error} from '@angular/compiler/src/util';

const userUrl = '/api/v1/users/';

@Injectable()
export class UserService {

  constructor(private http: HttpClient) {
  }

  listUsers(): Observable<User[]> {
    return this.http.get<User[]>(userUrl);
  }

  createUser(user: User): Observable<User> {
    return this.http.post<User>(userUrl, user);
  }

  activeUser(user: User): Observable<User> {
    return this.http.patch<User>(userUrl + user.id + '/', {is_active: user.is_active});
  }

  supperUser(user: User): Observable<User> {
    return this.http.patch<User>(userUrl + user.id + '/', {is_superuser: user.is_superuser});
  }

  deleteUser(userId): Observable<any> {
    return this.http.delete(userUrl + userId + '/');
  }

  updateUser(user: User): Observable<User> {
    return this.http.patch<User>(userUrl + user.id + '/', user);
  }

  syncUserFromLDAP(): Observable<any> {
    return this.http.post<any>(userUrl + 'sync/', {});
  }

  getUser(id: number): Observable<User> {
    return this.http.get<User>(userUrl + id + '/');
  }
}
