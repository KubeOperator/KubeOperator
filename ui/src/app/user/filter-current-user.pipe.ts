import {Pipe, PipeTransform} from '@angular/core';
import {User} from './user';
import {SessionUser} from '../shared/session-user';

@Pipe({
  name: 'filterCurrentUser'
})
export class FilterCurrentUserPipe implements PipeTransform {

  transform(users: User[], currentUser: SessionUser): any {
    const result = [];
    for (const user of users) {
      if (user.username !== currentUser.username) {
        result.push(user);
      }
    }
    return result;
  }

}
